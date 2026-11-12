package controller

import (
	"io"
	"log/slog"
	"os"
	"regexp"
	"sort"
	"strings"
	"testing"

	"github.com/victorhisamitsu/notewise/api/internal/app/perfume/controller/handler"
)

// TestRoutesMatchOpenAPIContract fails when a route is added, removed or
// renamed without updating api/openapi.yaml — the contract the web app's
// TypeScript types are generated from.
func TestRoutesMatchOpenAPIContract(t *testing.T) {
	spec, err := os.ReadFile("../../../../openapi.yaml")
	if err != nil {
		t.Fatal(err)
	}

	pathLine := regexp.MustCompile(`^  (/\S+):\s*$`)
	methodLine := regexp.MustCompile(`^    (get|post|put|patch|delete):\s*$`)
	var documented []string
	current := ""
	for _, line := range strings.Split(string(spec), "\n") {
		if m := pathLine.FindStringSubmatch(line); m != nil {
			current = m[1]
			continue
		}
		if m := methodLine.FindStringSubmatch(line); m != nil && current != "" {
			documented = append(documented, strings.ToUpper(m[1])+" "+current)
		}
	}

	h := handler.NewPerfume(nil, nil, slog.New(slog.NewTextHandler(io.Discard, nil)))
	var registered []string
	for _, r := range Routes(h) {
		registered = append(registered, r.Method+" "+r.Path)
	}

	sort.Strings(documented)
	sort.Strings(registered)
	if strings.Join(documented, "\n") != strings.Join(registered, "\n") {
		t.Fatalf("routes and openapi.yaml disagree\nregistered:\n  %s\ndocumented:\n  %s",
			strings.Join(registered, "\n  "), strings.Join(documented, "\n  "))
	}
}
