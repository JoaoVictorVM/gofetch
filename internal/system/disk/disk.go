// Package disk collects storage usage figures for the main disk.
package disk

import (
	"context"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Collector collects disk information. It implements system.Collector.
type Collector struct{}

// New returns a disk Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "disk" }

// Collect implements system.Collector. The actual collection is
// OS-specific and provided by build-tagged files.
func (Collector) Collect(ctx context.Context) (system.Section, error) {
	return collect(ctx)
}
