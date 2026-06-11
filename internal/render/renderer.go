// Package render defines how collected system information is presented
// to the user.
package render

import (
	"io"

	"github.com/JoaoVictorVM/gofetch/internal/art"
	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// Renderer produces the final user-facing output for collected system
// information. Implementations decide the format: the v1 text renderer
// draws the logo and a styled info column side by side; future versions
// may add e.g. a JSON renderer without touching collection.
type Renderer interface {
	Render(w io.Writer, info *system.Info, logo art.Art) error
}
