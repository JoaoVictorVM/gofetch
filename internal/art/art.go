// Package art provides the ASCII-art logos displayed beside the system
// information, keyed by operating system.
package art

// Art is an ASCII-art logo.
type Art struct {
	// Lines holds the logo line by line, all visually the same width,
	// so the renderer can lay it out without re-measuring.
	Lines []string
}

// logos maps a runtime.GOOS value to its logo. Populated as OS support
// lands (windows in v1; linux and darwin in future versions).
var logos = map[string]Art{}

// For returns the logo registered for the given operating system (a
// runtime.GOOS value such as "windows"). The boolean reports whether
// the OS has a registered logo.
func For(goos string) (Art, bool) {
	a, ok := logos[goos]
	return a, ok
}
