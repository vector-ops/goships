package screens

import (
	"context"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/frames"
	"github.com/vector-ops/goships/utils"
)

func ShowWelcomeScreen(ctx context.Context, win *goncurses.Window) {

	select {
	case <-ctx.Done():
		return
	default:

		my, mx := win.MaxYX()

		win.Erase()
		for i, f := range frames.LOGO_LINES {
			width := len(f)
			height := 15

			y, x := (my-height)/2+i, (mx-width)/2

			win.MovePrint(y, x, f)
		}
		win.Refresh()
	}

	utils.Delay(1000)

}
