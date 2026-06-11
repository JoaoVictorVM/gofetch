package system

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/sync/errgroup"
)

// DefaultTimeout bounds the whole collection. It needs headroom over
// the CPU usage sampling window (~500ms), the slowest collector.
const DefaultTimeout = 2 * time.Second

// Aggregate runs all collectors concurrently, bounded by timeout, and
// merges the successful results into an Info.
//
// Failure is graceful: a collector that errors leaves its section at
// the zero value (rendered as "N/A") and never takes the others down.
// The returned Info is always usable; the error joins every individual
// collector failure, for callers that want to report them.
func Aggregate(ctx context.Context, timeout time.Duration, collectors ...Collector) (*Info, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	sections := make([]Section, len(collectors))
	errs := make([]error, len(collectors))

	var g errgroup.Group
	for i, c := range collectors {
		g.Go(func() error {
			section, err := c.Collect(ctx)
			if err != nil {
				errs[i] = fmt.Errorf("collecting %s: %w", c.Name(), err)
				return nil
			}
			sections[i] = section
			return nil
		})
	}
	// Collector failures are recorded in errs instead of returned, so
	// one bad collector cannot cancel the others; Wait only synchronizes.
	_ = g.Wait()

	var info Info
	for _, section := range sections {
		if section != nil {
			section.apply(&info)
		}
	}
	return &info, errors.Join(errs...)
}
