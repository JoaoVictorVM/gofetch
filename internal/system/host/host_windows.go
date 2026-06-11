//go:build windows

package host

import (
	"context"
	"fmt"
	"os/user"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v4/host"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

func collect(ctx context.Context) (system.Section, error) {
	info, err := host.InfoWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("querying host info: %w", err)
	}

	username := ""
	if u, err := user.Current(); err == nil {
		// Windows reports DOMAIN\user; keep only the user part.
		username = u.Username[strings.LastIndexByte(u.Username, '\\')+1:]
	}

	return system.Host{
		OS:       formatOS(info.Platform, info.KernelVersion),
		Hostname: info.Hostname,
		User:     username,
		Uptime:   time.Duration(info.Uptime) * time.Second,
	}, nil
}

// formatOS turns gopsutil's platform name and kernel version, e.g.
// "Microsoft Windows 11 Pro" and "10.0.22631.4037 Build 22631.4037",
// into "Windows 11 Pro (build 22631)". The build number is the third
// component of the dotted kernel version.
func formatOS(platform, kernelVersion string) string {
	os := strings.TrimPrefix(platform, "Microsoft ")
	if fields := strings.Fields(kernelVersion); len(fields) > 0 {
		if parts := strings.Split(fields[0], "."); len(parts) >= 3 {
			return fmt.Sprintf("%s (build %s)", os, parts[2])
		}
	}
	return os
}
