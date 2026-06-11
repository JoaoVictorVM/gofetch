// Package memory collects RAM usage figures: used, total and percent.
package memory

import (
	"context"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Collector collects memory information. It implements system.Collector.
type Collector struct{}

// New returns a memory Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "memory" }

// Collect implements system.Collector. The actual collection is
// OS-specific and provided by build-tagged files.
func (Collector) Collect(ctx context.Context) (system.Section, error) {
	return collect(ctx)
}
