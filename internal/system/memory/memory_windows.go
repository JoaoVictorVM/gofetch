//go:build windows

package memory

import (
	"context"
	"fmt"

	"github.com/shirou/gopsutil/v4/mem"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

func collect(ctx context.Context) (system.Section, error) {
	vm, err := mem.VirtualMemoryWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying virtual memory: %w", err)
	}

	return system.Memory{
		UsedBytes:   vm.Used,
		TotalBytes:  vm.Total,
		UsedPercent: vm.UsedPercent,
	}, nil
}
