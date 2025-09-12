package types

type GameType string

const (
	PVE        GameType = "PVE"
	PVP        GameType = "PVP"
	UNSELECTED GameType = "unselected"
	QUIT       GameType = "quit"
)

type ShipType int

const (
	BATTLESHIP ShipType = iota
	CRUISER
	DESTROYER
	SUBMARINE
	AIRCRAFT_CARRIER
)

type Orientation int

const (
	VERTICAL Orientation = iota
	HORIZONTAL
)

type WindowType int

const (
	PLAYER WindowType = iota
	ENEMY
	SCORE
	BUFFER
	MENU
)

type CellType int

const (
	CELL_BLANK CellType = iota
	CELL_WALL_HORIZONTAL
	CELL_WALL_VERTICAL
	CELL_WALL_CORNER
	CELL_WALL_TOP_LEFT
	CELL_WALL_TOP_RIGHT
	CELL_WALL_BOTTOM_LEFT
	CELL_WALL_BOTTOM_RIGHT
	CELL_WALL_TEE_LEFT
	CELL_WALL_TEE_RIGHT
	CELL_WALL_TEE_UP
	CELL_WALL_TEE_DOWN
	CELL_CURSOR
	CELL_CRUISER
	CELL_DESTROYER
	CELL_BATTLESHIP
	CELL_CARRIER
	CELL_SUBMARINE
	CELL_DESTROYED
	CELL_MISS
	CELL_WATER
)

var WALLS_BOX = map[CellType]rune{
	CELL_WALL_HORIZONTAL:   0x2500, // ─
	CELL_WALL_VERTICAL:     0x2502, // │
	CELL_WALL_CORNER:       0x253C, // ┼
	CELL_WALL_TOP_LEFT:     0x250C, // ┌
	CELL_WALL_TOP_RIGHT:    0x2510, // ┐
	CELL_WALL_BOTTOM_LEFT:  0x2514, // └
	CELL_WALL_BOTTOM_RIGHT: 0x2518, // ┘
	CELL_WALL_TEE_LEFT:     0x2524, // ┬
	CELL_WALL_TEE_RIGHT:    0x251C, // ┴
	CELL_WALL_TEE_UP:       0x2534, // ├
	CELL_WALL_TEE_DOWN:     0x252C, // ┤
}

var WALLS_ASCII = map[CellType]rune{
	CELL_WALL_HORIZONTAL:   '-', // ─
	CELL_WALL_VERTICAL:     '|', // │
	CELL_WALL_CORNER:       '+', // ┼
	CELL_WALL_TOP_LEFT:     '+', // ┌
	CELL_WALL_TOP_RIGHT:    '+', // ┐
	CELL_WALL_BOTTOM_LEFT:  '+', // └
	CELL_WALL_BOTTOM_RIGHT: '+', // ┘
	CELL_WALL_TEE_LEFT:     '+', // ┬
	CELL_WALL_TEE_RIGHT:    '+', // ┴
	CELL_WALL_TEE_UP:       '+', // ├
	CELL_WALL_TEE_DOWN:     '+', // ┤
}

const (
	COLOR_WATER int16 = iota + 1
	COLOR_WALL
	COLOR_SHIP
	COLOR_FLAMES
	COLOR_MISS
	COLOR_CURSOR
)

const (
	MAINWINDOW = iota + 1
)

type Position struct {
	X, Y int
}

type Entity struct {
	Type          ShipType
	CellType      CellType
	StartPosition Position
	EndPosition   Position
	Color         int16
	Sprite        map[Orientation][]rune
}

var (
	CRUISER_SPRITE    = []rune{'C', 'R', 'U', 'Z'}
	BATTLESHIP_SPRITE = []rune{'B', 'B', 'S', 'H', 'P'}
	DESTROYER_SPRITE  = []rune{'D', 'D', 'S'}
	SUBMARINE_SPRITE  = []rune{'S', 'S'}
	CARRIER_SPRITE    = []rune{'C', 'V', 'S', 'H', 'P'}
)
