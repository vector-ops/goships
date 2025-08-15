package screens

import (
	"context"
	"time"

	"github.com/rthornton128/goncurses"
)

type Menu struct {
	window  *goncurses.Window
	options []string
}

func NewMenuScreen(win *goncurses.Window, options []string) *Menu {
	s := &Menu{
		window:  win,
		options: options,
	}

	return s
}

func (s *Menu) Show(ctx context.Context) {
	loading := true

	loadingCtx, cancel := context.WithCancel(ctx)
	go func() {
		time.Sleep(60 * time.Second)
		loading = false
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			if loading {
				ShowLoadingScreen(loadingCtx, s.window, nil)
			} else {
				f := "Welcome to GoShips"
				my, mx := s.window.MaxYX()
				width := len(f)
				height := 1

				y, x := (my-height)/2, (mx-width)/2

				s.window.MovePrint(y, x, f)
				s.window.Refresh()
				time.Sleep(300 * time.Millisecond)
				s.window.Erase()
			}
		}
	}
}
