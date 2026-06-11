//go:build windows

package board

import (
	"context"
	"testing"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// TestCollectSmoke runs the real collector and logs what it sees, so the
// values can be checked against the actual machine with `go test -v`.
func TestCollectSmoke(t *testing.T) {
	section, err := New().Collect(context.Background())
	if err != nil {
		t.Fatalf("Collect() error = %v", err)
	}

	b, ok := section.(system.Board)
	if !ok {
		t.Fatalf("Collect() returned %T, want system.Board", section)
	}
	if b.Manufacturer == "" && b.Model == "" {
		t.Errorf("Collect() returned empty board: %+v", b)
	}

	t.Logf("Manufacturer: %s", b.Manufacturer)
	t.Logf("Model:        %s", b.Model)
}

// TestCollectCanceledContext ensures a dead context aborts the collection
// with an error instead of hanging or panicking.
func TestCollectCanceledContext(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := New().Collect(ctx); err == nil {
		t.Error("Collect() with canceled context returned nil error")
	}
}
