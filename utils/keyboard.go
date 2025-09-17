package utils

import (
	"context"

	gc "github.com/rthornton128/goncurses"
)

func HandleKeyboardEvent(win *gc.Window, cancel context.CancelFunc, keyEventChan chan<- gc.Key) {
	for {
		ch := win.GetChar()
		switch ch {
		case 'q':
			cancel()
		default:
			keyEventChan <- ch
		}
	}
}
