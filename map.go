package main

import (
	"context"
	"fmt"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

const (
	DEFAULT_GRID_WIDTH  = 11
	DEFAULT_GRID_HEIGHT = 11
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
	m.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	m.draw()
	m.win.NoutRefresh()
	return nil
}

func (m *Map) draw() error {
	my, mx := m.win.MaxYX()

	yfactor, xfactor := my/m.gridHeight, mx/m.gridWidth

	for x := 1; x <= m.gridWidth; x++ {
		for y := 1; y <= m.gridHeight; y++ {
			// var ch rune
			cell := (*m.grid)[types.Position{X: x, Y: y}]
			// switch cell {
			// case types.CELL_BLANK:
			// 	ch = ' '
			// 	m.win.ColorOn(types.COLOR_WATER)
			// case types.CELL_CURSOR:
			// 	// ch should previous ch but color reversed
			// 	ch = ' '
			// 	m.win.ColorOn(types.COLOR_CURSOR)
			// case types.CELL_CRUISER:
			// 	ch = types.CRUISER_SPRITE[(utils.AbsInt(y-x)-1)%len(types.CRUISER_SPRITE)]
			// 	m.win.ColorOn(types.COLOR_SHIP)
			// case types.CELL_DESTROYER:
			// 	ch = types.CRUISER_SPRITE[(utils.AbsInt(y-x)-1)%len(types.DESTROYER_SPRITE)]
			// 	m.win.ColorOn(types.COLOR_SHIP)
			// case types.CELL_BATTLESHIP:
			// 	ch = types.BATTLESHIP_SPRITE[(utils.AbsInt(y-x)-1)%len(types.BATTLESHIP_SPRITE)]
			// 	m.win.ColorOn(types.COLOR_SHIP)
			// case types.CELL_CARRIER:
			// 	ch = types.CARRIER_SPRITE[(utils.AbsInt(y-x)-1)%len(types.CARRIER_SPRITE)]
			// 	m.win.ColorOn(types.COLOR_SHIP)
			// case types.CELL_SUBMARINE:
			// 	ch = types.SUBMARINE_SPRITE[(utils.AbsInt(y-x)-1)%len(types.SUBMARINE_SPRITE)]
			// 	m.win.ColorOn(types.COLOR_SHIP)
			// case types.CELL_DESTROYED:
			// 	ch = '#'
			// 	m.win.ColorOn(types.COLOR_FLAMES)
			// case types.CELL_MISS:
			// 	ch = 'o'
			// 	m.win.ColorOn(types.COLOR_MISS)
			// default:
			// 	wallch, ok := types.WALLS_ASCII[cell]
			// 	if ok {
			// 		ch = wallch
			// 	} else {
			// 		ch = '|'
			// 	}
			// 	m.win.ColorOn(types.COLOR_WALL)
			// }

			// m.win.MoveAddChar(y*yfactor+3, x*xfactor, goncurses.Char(ch))

			m.win.ColorOn(cell.color)
			m.win.MovePrint(y*yfactor+3, x*xfactor, cell.content)
			m.win.ColorOff(cell.color)
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

	maxSize := len(entity.Sprite[o])
	s := 0
	switch o {
	case types.VERTICAL:
		for y := entity.StartPosition.Y; y <= entity.EndPosition.Y; y++ {
			if s > maxSize {
				panic(fmt.Sprintf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize))
			}
			(*m.grid)[types.Position{X: entity.StartPosition.X, Y: y}] = Cell{
				content:  fmt.Sprintf("  %c", entity.Sprite[o][s]),
				cellType: entity.CellType,
				color:    entity.Color,
			}
			s++
		}
	case types.HORIZONTAL:
		for x := entity.StartPosition.X; x <= utils.ExpectedEndPosition(entity.StartPosition.X, entity.Sprite[o]); x++ {
			if s > maxSize {
				panic(fmt.Sprintf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize))
			}
			(*m.grid)[types.Position{X: x, Y: entity.StartPosition.Y}] = Cell{
				content:  fmt.Sprintf("  %c", entity.Sprite[o][s]),
				cellType: entity.CellType,
				color:    entity.Color,
			}
			s++
		}
	}
}

func (m *Map) createEmptyGrid() *map[types.Position]Cell {
	grid := make(map[types.Position]Cell)

	for x := 1; x <= m.gridWidth; x++ {
		for y := 1; y <= m.gridHeight; y++ {
			grid[types.Position{X: x, Y: y}] = Cell{
				content:  "   ",
				cellType: types.CELL_WATER,
				color:    types.COLOR_WATER,
			}
		}
	}

	createBorders(&grid, m.gridWidth, m.gridHeight)

	return &grid
}

func createBorders(grid *map[types.Position]Cell, gridWidth, gridHeight int) {
	for x := 1; x <= gridWidth; x++ {
		for y := 1; y <= gridHeight; y++ {
			if x == 1 && y == 1 {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TOP_LEFT]),
					cellType: types.CELL_WALL_TOP_LEFT,
					color:    types.COLOR_WALL,
				}
			} else if x == 1 && y == gridHeight {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_BOTTOM_LEFT]),
					cellType: types.CELL_WALL_BOTTOM_LEFT,
					color:    types.COLOR_WALL,
				}
			} else if x == gridWidth && y == 1 {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TOP_RIGHT]),
					cellType: types.CELL_WALL_TOP_RIGHT,
					color:    types.COLOR_WALL,
				}
			} else if x == gridWidth && y == gridHeight {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_BOTTOM_RIGHT]),
					cellType: types.CELL_WALL_BOTTOM_RIGHT,
					color:    types.COLOR_WALL,
				}
			} else if x == 1 {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TEE_DOWN]),
					cellType: types.CELL_WALL_TEE_DOWN,
					color:    types.COLOR_WALL,
				}
			} else if x == gridWidth {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TEE_UP]),
					cellType: types.CELL_WALL_TEE_UP,
					color:    types.COLOR_WALL,
				}
			} else if y == 1 {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TEE_RIGHT]),
					cellType: types.CELL_WALL_TEE_RIGHT,
					color:    types.COLOR_WALL,
				}
			} else if y == gridHeight {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_TEE_LEFT]),
					cellType: types.CELL_WALL_TEE_LEFT,
					color:    types.COLOR_WALL,
				}
			} else {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_CORNER]),
					cellType: types.CELL_WALL_CORNER,
					color:    types.COLOR_WALL,
				}
			}

			// draw outer borders
			// outer borders must be drawn when x
			//
			// scratch all that
			// since borders are constant why should i calculate it or even place it in the grid?
			// why not just draw the border and only store variable cells in the grid?
			if x < gridWidth {
				for i := 1; i < 4; i++ {
					// win.MoveAddChar(y, x+i, horizontal)
					(*grid)[types.Position{X: i, Y: y}] = Cell{
						content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_HORIZONTAL]),
						cellType: types.CELL_WALL_HORIZONTAL,
						color:    types.COLOR_WALL,
					}
				}
			}

			if y < gridHeight {
				(*grid)[types.Position{X: x, Y: y}] = Cell{
					content:  fmt.Sprintf(" %c", types.WALLS_ASCII[types.CELL_WALL_VERTICAL]),
					cellType: types.CELL_WALL_VERTICAL,
					color:    types.COLOR_WALL,
				}
			}
			// win.ColorOff(WHITE_BLACK)
		}
	}
}
