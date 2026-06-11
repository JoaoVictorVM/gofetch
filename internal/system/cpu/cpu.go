// Package cpu collects processor information: model, physical and
// logical core counts, clock and instantaneous usage.
package cpu

import (
	"context"
	"time"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// sampleInterval is how long usage sampling observes the CPU. Shorter
// reads are noisier; longer ones delay the whole program (this is the
// slowest collector and dictates total runtime).
const sampleInterval = 500 * time.Millisecond

// Collector collects CPU information. It implements system.Collector.
type Collector struct{}

// New returns a cpu Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "cpu" }

// Collect implements system.Collector. The actual collection is
// OS-specific and provided by build-tagged files.
func (Collector) Collect(ctx context.Context) (system.Section, error) {
	return collect(ctx)
}
