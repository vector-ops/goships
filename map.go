package main

import (
	"context"
	"fmt"
	"math/rand"
	"slices"
	"sync/atomic"

	"github.com/rthornton128/goncurses"
	"github.com/vector-ops/goships/logger"
	"github.com/vector-ops/goships/types"
	"github.com/vector-ops/goships/utils"
)

const (
	DefaultGridWidth  = 10
	DefaultGridHeight = 7
	CellHeight        = 2
	CellWidth         = 4
	DefaultTotalShips = 5
)

// Events
const (
	Win = iota
	ShipsPlaced
	Turn
	Hit
	Miss
)

type Map struct {
	win   *goncurses.Window
	debug bool

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
	stats         *MapStats
	turn          int64

	logEvents []int

	closeCh chan struct{}
	logger  *logger.Logger
}

type MapStats struct {
	Hits   int
	Misses int
	Ships  map[types.ShipType]*ShipStatus
}

func (ms *MapStats) GetShipsDestroyed() int {
	shipsDestroyed := 0
	for _, ship := range ms.Ships {
		if ship.destroyed {
			shipsDestroyed++
		}
	}
	return shipsDestroyed
}

type Cursor struct {
	startPosition types.Position
	endPosition   types.Position
	orientation   types.Orientation
	content       []rune
	shipType      *types.ShipType
}

type ShipStatus struct {
	totalCells int
	hitCells   int
	destroyed  bool
}

func NewMap(win *goncurses.Window, isPlayerMap bool, title string, titleColor int16, startingGrid *map[types.Position]types.Cell, gridWidth, gridHeight *int, enableKeyboard bool, debug bool, l *logger.Logger) *Map {
	m := &Map{
		win:            win,
		debug:          debug,
		title:          title,
		titleColor:     titleColor,
		gridHeight:     DefaultGridHeight,
		gridWidth:      DefaultGridWidth,
		cursor:         &Cursor{startPosition: types.Position{X: 0, Y: 0}, endPosition: types.Position{X: 0, Y: 0}, orientation: types.Horizontal, content: []rune{' '}},
		enableCursor:   false,
		enableKeyboard: enableKeyboard,
		isPlayerMap:    isPlayerMap,
		unplacedShips:  []types.ShipType{types.AircraftCarrier, types.Battleship, types.Cruiser, types.Destroyer, types.Submarine},
		logger:         l,
		closeCh:        make(chan struct{}),
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

	ships := make(map[types.ShipType]*ShipStatus)
	for _, shipType := range m.unplacedShips {
		ships[shipType] = &ShipStatus{
			totalCells: len(utils.GetEntitySprite(shipType)),
			hitCells:   0,
			destroyed:  false,
		}
	}
	mapStats := &MapStats{
		Ships: ships,
	}

	m.stats = mapStats

	go m.eventLogger()

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

func (m *Map) GetStats() *MapStats {
	return m.stats
}

func (m *Map) GetTurn() int64 {
	return m.turn
}

func (m *Map) HandleKeyInput(key goncurses.Key) {
	if !m.enableKeyboard || !m.enableCursor {
		return
	}

	if m.isPlayerMap && len(m.unplacedShips) == 0 {
		return
	}

	switch key {
	case 'u':
		m.SaveState()
	case goncurses.KEY_UP, 'w':
		if m.cursor.startPosition.Y > 0 {
			m.cursor.startPosition.Y--
			m.cursor.endPosition.Y--
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_DOWN, 's':
		if (m.cursor.orientation == types.Horizontal && m.cursor.startPosition.Y < m.gridHeight-1) || (m.cursor.orientation == types.Vertical && m.cursor.endPosition.Y < m.gridHeight-1) {
			m.cursor.startPosition.Y++
			m.cursor.endPosition.Y++
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_LEFT, 'a':
		if m.cursor.startPosition.X > 0 {
			m.cursor.startPosition.X--
			m.cursor.endPosition.X--
			if !m.isPlayerMap {
				m.cursor.content = []rune{(*m.grid)[m.cursor.startPosition].Content}
			}
		}
	case goncurses.KEY_RIGHT, 'd':
		if (m.cursor.orientation == types.Horizontal && m.cursor.endPosition.X < m.gridWidth-1) || (m.cursor.orientation == types.Vertical && m.cursor.startPosition.X < m.gridWidth-1) {
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
		if newOrientation == types.Horizontal {
			newOrientation = types.Vertical
		} else {
			newOrientation = types.Horizontal
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
				Color:         types.ColorShip,
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
			// cell := (*m.grid)[m.cursor.startPosition]

			if hit := m.hitCell(m.cursor.startPosition.X, m.cursor.startPosition.Y); hit {
				m.turn++
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
	close(m.closeCh)
	return m.win.Delete()
}

func (m *Map) SaveState() {
	unplacedShips := []string{}
	for _, ship := range m.unplacedShips {
		unplacedShips = append(unplacedShips, utils.GetShipType(ship))
	}

	utils.SaveMapState(m.title, *m.grid, unplacedShips)
}

func (m *Map) PlaceRandomShips() error {
	if len(m.unplacedShips) == 0 {
		return fmt.Errorf("no ships to place")
	}

	for _, shipType := range m.unplacedShips {

		var orientation types.Orientation

		sprite := utils.GetEntitySprite(shipType)

		x := rand.Intn(m.gridWidth)
		y := rand.Intn(m.gridHeight)

		startPosition := types.Position{X: x, Y: y}

		endPosition := utils.ExpectedEndPosition(startPosition, sprite, orientation)

		for utils.CheckOverlap((*m.grid), types.Ship{StartPosition: startPosition, EndPosition: endPosition}) || !utils.ValidEntityPosition(types.Ship{StartPosition: startPosition, EndPosition: endPosition}, m.gridHeight, m.gridWidth) {
			x := rand.Intn(m.gridWidth)
			y := rand.Intn(m.gridHeight)

			randOrientation := rand.Intn(2)
			if randOrientation == 0 {
				orientation = types.Horizontal
			} else {
				orientation = types.Vertical
			}

			startPosition = types.Position{X: x, Y: y}
			endPosition = utils.ExpectedEndPosition(startPosition, sprite, orientation)

		}

		ship := types.Ship{
			Type:          shipType,
			StartPosition: startPosition,
			EndPosition:   endPosition,
			Color:         types.ColorShip,
		}

		err := m.placeShip(ship, orientation)
		if err != nil {
			return err
		}

	}
	m.unplacedShips = nil
	return nil
}

func (m *Map) HitRandomSpot() {
	if !m.hasEmptyCells() {
		return
	}

	hit := false

	for !hit {

		x := rand.Intn(m.gridWidth)
		y := rand.Intn(m.gridHeight)

		hit = m.hitCell(x, y)
	}
	m.turn++
}

func (m *Map) hitCell(x, y int) bool {
	cell := (*m.grid)[types.Position{X: x, Y: y}]

	if cell.Type == types.CellDestroyed || cell.Type == types.CellMiss {
		return false
	}

	if cell.Type == types.CellWater || cell.Type == types.CellCursor {
		(*m.grid)[types.Position{X: x, Y: y}] = types.Cell{
			Type:    types.CellMiss,
			Color:   types.ColorMiss,
			Content: 'x',
			Hit:     true,
		}

		m.stats.Misses++
	}
	if cell.Type == types.CellShip {
		(*m.grid)[types.Position{X: x, Y: y}] = types.Cell{
			ShipType: cell.ShipType,
			Type:     types.CellDestroyed,
			Color:    types.ColorHit,
			Content:  cell.Content,
			Hit:      true,
		}
		m.stats.Ships[cell.ShipType].hitCells++
		if m.stats.Ships[cell.ShipType].hitCells == m.stats.Ships[cell.ShipType].totalCells {
			m.stats.Ships[cell.ShipType].destroyed = true
		}
		m.stats.Hits++
	}

	return true
}

func (m *Map) draw() error {
	my, mx := m.win.MaxYX()
	offsetX := CellWidth/2 - 1
	offsetY := CellHeight / 2

	startX := (mx - m.gridWidth*CellWidth) / 2
	startY := (my - m.gridHeight*CellHeight) / 2

	for col := 0; col < m.gridWidth; col++ {
		for row := 0; row < m.gridHeight; row++ {

			cell := (*m.grid)[types.Position{X: col, Y: row}]
			if m.enableCursor && col >= m.cursor.startPosition.X && col <= m.cursor.endPosition.X && row >= m.cursor.startPosition.Y && row <= m.cursor.endPosition.Y {

				var relativeIndex int
				if m.cursor.orientation == types.Horizontal {
					relativeIndex = col - m.cursor.startPosition.X
				} else {
					relativeIndex = row - m.cursor.startPosition.Y
				}

				color := types.ColorCursor

				if utils.CheckOverlap((*m.grid), types.Ship{StartPosition: m.cursor.startPosition, EndPosition: m.cursor.endPosition}) && m.isPlayerMap {
					color = types.BlackRed
				}

				hit := cell.Hit
				cell = types.Cell{
					Type:    types.CellCursor,
					Color:   color,
					Content: m.cursor.content[relativeIndex],
					Hit:     hit,
				}
			}

			x := (startX + offsetX) + col*CellWidth
			y := (startY + offsetY) + row*CellHeight

			color := cell.Color
			content := cell.Content
			if !m.isPlayerMap && cell.Type == types.CellShip {
				color = types.ColorWater
				content = ' '
			}

			if !m.isPlayerMap && (cell.Type == types.CellCursor || cell.Type == types.CellShip) && !cell.Hit {
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
	if !utils.ValidEntityPosition(entity, m.gridHeight, m.gridWidth) {
		return fmt.Errorf("invalid entity position: %s, start: %d,%d", utils.GetShipType(entity.Type), entity.StartPosition.X, entity.StartPosition.Y)
	}

	sprite := utils.GetEntitySprite(entity.Type)
	maxSize := len(sprite)
	s := 0
	switch o {
	case types.Vertical:
		for y := entity.StartPosition.Y; y <= utils.ExpectedEndCoordinate(entity.StartPosition.Y, sprite); y++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: entity.StartPosition.X, Y: y}] = types.Cell{
				Content:  sprite[s],
				Type:     types.CellShip,
				Color:    entity.Color,
				ShipType: entity.Type,
			}
			s++
		}
	case types.Horizontal:
		for x := entity.StartPosition.X; x <= utils.ExpectedEndCoordinate(entity.StartPosition.X, sprite); x++ {
			if s > maxSize {
				return fmt.Errorf("sprite smaller than entity size, s: %d, entity size: %d", s, maxSize)
			}
			(*m.grid)[types.Position{X: x, Y: entity.StartPosition.Y}] = types.Cell{
				Content:  sprite[s],
				Type:     types.CellShip,
				Color:    entity.Color,
				ShipType: entity.Type,
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
				Content:  ' ',
				Type:     types.CellWater,
				Color:    types.ColorWater,
				ShipType: types.None,
			}
		}
	}

	return &grid
}

func (m *Map) drawBorders() error {
	my, mx := m.win.MaxYX()

	if mx < m.gridWidth*CellWidth || my < m.gridHeight*CellHeight {
		return fmt.Errorf("window size is too small to accommodate the map. Map is %dx%d, window is %dx%d", m.gridWidth*CellWidth, m.gridHeight*CellHeight, mx, my)
	}

	// calculate the starting position for the drawBorders
	// borders start is (max size - total size of all types.Cells) / 2
	// this gives the starting position of the grid such that it is centered in the window
	startX := (mx - m.gridWidth*CellWidth) / 2
	startY := (my - m.gridHeight*CellHeight) / 2

	for row := 0; row <= m.gridHeight; row++ {
		for col := 0; col <= m.gridWidth; col++ {
			y := startY + row*CellHeight
			x := startX + col*CellWidth

			m.win.ColorOn(types.ColorWall)
			m.win.MoveAddChar(y, x, goncurses.Char(types.WallsASCII[types.CellWallCorner]))

			if col < m.gridWidth {
				for i := 1; i < CellWidth; i++ {
					m.win.MoveAddChar(y, x+i, goncurses.Char(types.WallsASCII[types.CellWallHorizontal]))
				}
			}

			if row < m.gridHeight {
				m.win.MoveAddChar(y+1, x, goncurses.Char(types.WallsASCII[types.CellWallVertical]))
			}
			m.win.ColorOff(types.ColorWall)
		}
	}

	return nil
}

func (m *Map) hasEmptyCells() bool {
	for _, cell := range *m.grid {
		if cell.Type == types.CellBlank || cell.Type == types.CellWater {
			return true
		}
	}
	return false
}

func (m *Map) eventLogger() {

	loggedEvents := make([]int, 0)

	prevVals := make(map[string]any)

	for {
		select {
		case <-m.closeCh:
			return
		default:
			for _, event := range m.logEvents {
				switch event {
				case ShipsPlaced:
					if len(m.unplacedShips) == 0 && !slices.Contains(loggedEvents, ShipsPlaced) {
						msg := fmt.Sprintf("%s has placed ships", m.title)
						m.logger.Infof(msg)
						loggedEvents = append(loggedEvents, ShipsPlaced)
					}
				case Turn:

					prevTurn, ok := prevVals["turn"].(int64)
					if !ok {
						prevTurn = 0
					}

					turn :=
						atomic.LoadInt64(&(m.turn))
					if prevTurn != turn {
						msg := fmt.Sprintf("%s's turn", m.title)
						m.logger.Infof(msg)
						prevVals["turn"] = turn
					}

				case Hit:
					prevHits, ok := prevVals["hits"].(int)
					if !ok {
						prevHits = 0
					}

					hits := m.stats.Hits

					if prevHits < hits {
						msg := fmt.Sprintf("%s's ship was hit", m.title)
						m.logger.Infof(msg)
						prevVals["hits"] = hits
					}

				}

			}
		}
	}
}

func (m *Map) LogOn(event int) {
	m.logEvents = append(m.logEvents, event)
}
