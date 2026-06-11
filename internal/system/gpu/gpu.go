// Package gpu reserves the GPU field in the output. In v1 it is a
// placeholder: the renderer prints only the label, with no data. Model
// and usage collection come in v1.1 (WMI Win32_VideoController).
package gpu

import (
	"context"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Collector is the v1 GPU placeholder. It implements system.Collector
// and always succeeds with an empty system.GPU.
type Collector struct{}

// New returns a gpu Collector.
func New() Collector { return Collector{} }

// Name implements system.Collector.
func (Collector) Name() string { return "gpu" }

// Collect implements system.Collector.
func (Collector) Collect(_ context.Context) (system.Section, error) {
	return system.GPU{}, nil
}
