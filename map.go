package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

const (
	DEFAULT_GRID_WIDTH  = 10
	DEFAULT_GRID_HEIGHT = 7
	CELL_HEIGHT         = 2
	CELL_WIDTH          = 4
)

type Map struct {
	win *goncurses.Window

	title          string
	titleColor     int16
	grid           *map[types.Position]types.Cell
	gridHeight     int
	gridWidth      int
	cursor         *Cursor
	enableCursor   bool
	enableKeyboard bool
	isPlayerMap    bool

	unplacedShips []types.ShipType
}

type Cursor struct {
	startPosition types.Position
	endPosition   types.Position
	orientation   types.Orientation
	content       []rune
	shipType      *types.ShipType
}

func NewMap(win *goncurses.Window, isPlayerMap bool, title string, titleColor int16, startingGrid *map[types.Position]types.Cell, gridWidth, gridHeight *int, enableKeyboard bool) *Map {
	m := &Map{
		win:            win,
		title:          title,
		titleColor:     titleColor,
		gridHeight:     DEFAULT_GRID_HEIGHT,
		gridWidth:      DEFAULT_GRID_WIDTH,
		cursor:         &Cursor{startPosition: types.Position{X: 0, Y: 0}, endPosition: types.Position{X: 0, Y: 0}, orientation: types.HORIZONTAL, content: []rune{' '}},
		enableCursor:   false,
		enableKeyboard: enableKeyboard,
		isPlayerMap:    isPlayerMap,
		unplacedShips:  []types.ShipType{types.AIRCRAFT_CARRIER, types.BATTLESHIP, types.CRUISER, types.DESTROYER, types.SUBMARINE},
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

func (m *Map) EnableCursor(enable bool) {
	if enable && m.cursor == nil {
		m.cursor = &Cursor{startPosition: types.Position{X: 0, Y: 0}, endPosition: types.Position{X: 0, Y: 0}, content: []rune{' '}}
	}
	if m.isPlayerMap && len(m.unplacedShips) > 0 && m.cursor.shipType == nil {
		m.cursor.shipType = &m.unplacedShips[0]
		m.cursor.endPosition = utils.ExpectedEndPosition(m.cursor.startPosition, utils.GetEntitySprite(*m.cursor.shipType), m.cursor.orientation)
		m.cursor.content = utils.GetEntitySprite(*m.cursor.shipType)
	}
	m.enableCursor = enable
}

func (m *Map) HasPlacedShips() bool {
	return len(m.unplacedShips) == 0
}

func (m *Map) HandleKeyInput(key goncurses.Key) {
	if !m.enableKeyboard || !m.enableCursor {
		return
	}

	if m.isPlayerMap && len(m.unplacedShips) == 0 {
		return
	}

	switch key {
	case goncurses.KEY_UP:
		if m.cursor.startPosition.Y > 0 {
			m.cursor.startPosition.Y--
			m.cursor.endPosition.Y--
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_DOWN:
		if (m.cursor.orientation == types.HORIZONTAL && m.cursor.startPosition.Y < m.gridHeight-1) || (m.cursor.orientation == types.VERTICAL && m.cursor.endPosition.Y < m.gridHeight-1) {
			m.cursor.startPosition.Y++
			m.cursor.endPosition.Y++
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_LEFT:
		if m.cursor.startPosition.X > 0 {
			m.cursor.startPosition.X--
			m.cursor.endPosition.X--
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_RIGHT:
		if (m.cursor.orientation == types.HORIZONTAL && m.cursor.endPosition.X < m.gridWidth-1) || (m.cursor.orientation == types.VERTICAL && m.cursor.startPosition.X < m.gridWidth-1) {
			m.cursor.startPosition.X++
			m.cursor.endPosition.X++
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}

	case ' ':
		if !m.isPlayerMap {
			return
		}
		newOrientation := m.cursor.orientation
		if newOrientation == types.HORIZONTAL {
			newOrientation = types.VERTICAL
		} else {
			newOrientation = types.HORIZONTAL
		}

		sprite := m.cursor.content

		expectedEndPosition := utils.ExpectedEndPosition(m.cursor.startPosition, sprite, newOrientation)

		if expectedEndPosition.X < m.gridWidth && expectedEndPosition.Y < m.gridHeight {
			m.cursor.endPosition = expectedEndPosition
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
			m.cursor.orientation = newOrientation
		}

	case goncurses.KEY_ENTER, goncurses.KEY_RETURN:
		if m.isPlayerMap {
			if m.cursor.shipType == nil {
				return
			}

			if utils.CheckOverlap((*m.grid), types.Ship{StartPosition: m.cursor.startPosition, EndPosition: m.cursor.endPosition}) {
				return
			}

			entity := types.Ship{
				Type:          *m.cursor.shipType,
				StartPosition: m.cursor.startPosition,
				EndPosition:   m.cursor.endPosition,
				Color:         types.COLOR_SHIP,
			}

			m.placeShip(entity, m.cursor.orientation)

			m.unplacedShips = m.unplacedShips[1:]
			if len(m.unplacedShips) == 0 {
				m.cursor.shipType = nil
			} else {
				m.cursor.shipType = &m.unplacedShips[0]
				m.cursor.endPosition = utils.ExpectedEndPosition(m.cursor.startPosition, utils.GetEntitySprite(*m.cursor.shipType), m.cursor.orientation)
				m.cursor.content = utils.GetEntitySprite(*m.cursor.shipType)
			}
		} else {

			cell := (*m.grid)[m.cursor.startPosition]
			if cell.Type == types.CELL_WATER || cell.Type == types.CELL_CURSOR {
				(*m.grid)[m.cursor.startPosition] = types.Cell{
					Type:    types.CELL_MISS,
					Color:   types.COLOR_MISS,
					Content: 'x',
				}
			}
			if cell.Type == types.CELL_SHIP {
				(*m.grid)[m.cursor.startPosition] = types.Cell{
					Type:    types.CELL_DESTROYED,
					Color:   types.COLOR_HIT,
					Content: cell.Content,
				}
			}
		}

	default:

	}
}

func (m *Map) Render(ctx context.Context) error {
	m.win.Erase()
	err := m.win.Box(goncurses.ACS_VLINE, goncurses.ACS_HLINE)
	if err != nil {
		return err
	}

	_, mx := m.win.MaxYX()

	m.win.ColorOn(m.titleColor)
	m.win.MovePrint(1, (mx/2)-len(m.title)/2, m.title)
	m.win.ColorOff(m.titleColor)

	err = m.drawBorders()
	if err != nil {
		return err
	}
	m.draw()
	m.win.NoutRefresh()
	return nil
}

func (m *Map) Close() error {
	return m.win.Delete()
}

func (m *Map) PlaceRandomShips() error {
	if len(m.unplacedShips) == 0 {
		return nil
	}
	for _, shipType := range m.unplacedShips {

		var orientation types.Orientation

		sprite := utils.GetEntitySprite(shipType)

		var startPosition types.Position
		var endPosition types.Position

		for utils.CheckOverlap((*m.grid), types.Ship{StartPosition: startPosition, EndPosition: endPosition}) {
			x := rand.Intn(m.gridWidth)
			y := rand.Intn(m.gridHeight)

			randOrientation := rand.Intn(2)
			if randOrientation == 0 {
				orientation = types.HORIZONTAL
			} else {
				orientation = types.VERTICAL
			}

			startPosition = types.Position{X: x, Y: y}
			endPosition = utils.ExpectedEndPosition(startPosition, sprite, orientation)
		}

		ship := types.Ship{
			Type:          shipType,
			StartPosition: startPosition,
			EndPosition:   endPosition,
			Color:         types.COLOR_SHIP,
		}

		err := m.placeShip(ship, orientation)
		if err != nil {
			return err
		}

	}
	m.unplacedShips = nil
	return nil
}

func (m *Map) draw() error {
	my, mx := m.win.MaxYX()
	offsetX := CELL_WIDTH/2 - 1
	offsetY := CELL_HEIGHT / 2

	startX := (mx - m.gridWidth*CELL_WIDTH) / 2
	startY := (my - m.gridHeight*CELL_HEIGHT) / 2

	for col := 0; col < m.gridWidth; col++ {
		for row := 0; row < m.gridHeight; row++ {

			cell := (*m.grid)[types.Position{X: col, Y: row}]
			if m.enableCursor && col >= m.cursor.startPosition.X && col <= m.cursor.endPosition.X && row >= m.cursor.startPosition.Y && row <= m.cursor.endPosition.Y {

				var relativeIndex int
				if m.cursor.orientation == types.HORIZONTAL {
					relativeIndex = col - m.cursor.startPosition.X
				} else {
					relativeIndex = row - m.cursor.startPosition.Y
				}

				color := types.COLOR_CURSOR

				if utils.CheckOverlap((*m.grid), types.Ship{StartPosition: m.cursor.startPosition, EndPosition: m.cursor.endPosition}) && m.isPlayerMap {
					color = types.BLACK_RED
				}

				cell = types.Cell{
					Type:    types.CELL_CURSOR,
					Color:   color,
					Content: m.cursor.content[relativeIndex],
				}
			}

			x := (startX + offsetX) + col*CELL_WIDTH
			y := (startY + offsetY) + row*CELL_HEIGHT

			color := cell.Color
			content := cell.Content
			if !m.isPlayerMap && cell.Type == types.CELL_SHIP {
				color = types.COLOR_WATER
				content = ' '
			}

			if !m.isPlayerMap && cell.Type == types.CELL_CURSOR {
				content = ' '
			}

			m.win.ColorOn(color)
			m.win.MovePrint(y, x, fmt.Sprintf(" %c ", content))
			m.win.ColorOff(color)
		}
	}
	return nil
}

func (m *Map) placeShip(entity types.Ship, o types.Orientation) error {
	if !utils.ValidateEntityPosition(entity, m.gridHeight, m.gridWidth) {
		return errors.New("Invalid entity position")
	}

	sprite := utils.GetEntitySprite(entity.Type)
	maxSize := len(sprite)
	s := 0
	switch o {
	case types.VERTICAL:
		for y := entity.StartPosition.Y; y <= utils.ExpectedEndCoordinate(entity.StartPosition.Y, sprite); y++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: entity.StartPosition.X, Y: y}] = types.Cell{
				Content: sprite[s],
				Type:    types.CELL_SHIP,
				Color:   entity.Color,
			}
			s++
		}
	case types.HORIZONTAL:
		for x := entity.StartPosition.X; x <= utils.ExpectedEndCoordinate(entity.StartPosition.X, sprite); x++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: x, Y: entity.StartPosition.Y}] = types.Cell{
				Content: sprite[s],
				Type:    types.CELL_SHIP,
				Color:   entity.Color,
			}
			s++
		}
	}

	return nil
}

func (m *Map) createEmptyGrid() *map[types.Position]types.Cell {
	grid := make(map[types.Position]types.Cell)

	for x := 0; x < m.gridWidth; x++ {
		for y := 0; y < m.gridHeight; y++ {
			grid[types.Position{X: x, Y: y}] = types.Cell{
				Content: ' ',
				Type:    types.CELL_WATER,
				Color:   types.COLOR_WATER,
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
	// borders start is (max size - total size of all types.Cells) / 2
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
