package render

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/JoaoVictorVM/gofetch/internal/art"
	"github.com/JoaoVictorVM/gofetch/internal/system"
)

func sampleInfo() *system.Info {
	return &system.Info{
		Host: system.Host{
			OS:       "Windows 11 Pro (build 22631)",
			Hostname: "DESKTOP-A1B2C3",
			User:     "user",
			Uptime:   3*time.Hour + 12*time.Minute,
		},
		CPU: system.CPU{
			Model:         "AMD Ryzen 5 5600X 6-Core Processor",
			PhysicalCores: 6,
			LogicalCores:  12,
			MHz:           3700,
			UsagePercent:  14,
		},
		Memory: system.Memory{
			UsedBytes:   8 << 30,
			TotalBytes:  32 << 30,
			UsedPercent: 24,
		},
		Disk: system.Disk{
			Mount:       "C:",
			UsedBytes:   412 << 30,
			TotalBytes:  931 << 30,
			UsedPercent: 44,
		},
		Board: system.Board{Manufacturer: "ASUS", Model: "TUF GAMING B550-PLUS"},
	}
}

// TestRenderSmoke renders mockup-like data and prints the layout for
// visual inspection with `go test -v`.
func TestRenderSmoke(t *testing.T) {
	logo, ok := art.For("windows")
	if !ok {
		t.Fatal("windows logo not registered")
	}

	var buf bytes.Buffer
	if err := NewText(WithColor(false)).Render(&buf, sampleInfo(), logo); err != nil {
		t.Fatalf("Render() error = %v", err)
	}
	out := buf.String()

	for _, want := range []string{
		"user@DESKTOP-A1B2C3",
		"OS",
		"Windows 11 Pro (build 22631)",
		"AMD Ryzen 5 5600X (12) @ 3.7GHz",
		"14%",
		"8.0 GiB / 32.0 GiB (24%)",
		"Disk (C:)",
		"412.0 GiB / 931.0 GiB (44%)",
		"ASUS TUF GAMING B550-PLUS",
		"3 hours, 12 mins",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q", want)
		}
	}

	t.Logf("rendered output:\n%s", out)
}

// TestRenderNotAvailable ensures failed sections degrade to N/A instead
// of empty or broken lines.
func TestRenderNotAvailable(t *testing.T) {
	logo, _ := art.For("windows")

	var buf bytes.Buffer
	if err := NewText(WithColor(false)).Render(&buf, &system.Info{}, logo); err != nil {
		t.Fatalf("Render() error = %v", err)
	}

	if got := strings.Count(buf.String(), "N/A"); got < 7 {
		t.Errorf("output has %d N/A entries, want at least 7 (all sections empty)\n%s", got, buf.String())
	}
}
