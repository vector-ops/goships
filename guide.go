package main

import (
	"context"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
)

type Guide struct {
	win   *goncurses.Window
	debug bool

	title      string
	titleColor int16
}

func NewGuide(win *goncurses.Window, debug bool) *Guide {
	return &Guide{
		win:        win,
		debug:      debug,
		title:      "GUIDE",
		titleColor: types.BlueBlack,
	}
}

func (g *Guide) Render(ctx context.Context) error {
	g.win.Erase()
	err := g.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}

	_, mx := g.win.MaxYX()

	g.win.ColorOn(g.titleColor)
	g.win.MovePrint(1, (mx/2)-len(g.title)/2, g.title)
	g.win.ColorOff(g.titleColor)

	g.draw()
	g.win.NoutRefresh()

	return nil
}

func (g *Guide) Close() error {
	g.win.Delete()
	return nil
}

func (g *Guide) draw() {
	// my, mx := g.win.MaxYX()

	offsetX := 2

	// Water
	g.win.ColorOn(types.ColorWater)
	g.win.MovePrint(2, 0+offsetX, "   ")
	g.win.ColorOff(types.ColorWater)
	g.win.MovePrint(2, 4+offsetX, "- Water")

	// Hit
	g.win.ColorOn(types.ColorHit)
	g.win.MovePrint(3, 0+offsetX, "   ")
	g.win.ColorOff(types.ColorHit)
	g.win.MovePrint(3, 4+offsetX, "- Hit")

	// Miss
	g.win.ColorOn(types.ColorMiss)
	g.win.MovePrint(4, 0+offsetX, " x ")
	g.win.ColorOff(types.ColorMiss)
	g.win.MovePrint(4, 4+offsetX, "- Miss")

	// Ship
	g.win.ColorOn(types.ColorShip)
	g.win.MovePrint(5, 0+offsetX, "   ")
	g.win.ColorOff(types.ColorShip)
	g.win.MovePrint(5, 4+offsetX, "- Ship")

	// Cursor
	g.win.ColorOn(types.ColorCursor)
	g.win.MovePrint(6, 0+offsetX, "   ")
	g.win.ColorOff(types.ColorCursor)
	g.win.MovePrint(6, 4+offsetX, "- Cursor")

	// Ships
	g.win.MovePrint(7, 0+offsetX, "BBSHP - Battleship")
	g.win.MovePrint(8, 0+offsetX, "CRUZ -  Cruizer")
	g.win.MovePrint(9, 0+offsetX, "SS - Submarine")
	g.win.MovePrint(10, 0+offsetX, "DDS - Destroyer")
	g.win.MovePrint(11, 0+offsetX, "CVSHP - Aircraft Carrier")
}
