// Package board collects motherboard identification: manufacturer and
// model.
package board

import (
	"context"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Collector collects motherboard information. It implements
// system.Collector.
type Collector struct{}

// New returns a board Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "board" }

// Collect implements system.Collector. The actual collection is
// OS-specific and provided by build-tagged files.
func (Collector) Collect(ctx context.Context) (system.Section, error) {
	return collect(ctx)
}
