// Package host collects operating system and machine identity
// information: OS name and build, hostname, current user and uptime.
package host

import (
	"context"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Collector collects host information. It implements system.Collector.
type Collector struct{}

// New returns a host Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "host" }

// Collect implements system.Collector. The actual collection is
// OS-specific and provided by build-tagged files.
func (Collector) Collect(ctx context.Context) (system.Section, error) {
	return collect(ctx)
}
