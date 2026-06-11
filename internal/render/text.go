package render

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/JoaoVictorVM/gofetch/internal/art"
	"github.com/JoaoVictorVM/gofetch/internal/system"
)

const (
	notAvailable = "N/A"
	gpuDash      = "—" // v1 placeholder value for the GPU row

	// accentColor is a basic ANSI bright color, so it degrades
	// gracefully from truecolor down to 16-color terminals.
	accentColor = lipgloss.Color("12")

	logoGap = 3 // blank columns between logo and info
)

// Text renders the system information as styled terminal text: the OS
// logo on the left and the information column on the right. It
// implements Renderer.
type Text struct {
	color bool
}

// TextOption configures a Text renderer.
type TextOption func(*Text)

// WithColor enables or disables colors and styles (enabled by default).
func WithColor(enabled bool) TextOption {
	return func(t *Text) { t.color = enabled }
}

// NewText returns a Text renderer configured by opts.
func NewText(opts ...TextOption) *Text {
	t := &Text{color: true}
	for _, opt := range opts {
		opt(t)
	}
	return t
}

// Render implements Renderer.
func (t *Text) Render(w io.Writer, info *system.Info, logo art.Art) error {
	accent := lipgloss.NewStyle()
	if t.color {
		accent = accent.Foreground(accentColor).Bold(true)
	}

	infoColumn := t.infoColumn(info, accent)
	logoColumn := accent.MarginRight(logoGap).Render(strings.Join(logo.Lines, "\n"))

	out := lipgloss.JoinHorizontal(lipgloss.Top, logoColumn, infoColumn)
	if _, err := fmt.Fprintln(w, out); err != nil {
		return fmt.Errorf("writing output: %w", err)
	}
	return nil
}

type row struct {
	label string
	value string
}

func (t *Text) infoColumn(info *system.Info, accent lipgloss.Style) string {
	rows := []row{
		{"OS", valueOr(info.Host.OS)},
		{"Host", valueOr(info.Host.Hostname)},
		{"Uptime", formatUptime(info.Host.Uptime)},
		{"CPU", formatCPU(info.CPU)},
		{"CPU Usage", formatCPUUsage(info.CPU)},
		{"GPU", gpuDash},
		{"Memory", formatUsage(info.Memory.UsedBytes, info.Memory.TotalBytes, info.Memory.UsedPercent)},
		{diskLabel(info.Disk), formatUsage(info.Disk.UsedBytes, info.Disk.TotalBytes, info.Disk.UsedPercent)},
		{"Board", formatBoard(info.Board)},
	}

	labelWidth := 0
	for _, r := range rows {
		if len(r.label) > labelWidth {
			labelWidth = len(r.label)
		}
	}

	header := headerLine(info.Host)
	lines := make([]string, 0, len(rows)+2)
	lines = append(lines,
		accent.Render(header),
		strings.Repeat("─", lipgloss.Width(header)),
	)
	for _, r := range rows {
		label := fmt.Sprintf("%-*s", labelWidth+2, r.label)
		lines = append(lines, accent.Render(label)+r.value)
	}
	return strings.Join(lines, "\n")
}

// headerLine builds the "user@host" banner, falling back to whatever
// half is available, or the program name when the host collector failed.
func headerLine(h system.Host) string {
	switch {
	case h.User != "" && h.Hostname != "":
		return h.User + "@" + h.Hostname
	case h.Hostname != "":
		return h.Hostname
	case h.User != "":
		return h.User
	default:
		return "gofetch"
	}
}

func valueOr(s string) string {
	if s == "" {
		return notAvailable
	}
	return s
}

// formatCPU renders e.g. "AMD Ryzen 5 5600X (12) @ 3.7GHz". Vendor
// noise embedded in the model name — Intel's "CPU @ 2.90GHz" clock,
// AMD's "6-Core Processor" — is stripped so it is not shown twice.
func formatCPU(c system.CPU) string {
	if c.Model == "" {
		return notAvailable
	}
	model := c.Model
	if at := strings.IndexByte(model, '@'); at >= 0 {
		model = model[:at]
	}
	if i := strings.Index(model, "-Core"); i >= 0 {
		if sp := strings.LastIndexByte(model[:i], ' '); sp >= 0 {
			model = model[:sp]
		}
	}
	model = strings.TrimSuffix(strings.TrimSpace(model), " CPU")

	out := model
	if c.LogicalCores > 0 {
		out = fmt.Sprintf("%s (%d)", out, c.LogicalCores)
	}
	if c.MHz > 0 {
		out = fmt.Sprintf("%s @ %.1fGHz", out, c.MHz/1000)
	}
	return out
}

func formatCPUUsage(c system.CPU) string {
	if c.Model == "" {
		return notAvailable
	}
	return fmt.Sprintf("%.0f%%", c.UsagePercent)
}

// formatUsage renders e.g. "7.8 GiB / 31.9 GiB (24%)", or "N/A" when
// the section was not collected.
func formatUsage(used, total uint64, percent float64) string {
	if total == 0 {
		return notAvailable
	}
	return fmt.Sprintf("%s / %s (%.0f%%)", formatGiB(used), formatGiB(total), percent)
}

func formatGiB(bytes uint64) string {
	const gib = 1 << 30
	return fmt.Sprintf("%.1f GiB", float64(bytes)/gib)
}

// formatUptime renders e.g. "3 hours, 12 mins", including days when
// relevant and seconds only for sub-minute uptimes.
func formatUptime(d time.Duration) string {
	if d <= 0 {
		return notAvailable
	}
	if d < time.Minute {
		return plural(int(d.Seconds()), "sec")
	}

	var parts []string
	if days := int(d.Hours()) / 24; days > 0 {
		parts = append(parts, plural(days, "day"))
	}
	if hours := int(d.Hours()) % 24; hours > 0 {
		parts = append(parts, plural(hours, "hour"))
	}
	if mins := int(d.Minutes()) % 60; mins > 0 {
		parts = append(parts, plural(mins, "min"))
	}
	if len(parts) == 0 {
		return plural(int(d.Minutes()), "min")
	}
	return strings.Join(parts, ", ")
}

func plural(n int, unit string) string {
	if n == 1 {
		return fmt.Sprintf("1 %s", unit)
	}
	return fmt.Sprintf("%d %ss", n, unit)
}

func diskLabel(d system.Disk) string {
	if d.Mount == "" {
		return "Disk"
	}
	return fmt.Sprintf("Disk (%s)", d.Mount)
}

func formatBoard(b system.Board) string {
	switch {
	case b.Manufacturer != "" && b.Model != "":
		return b.Manufacturer + " " + b.Model
	case b.Model != "":
		return b.Model
	case b.Manufacturer != "":
		return b.Manufacturer
	default:
		return notAvailable
	}
}
