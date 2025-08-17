package main

import (
	"context"
	"log"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

const (
	GRID_WIDTH  = 11
	GRID_HEIGHT = 11
)

type Map struct {
	win *goncurses.Window

	grid *map[types.Position]types.CellType
}

func NewMap(win *goncurses.Window, startingGrid *map[types.Position]types.CellType) *Map {
	m := &Map{
		win: win,
	}

	if startingGrid != nil {
		m.grid = startingGrid
	} else {
		m.grid = m.createEmptyGrid()
	}

	return m
}

func (m *Map) Render(ctx context.Context) error {
	m.win.Erase()
	m.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	m.draw()
	m.win.NoutRefresh()
	return nil
}

func (m *Map) draw() error {
	my, mx := m.win.MaxYX()

	yfactor, xfactor := my/GRID_HEIGHT, mx/GRID_WIDTH

	for x := 1; x <= GRID_WIDTH; x++ {
		for y := 1; y <= GRID_HEIGHT; y++ {
			var ch rune
			switch (*m.grid)[types.Position{X: x, Y: y}] {
			case types.CELL_BLANK:
				ch = ' '
				m.win.ColorOn(types.COLOR_WATER)
			case types.CELL_WALL:
				ch = '|'
				m.win.ColorOn(types.COLOR_WALL)
			case types.CELL_CURSOR:
				// ch should previous ch but color reversed
				ch = ' '
				m.win.ColorOn(types.COLOR_CURSOR)
			case types.CELL_CRUISER:
				ch = '$'
				m.win.ColorOn(types.COLOR_SHIP)
			case types.CELL_DESTROYER:
				ch = '$'
				m.win.ColorOn(types.COLOR_SHIP)
			case types.CELL_BATTLESHIP:
				ch = '$'
				m.win.ColorOn(types.COLOR_SHIP)
			case types.CELL_CARRIER:
				ch = '$'
				m.win.ColorOn(types.COLOR_SHIP)
			case types.CELL_SUBMARINE:
				ch = '$'
				m.win.ColorOn(types.COLOR_SHIP)
			case types.CELL_DESTROYED:
				ch = '#'
				m.win.ColorOn(types.COLOR_FLAMES)
			case types.CELL_MISS:
				ch = 'o'
				m.win.ColorOn(types.COLOR_MISS)
			}

			m.win.MoveAddChar(y*yfactor+3, x*xfactor, goncurses.Char(ch))
		}
	}
	return nil
}

func (m *Map) Close() error {
	return m.win.Delete()
}

func (m *Map) SetEntity(entity types.Entity, o types.Orientation) {
	if !utils.ValidateEntityPosition(entity) {
		panic("Invalid entity position") // needs to be handled properly
	}

	if entity.Type == types.CRUISER {
		// chars := []rune(entity.Sprite[o])

		maxSize := len(entity.Sprite[o])
		s := 0
		switch o {
		case types.VERTICAL:
			for y := entity.StartPosition.Y; y <= entity.EndPosition.Y; y++ {
				if s > maxSize {
					log.Fatalf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
				}
				(*m.grid)[types.Position{X: entity.StartPosition.X, Y: y}] = types.CELL_DESTROYED
				s++
			}
		case types.HORIZONTAL:
			for x := entity.StartPosition.X; x <= entity.EndPosition.X; x++ {
				if s > maxSize {
					log.Fatalf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
				}
				(*m.grid)[types.Position{Y: entity.StartPosition.Y, X: x}] = types.CELL_CRUISER
				s++
			}
		}
	}
}

func (m *Map) createEmptyGrid() *map[types.Position]types.CellType {
	grid := make(map[types.Position]types.CellType)

	for x := 1; x <= GRID_WIDTH; x++ {
		for y := 1; y <= GRID_HEIGHT; y++ {
			grid[types.Position{X: x, Y: y}] = types.CELL_WALL
		}
	}

	return &grid
}
