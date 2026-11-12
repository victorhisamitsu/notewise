// Command api is the Notewise HTTP API. It is the composition root: it
// only wires dependencies, all behavior lives in internal/.
//
// On Vercel it runs with the Go framework preset, which builds this
// entrypoint and provides the PORT variable.
package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/handler"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/usecase"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/domain/service"
	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/repository/scento"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/config"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/httpx"
	"github.com/victorhisamitsu/notewise/api/internal/crosscutting/logger"
)

func main() {
	cfg := config.Load()
	log := logger.New(cfg.LogLevel)

	// Repository
	provider := scento.NewClient(cfg.ScentoBaseURL, cfg.ScentoTimeout, cfg.CacheTTL, cfg.UserAgent, log)

	// Domain services
	catalog := service.NewCatalog(provider, log)
	recommender := service.NewRecommender(provider, log)

	// Use cases and handler
	h := handler.NewPerfume(
		usecase.NewDiscovery(catalog, scento.Attribution),
		usecase.NewPersonal(recommender, scento.Attribution),
		log,
	)

	mux := http.NewServeMux()
	controller.RegisterRoutes(mux, h)
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, _ *http.Request) {
		httpx.WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
	})

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           httpx.Recover(log, httpx.Logging(log, httpx.CORS(cfg.AllowedOrigins, mux))),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
	}
	go func() {
		log.Info("notewise api listening", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("server stopped", "err", err)
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
