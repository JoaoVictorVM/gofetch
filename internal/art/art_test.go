package art

import (
	"strings"
	"testing"
)

// TestForWindows checks the windows logo is registered, consistent and
// prints it for visual inspection with `go test -v`.
func TestForWindows(t *testing.T) {
	logo, ok := For("windows")
	if !ok {
		t.Fatal(`For("windows") not found, want registered logo`)
	}
	if len(logo.Lines) == 0 {
		t.Fatal("windows logo has no lines")
	}

	width := len(logo.Lines[0])
	for i, line := range logo.Lines {
		if len(line) != width {
			t.Errorf("line %d has width %d, want %d (all lines padded equally)", i, len(line), width)
		}
	}

	t.Logf("windows logo (%dx%d):\n%s", width, len(logo.Lines), strings.Join(logo.Lines, "\n"))
}

// TestForUnknown ensures unsupported operating systems report absence
// instead of returning garbage.
func TestForUnknown(t *testing.T) {
	if _, ok := For("plan9"); ok {
		t.Error(`For("plan9") = ok, want not found`)
	}
}
