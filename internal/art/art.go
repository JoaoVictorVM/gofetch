// Package art provides the ASCII-art logos displayed beside the system
// information, keyed by operating system.
package art

import (
	_ "embed"
	"strings"
)

// Art is an ASCII-art logo.
type Art struct {
	// Lines holds the logo line by line, all padded to the same width,
	// so the renderer can lay it out without re-measuring.
	Lines []string
}

//go:embed logos/windows.txt
var windowsLogo string

// logos maps a runtime.GOOS value to its logo. Populated as OS support
// lands (windows in v1; linux and darwin in future versions).
var logos = map[string]Art{
	"windows": parse(windowsLogo),
}

// For returns the logo registered for the given operating system (a
// runtime.GOOS value such as "windows"). The boolean reports whether
// the OS has a registered logo.
func For(goos string) (Art, bool) {
	a, ok := logos[goos]
	return a, ok
}

// parse splits a raw logo into lines, drops trailing blank lines and
// right-pads every line to the width of the longest one.
func parse(raw string) Art {
	lines := strings.Split(strings.TrimRight(raw, "\n"), "\n")

	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	for i, line := range lines {
		lines[i] = line + strings.Repeat(" ", width-len(line))
	}
	return Art{Lines: lines}
}
