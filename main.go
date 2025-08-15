package main

import (
	"context"
	"log"

	gc "github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/screens"
	"github.com/vector-ops/goships/utils"
)

func main() {
	stdscr, err := gc.Init()
	if err != nil {
		log.Printf("Yes this is error")
		log.Fatal(err)
	}
	defer gc.End()

	gc.Echo(false)
	gc.CBreak(true)
	gc.Cursor(0)

	stdscr.Refresh()
	stdscr.Keypad(true)

	ctx, cancel := context.WithCancel(context.Background())

	go utils.QuitOnQ(stdscr, cancel)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			menu := screens.NewMenuScreen(stdscr, []string{})
			menu.Show(ctx)
		}
	}

}
