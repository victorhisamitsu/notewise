package cache

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestGetOrLoad_CachesAndExpires(t *testing.T) {
	c := NewTTL[int](time.Minute, 10)
	now := time.Now()
	c.now = func() time.Time { return now }

	var calls int
	load := func(context.Context) (int, error) { calls++; return 42, nil }

	for i := 0; i < 3; i++ {
		v, _, err := c.GetOrLoad(context.Background(), "k", load)
		if err != nil || v != 42 {
			t.Fatalf("got %d, %v", v, err)
		}
	}
	if calls != 1 {
		t.Fatalf("calls = %d, want 1", calls)
	}

	now = now.Add(2 * time.Minute)
	if _, hit, _ := c.GetOrLoad(context.Background(), "k", load); hit || calls != 2 {
		t.Fatalf("expected reload after expiry (hit=%v calls=%d)", hit, calls)
	}
}

func TestGetOrLoad_CoalescesConcurrentLoads(t *testing.T) {
	c := NewTTL[string](time.Minute, 10)
	var calls atomic.Int64
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, _, _ = c.GetOrLoad(context.Background(), "hot", func(context.Context) (string, error) {
				calls.Add(1)
				time.Sleep(30 * time.Millisecond)
				return "v", nil
			})
		}()
	}
	wg.Wait()
	if calls.Load() != 1 {
		t.Fatalf("calls = %d, want 1", calls.Load())
	}
}

func TestGetOrLoad_DoesNotCacheErrors(t *testing.T) {
	c := NewTTL[int](time.Minute, 10)
	boom := errors.New("boom")
	if _, _, err := c.GetOrLoad(context.Background(), "k", func(context.Context) (int, error) { return 0, boom }); !errors.Is(err, boom) {
		t.Fatalf("err = %v", err)
	}
	if _, ok := c.Get("k"); ok {
		t.Fatal("errors must not be cached")
	}
}

func TestSet_IsBounded(t *testing.T) {
	c := NewTTL[int](time.Minute, 3)
	for i, k := range []string{"a", "b", "c", "d"} {
		c.Set(k, i)
	}
	if len(c.items) > 3 {
		t.Fatalf("len = %d, want <= 3", len(c.items))
	}
}
