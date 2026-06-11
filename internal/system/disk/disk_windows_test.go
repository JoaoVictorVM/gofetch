//go:build windows

package disk

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

	d, ok := section.(system.Disk)
	if !ok {
		t.Fatalf("Collect() returned %T, want system.Disk", section)
	}
	if d.Mount == "" || d.TotalBytes == 0 || d.UsedBytes == 0 || d.UsedBytes > d.TotalBytes {
		t.Errorf("Collect() returned inconsistent data: %+v", d)
	}
	if d.UsedPercent <= 0 || d.UsedPercent > 100 {
		t.Errorf("UsedPercent = %v, want within (0, 100]", d.UsedPercent)
	}

	const gib = 1 << 30
	t.Logf("Mount: %s", d.Mount)
	t.Logf("Used:  %.0f GiB", float64(d.UsedBytes)/gib)
	t.Logf("Total: %.0f GiB", float64(d.TotalBytes)/gib)
	t.Logf("Usage: %.1f%%", d.UsedPercent)
}
