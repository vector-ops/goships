package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

const (
	DEFAULT_GRID_WIDTH  = 9
	DEFAULT_GRID_HEIGHT = 7
	CELL_HEIGHT         = 2
	CELL_WIDTH          = 5
)

type Cell struct {
	cellType types.CellType
	color    int16
	content  string
}

type Map struct {
	win *goncurses.Window

	grid       *map[types.Position]Cell
	gridHeight int
	gridWidth  int
}

func NewMap(win *goncurses.Window, startingGrid *map[types.Position]Cell, gridWidth, gridHeight *int) *Map {
	m := &Map{
		win:        win,
		gridHeight: DEFAULT_GRID_HEIGHT,
		gridWidth:  DEFAULT_GRID_WIDTH,
	}

	if startingGrid != nil {
		m.grid = startingGrid
	} else {
		m.grid = m.createEmptyGrid()
	}

	if gridWidth != nil {
		m.gridWidth = *gridWidth
	}
	if gridHeight != nil {
		m.gridHeight = *gridHeight
	}

	return m
}

func (m *Map) Render(ctx context.Context) error {
	m.win.Erase()
	err := m.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}
	err = m.drawBorders()
	if err != nil {
		return err
	}
	m.draw()
	m.win.NoutRefresh()
	return nil
}

func (m *Map) draw() error {
	my, mx := m.win.MaxYX()
	xOffset := CELL_WIDTH / 2
	yOffset := CELL_HEIGHT / 2

	startX := (mx - m.gridWidth*CELL_WIDTH) / 2
	startY := (my - m.gridHeight*CELL_HEIGHT) / 2

	for col := 0; col < m.gridWidth; col++ {
		for row := 0; row < m.gridHeight; row++ {

			cell := (*m.grid)[types.Position{X: col, Y: row}]

			x := (startX + xOffset) + col*CELL_WIDTH
			y := (startY + yOffset) + row*CELL_HEIGHT

			m.win.ColorOn(cell.color)
			m.win.MovePrint(y, x, cell.content)
			m.win.ColorOff(cell.color)
		}
	}
	return nil
}

func (m *Map) Close() error {
	return m.win.Delete()
}

func (m *Map) SetEntity(entity types.Entity, o types.Orientation) error {
	if !utils.ValidateEntityPosition(entity) {
		return errors.New("Invalid entity position")
	}

	maxSize := len(entity.Sprite[o])
	s := 0
	switch o {
	case types.VERTICAL:
		for y := entity.StartPosition.Y; y <= entity.EndPosition.Y; y++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: entity.StartPosition.X, Y: y}] = Cell{
				content:  fmt.Sprintf(" %c ", entity.Sprite[o][s]),
				cellType: entity.CellType,
				color:    entity.Color,
			}
			s++
		}
	case types.HORIZONTAL:
		for x := entity.StartPosition.X; x <= utils.ExpectedEndPosition(entity.StartPosition.X, entity.Sprite[o]); x++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: x, Y: entity.StartPosition.Y}] = Cell{
				content:  fmt.Sprintf(" %c ", entity.Sprite[o][s]),
				cellType: entity.CellType,
				color:    entity.Color,
			}
			s++
		}
	}

	return nil
}

func (m *Map) createEmptyGrid() *map[types.Position]Cell {
	grid := make(map[types.Position]Cell)

	for x := 0; x < m.gridWidth; x++ {
		for y := 0; y < m.gridHeight; y++ {
			grid[types.Position{X: x, Y: y}] = Cell{
				content:  "   ",
				cellType: types.CELL_WATER,
				color:    types.COLOR_WATER,
			}
		}
	}

	return &grid
}

func (m *Map) drawBorders() error {
	my, mx := m.win.MaxYX()

	if mx < m.gridWidth*CELL_WIDTH || my < m.gridHeight*CELL_HEIGHT {
		return fmt.Errorf("Window size is too small to accommodate the map. Map is %dx%d, window is %dx%d", m.gridWidth*CELL_WIDTH, m.gridHeight*CELL_HEIGHT, mx, my)
	}

	// calculate the starting position for the drawBorders
	// borders start is (max size - (grid size * cell size)) / 2
	// this gives the starting position of the grid such that it is centered in the window
	startX := (mx - m.gridWidth*CELL_WIDTH) / 2
	startY := (my - m.gridHeight*CELL_HEIGHT) / 2

	for row := 0; row <= m.gridHeight; row++ {
		for col := 0; col <= m.gridWidth; col++ {
			y := startY + row*CELL_HEIGHT
			x := startX + col*CELL_WIDTH

			m.win.ColorOn(types.COLOR_WALL)
			m.win.MoveAddChar(y, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_CORNER]))

			if col < m.gridWidth {
				for i := 1; i < CELL_WIDTH; i++ {
					m.win.MoveAddChar(y, x+i, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_HORIZONTAL]))
				}
			}

			if row < m.gridHeight {
				m.win.MoveAddChar(y+1, x, goncurses.Char(types.WALLS_ASCII[types.CELL_WALL_VERTICAL]))
			}
			m.win.ColorOff(types.COLOR_WALL)
		}
	}

	return nil
}
