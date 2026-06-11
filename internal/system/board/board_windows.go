//go:build windows

package board

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/yusufpapurcu/wmi"

	"github.com/JoaoVictorVM/gofetch/internal/system"
)

// baseBoard receives the fields gofetch needs from Win32_BaseBoard.
type baseBoard struct {
	Manufacturer string
	Product      string
}

const query = "SELECT Manufacturer, Product FROM Win32_BaseBoard"

func collect(ctx context.Context) (system.Section, error) {
	// The wmi package has no context support, so the query runs in a
	// goroutine and the result is discarded if ctx expires first.
	type result struct {
		boards []baseBoard
		err    error
	}
	ch := make(chan result, 1)
	go func() {
		var boards []baseBoard
		err := wmi.Query(query, &boards)
		ch <- result{boards, err}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("querying baseboard: %w", ctx.Err())
	case r := <-ch:
		if r.err != nil {
			return nil, fmt.Errorf("querying baseboard: %w", r.err)
		}
		if len(r.boards) == 0 {
			return nil, errors.New("querying baseboard: no baseboard reported")
		}
		return system.Board{
			Manufacturer: strings.TrimSpace(r.boards[0].Manufacturer),
			Model:        strings.TrimSpace(r.boards[0].Product),
		}, nil
	}
}
