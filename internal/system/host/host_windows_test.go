//go:build windows

package host

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

	h, ok := section.(system.Host)
	if !ok {
		t.Fatalf("Collect() returned %T, want system.Host", section)
	}
	if h.OS == "" || h.Hostname == "" || h.User == "" || h.Uptime <= 0 {
		t.Errorf("Collect() returned incomplete data: %+v", h)
	}

	t.Logf("OS:       %s", h.OS)
	t.Logf("Hostname: %s", h.Hostname)
	t.Logf("User:     %s", h.User)
	t.Logf("Uptime:   %s", h.Uptime)
}
