package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	_ "net/http/pprof"

	gc "github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/game"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

func main() {
	var debug bool
	args := os.Args[1:]

	if len(args) != 0 {
		if args[0] == "--debug" || args[0] == "-d" {
			wd, err := os.Getwd()
			if err != nil {
				wd = "unknown"
			}
			debug = true
			fmt.Println("Debug mode enabled. You can save game state by pressing 'u'")
			fmt.Printf("Find logs in the logs directory: %s/logs\n", wd)
			if len(args) == 2 && args[1] == "0" {
			} else {
				fmt.Println("Starting game in 5 seconds...")
				fmt.Println("Use 0 with --debug flag to skip this message")
				time.Sleep(time.Second * 5)
			}
		}
	}

	stdscr, err := gc.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer gc.End()

	h, w := gc.StdScr().MaxYX()
	if h < 32 || w < 40 {
		gc.End() // restore terminal before printing errors
		fmt.Fprintf(os.Stderr, "Terminal too small. Minimum size: 32x40. Current size: %dx%d\n", h, w)
		return
	}

	if err := gc.StartColor(); err != nil {
		gc.End()
		fmt.Fprintf(os.Stderr, "Failed to start color mode: %v\n", err)
		return
	}

	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)

	stdscr.Refresh()
	stdscr.Keypad(true)

	gc.UseDefaultColors()

	gc.InitPair(types.ColorWater, gc.C_BLUE, gc.C_BLUE)
	gc.InitPair(types.ColorCursor, gc.C_BLACK, gc.C_WHITE)
	gc.InitPair(types.ColorHit, gc.C_BLACK, gc.C_RED)
	gc.InitPair(types.ColorWall, -1, -1)
	gc.InitPair(types.ColorShip, gc.C_BLACK, gc.C_GREEN)
	gc.InitPair(types.ColorMiss, gc.C_BLACK, gc.C_BLUE)
	gc.InitPair(types.WhiteBlack, gc.C_WHITE, -1)
	gc.InitPair(types.RedBlack, gc.C_RED, -1)
	gc.InitPair(types.GreenBlack, gc.C_GREEN, -1)
	gc.InitPair(types.BlueBlack, gc.C_BLUE, -1)
	gc.InitPair(types.YellowBlack, gc.C_YELLOW, -1)
	gc.InitPair(types.MagentaBlack, gc.C_MAGENTA, -1)
	gc.InitPair(types.CyanBlack, gc.C_CYAN, -1)
	gc.InitPair(types.WhiteRed, gc.C_WHITE, gc.C_RED)
	gc.InitPair(types.WhiteGreen, gc.C_WHITE, gc.C_GREEN)
	gc.InitPair(types.WhiteBlue, gc.C_WHITE, gc.C_BLUE)
	gc.InitPair(types.BlackYellow, gc.C_BLACK, gc.C_YELLOW)
	gc.InitPair(types.BlackCyan, gc.C_BLACK, gc.C_CYAN)
	gc.InitPair(types.BlackMagenta, gc.C_BLACK, gc.C_MAGENTA)
	gc.InitPair(types.BlackRed, gc.C_BLACK, gc.C_RED)

	ctx, cancel := context.WithCancel(context.Background())

	ch := make(chan gc.Key)

	go utils.HandleKeyboardEvent(stdscr, cancel, ch)

	game := game.NewGameState(stdscr, ch, debug)

	// clear previous logs
	utils.RemoveFilesByPattern("logs/*.log")
	utils.RemoveFilesByPattern("logs/*.json")

	if err := game.Render(ctx, cancel); err != nil {
		game.CloseResources()
		cancel()
		gc.End() // to restore terminal
		fmt.Fprintf(os.Stderr, "Error rendering game: %v\n", err)
	}
}
