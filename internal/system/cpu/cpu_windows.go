//go:build windows

package cpu

import (
	"context"
	"errors"
	"fmt"

	"github.com/shirou/gopsutil/v4/cpu"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

func collect(ctx context.Context) (system.Section, error) {
	infos, err := cpu.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying cpu info: %w", err)
	}
	if len(infos) == 0 {
		return nil, errors.New("querying cpu info: no cpus reported")
	}

	physical, err := cpu.CountsWithContext(ctx, false)
	if err != nil {
		return nil, fmt.Errorf("counting physical cores: %w", err)
	}
	logical, err := cpu.CountsWithContext(ctx, true)
	if err != nil {
		return nil, fmt.Errorf("counting logical cores: %w", err)
	}

	// A non-zero interval samples usage over that window; interval 0
	// would report usage since boot instead of the instantaneous value.
	percents, err := cpu.PercentWithContext(ctx, sampleInterval, false)
	if err != nil {
		return nil, fmt.Errorf("sampling cpu usage: %w", err)
	}
	if len(percents) == 0 {
		return nil, errors.New("sampling cpu usage: no measurement reported")
	}

	return system.CPU{
		Model:         infos[0].ModelName,
		PhysicalCores: physical,
		LogicalCores:  logical,
		MHz:           infos[0].Mhz,
		UsagePercent:  percents[0],
	}, nil
}
