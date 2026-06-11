//go:build windows

package system_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/JoaoVictorVM/gofetch/internal/system"
	"github.com/JoaoVictorVM/gofetch/internal/system/board"
	"github.com/JoaoVictorVM/gofetch/internal/system/cpu"
	"github.com/JoaoVictorVM/gofetch/internal/system/disk"
	"github.com/JoaoVictorVM/gofetch/internal/system/gpu"
	"github.com/JoaoVictorVM/gofetch/internal/system/host"
	"github.com/JoaoVictorVM/gofetch/internal/system/memory"
)

func realCollectors() []system.Collector {
	return []system.Collector{
		host.New(), cpu.New(), memory.New(), disk.New(), board.New(), gpu.New(),
	}
}

// TestAggregateSmoke runs every real collector concurrently and logs the
// merged result and elapsed time for manual inspection with `go test -v`.
// Total time must be close to the slowest collector (the ~500ms CPU
// sampling), not the sum of all collectors.
func TestAggregateSmoke(t *testing.T) {
	start := time.Now()
	info, err := system.Aggregate(context.Background(), system.DefaultTimeout, realCollectors()...)
	elapsed := time.Since(start)

	if err != nil {
		t.Errorf("Aggregate() reported collector failures: %v", err)
	}
	if info.Host.OS == "" || info.CPU.Model == "" || info.Memory.TotalBytes == 0 ||
		info.Disk.TotalBytes == 0 || info.Board.Model == "" {
		t.Errorf("Aggregate() returned incomplete info: %+v", info)
	}
	if elapsed > 1500*time.Millisecond {
		t.Errorf("Aggregate() took %v, want well under 1.5s (concurrent, not sequential)", elapsed)
	}

	t.Logf("elapsed: %v", elapsed)
	t.Logf("info: %+v", info)
}

// failingCollector simulates a broken subsystem.
type failingCollector struct{}

func (failingCollector) Name() string { return "broken" }
func (failingCollector) Collect(context.Context) (system.Section, error) {
	return nil, errors.New("boom")
}

// TestAggregateGracefulFailure ensures one broken collector neither stops
// the others nor panics: its section stays zero and the error is reported.
func TestAggregateGracefulFailure(t *testing.T) {
	collectors := append(realCollectors(), failingCollector{})

	info, err := system.Aggregate(context.Background(), system.DefaultTimeout, collectors...)

	if err == nil {
		t.Error("Aggregate() with a failing collector returned nil error")
	}
	if info.Host.OS == "" || info.CPU.Model == "" {
		t.Errorf("Aggregate() lost healthy sections after one failure: %+v", info)
	}
}
