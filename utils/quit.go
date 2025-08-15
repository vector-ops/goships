package utils

import (
	"context"

	gc "github.com/rthornton128/goncurses"
)

func QuitOnQ(win *gc.Window, cancel context.CancelFunc) {
	for {
		switch win.GetChar() {
		case 'q':
			cancel()
		}
	}
}
