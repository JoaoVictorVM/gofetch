package system

import "context"

// Info is a piece of system information produced by a Collector.
// Implementations merge themselves into a SystemInfo; the unexported
// method keeps the set of implementations closed to this package.
type Info interface {
	apply(*SystemInfo)
}

// Collector gathers the information of a single subsystem (host, cpu,
// memory, ...). Implementations must honor ctx cancellation and return
// an error instead of panicking, so a failing subsystem degrades to
// "N/A" without taking the whole output down.
type Collector interface {
	// Name identifies the collector in error messages, e.g. "cpu".
	Name() string
	// Collect gathers the subsystem information.
	Collect(ctx context.Context) (Info, error)
}
