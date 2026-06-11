//go:build windows

package disk

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/shirou/gopsutil/v4/disk"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

func collect(ctx context.Context) (system.Section, error) {
	drive := systemDrive()

	usage, err := disk.UsageWithContext(ctx, drive+`\`)
	if err != nil {
		return nil, fmt.Errorf("querying usage of %s: %w", drive, err)
	}

	return system.Disk{
		Mount:       drive,
		UsedBytes:   usage.Used,
		TotalBytes:  usage.Total,
		UsedPercent: usage.UsedPercent,
	}, nil
}

// systemDrive returns the Windows system drive (e.g. "C:"), which is
// what gofetch treats as the main disk.
func systemDrive() string {
	if d := strings.TrimSuffix(os.Getenv("SystemDrive"), `\`); d != "" {
		return d
	}
	return "C:"
}
