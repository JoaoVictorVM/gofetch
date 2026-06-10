package system

import "context"

// Section is one piece of system information produced by a Collector.
// Implementations merge themselves into an Info; the unexported method
// keeps the set of implementations closed to this package.
type Section interface {
	apply(*Info)
}

// Collector gathers the information of a single subsystem (host, cpu,
// memory, ...). Implementations must honor ctx cancellation and return
// an error instead of panicking, so a failing subsystem degrades to
// "N/A" without taking the whole output down.
type Collector interface {
	// Name identifies the collector in error messages, e.g. "cpu".
	Name() string
	// Collect gathers the subsystem information.
	Collect(ctx context.Context) (Section, error)
}
