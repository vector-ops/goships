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
		log.Fatal(err)
	}
	defer gc.End()

	h, w := gc.StdScr().MaxYX()
	if h < 32 || w < 40 {
		gc.End() // restore terminal
		fmt.Fprintf(os.Stderr, "Terminal too small. Minimum size: 32x40. Current size: %dx%d\n", h, w)
		return
	}

	if err := gc.StartColor(); err != nil {
		gc.End() // restore terminal
		fmt.Fprintf(os.Stderr, "Failed to start color mode: %v\n", err)
		return
	}

	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)

	stdscr.Refresh()
	stdscr.Keypad(true)

	gc.UseDefaultColors()

	gc.InitPair(types.COLOR_WATER, gc.C_BLUE, gc.C_BLUE)
	gc.InitPair(types.COLOR_CURSOR, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(types.COLOR_HIT, gc.C_RED, gc.C_RED)
	gc.InitPair(types.COLOR_WALL, -1, -1)
	gc.InitPair(types.COLOR_SHIP, gc.C_BLACK, gc.C_GREEN)
	gc.InitPair(types.COLOR_MISS, gc.C_BLACK, gc.C_BLUE)
	gc.InitPair(types.COLOR_TITLE_ENEMY, gc.C_RED, -1)
	gc.InitPair(types.COLOR_TITLE_PLAYER, gc.C_GREEN, -1)
	gc.InitPair(types.WHITE_BLACK, gc.C_WHITE, -1)
	gc.InitPair(types.RED_BLACK, gc.C_RED, -1)
	gc.InitPair(types.GREEN_BLACK, gc.C_GREEN, -1)
	gc.InitPair(types.BLUE_BLACK, gc.C_BLUE, -1)
	gc.InitPair(types.YELLOW_BLACK, gc.C_YELLOW, -1)
	gc.InitPair(types.MAGENTA_BLACK, gc.C_MAGENTA, -1)
	gc.InitPair(types.CYAN_BLACK, gc.C_CYAN, -1)
	gc.InitPair(types.WHITE_RED, gc.C_WHITE, gc.C_RED)
	gc.InitPair(types.WHITE_GREEN, gc.C_WHITE, gc.C_GREEN)
	gc.InitPair(types.WHITE_BLUE, gc.C_WHITE, gc.C_BLUE)
	gc.InitPair(types.BLACK_YELLOW, gc.C_BLACK, gc.C_YELLOW)
	gc.InitPair(types.BLACK_CYAN, gc.C_BLACK, gc.C_CYAN)
	gc.InitPair(types.BLACK_MAGENTA, gc.C_BLACK, gc.C_MAGENTA)

	ctx, cancel := context.WithCancel(context.Background())

	go utils.QuitOnQ(stdscr, cancel)

	game := NewGameState(stdscr)

	if err := game.Render(ctx, cancel); err != nil {
		game.CloseResources()
		cancel()
		gc.End() // to restore terminal
		fmt.Fprintf(os.Stderr, "Error rendering game: %v\n", err)
	}
}
