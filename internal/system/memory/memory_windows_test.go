//go:build windows

package memory

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

	m, ok := section.(system.Memory)
	if !ok {
		t.Fatalf("Collect() returned %T, want system.Memory", section)
	}
	if m.TotalBytes == 0 || m.UsedBytes == 0 || m.UsedBytes > m.TotalBytes {
		t.Errorf("Collect() returned inconsistent data: %+v", m)
	}
	if m.UsedPercent <= 0 || m.UsedPercent > 100 {
		t.Errorf("UsedPercent = %v, want within (0, 100]", m.UsedPercent)
	}

	const gib = 1 << 30
	t.Logf("Used:  %.1f GiB", float64(m.UsedBytes)/gib)
	t.Logf("Total: %.1f GiB", float64(m.TotalBytes)/gib)
	t.Logf("Usage: %.1f%%", m.UsedPercent)
}
