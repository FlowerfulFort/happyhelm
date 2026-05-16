//go:build !windows

package tui

import (
	"fmt"
	"os"
)

func openTTY() (*os.File, error) {
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil, fmt.Errorf("open terminal for interactive picker: %w; use --no-tui when no terminal is available", err)
	}
	return tty, nil
}
