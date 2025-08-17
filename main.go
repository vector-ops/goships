package main

import (
	"context"
	"fmt"
	"log"
	"os"

	gc "github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

func main() {
	stdscr, err := gc.Init()
	if err != nil {
		log.Printf("Yes this is error")
		log.Fatal(err)
	}
	defer gc.End()

	h, w := gc.StdScr().MaxYX()
	if h < 32 || w < 40 {
		fmt.Fprintf(os.Stderr, "Terminal too small. Minimum size: 32x40. Current size: %dx%d\n", h, w)
		return
	}

	if err := gc.StartColor(); err != nil {
		log.Fatal(err)
	}

	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)

	stdscr.Refresh()
	stdscr.Keypad(true)

	gc.UseDefaultColors()

	gc.InitPair(types.COLOR_WATER, -1, -1)
	gc.InitPair(types.COLOR_CURSOR, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(types.COLOR_FLAMES, gc.C_RED, -1)
	gc.InitPair(types.COLOR_WALL, gc.C_WHITE, -1)
	gc.InitPair(types.COLOR_SHIP, gc.C_GREEN, -1)
	gc.InitPair(types.COLOR_MISS, gc.C_BLUE, -1)

	ctx, cancel := context.WithCancel(context.Background())

	go utils.QuitOnQ(stdscr, cancel)

	game := NewGameState(stdscr)

	if err := game.Render(ctx, cancel); err != nil {
		log.Fatal(err)
	}

}
