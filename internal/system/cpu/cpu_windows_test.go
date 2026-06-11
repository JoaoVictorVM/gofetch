//go:build windows

package cpu

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

	c, ok := section.(system.CPU)
	if !ok {
		t.Fatalf("Collect() returned %T, want system.CPU", section)
	}
	if c.Model == "" || c.PhysicalCores <= 0 || c.LogicalCores <= 0 || c.MHz <= 0 {
		t.Errorf("Collect() returned incomplete data: %+v", c)
	}
	if c.UsagePercent < 0 || c.UsagePercent > 100 {
		t.Errorf("UsagePercent = %v, want within [0, 100]", c.UsagePercent)
	}

	t.Logf("Model:    %s", c.Model)
	t.Logf("Cores:    %d physical, %d logical", c.PhysicalCores, c.LogicalCores)
	t.Logf("Clock:    %.0f MHz", c.MHz)
	t.Logf("Usage:    %.1f%%", c.UsagePercent)
}
