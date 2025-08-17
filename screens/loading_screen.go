package screens

import (
	"context"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/frames"
	"github.com/vector-ops/goships/utils"
)

func ShowLoadingScreen(ctx context.Context, win *goncurses.Window, customFrames *[]string) {

	renderFrames := frames.DOUBLE_TURRET

	if customFrames != nil {
		renderFrames = *customFrames
	}

	for {
		for _, f := range renderFrames {
			select {
			case <-ctx.Done():
				return
			default:
				my, mx := win.MaxYX()
				width := len(f)
				height := 1

				y, x := (my-height)/2, (mx-width)/2

				win.Erase()
				win.Refresh()

				win.MovePrint(y, x, f)
				win.Refresh()
				utils.Delay(300)
			}
		}
	}

}
