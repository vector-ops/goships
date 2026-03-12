package types

type GameType string

const (
	PVE        GameType = "PVE"
	PVP        GameType = "PVP"
	UNSELECTED GameType = "unselected"
	QUIT       GameType = "quit"
)

type ShipType int

func (st ShipType) String() string {
	switch st {
	case 0:
		return "none"
	case 1:
		return "cruiser"
	case 2:
		return "destroyer"
	case 3:
		return "submarine"
	case 4:
		return "aircraft_carrier"
	case 5:
		return "battleship"
	default:
		return "none"
	}
}

const (
	BATTLESHIP ShipType = iota
	CRUISER
	DESTROYER
	SUBMARINE
	AIRCRAFT_CARRIER
	NONE
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
	GUIDE
	MENU
)

type CellType int

func (ct CellType) String() string {
	switch ct {
	case 0:
		return "cursor"
	case 1:
		return "ship"
	case 2:
		return "destroyed"
	case 3:
		return "miss"
	case 4:
		return "water"
	case 5:
		return "blank"
	default:
		return "wall"
	}
}

const (
	CELL_CURSOR CellType = iota
	CELL_SHIP
	CELL_DESTROYED
	CELL_MISS
	CELL_WATER
	CELL_BLANK
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
	COLOR_HIT
	COLOR_MISS
	COLOR_CURSOR
	WHITE_BLACK
	RED_BLACK
	GREEN_BLACK
	BLUE_BLACK
	YELLOW_BLACK
	MAGENTA_BLACK
	CYAN_BLACK
	WHITE_RED
	WHITE_GREEN
	WHITE_BLUE
	BLACK_YELLOW
	BLACK_CYAN
	BLACK_MAGENTA
	BLACK_RED
)

const (
	MAINWINDOW = iota + 1
)

type Position struct {
	X int `json:"x"`
	Y int `json:"y"`
}

type Ship struct {
	Type          ShipType
	StartPosition Position
	EndPosition   Position
	Color         int16
}

var (
	CRUISER_SPRITE    = []rune{'C', 'R', 'U', 'Z'}
	BATTLESHIP_SPRITE = []rune{'B', 'B', 'S', 'H', 'P'}
	DESTROYER_SPRITE  = []rune{'D', 'D', 'S'}
	SUBMARINE_SPRITE  = []rune{'S', 'S'}
	CARRIER_SPRITE    = []rune{'C', 'V', 'S', 'H', 'P'}
)

type Cell struct {
	Type     CellType `json:"type"`
	ShipType ShipType `json:"ship_type"`
	Color    int16    `json:"color"`
	Content  rune     `json:"content"`
	Hit      bool     `json:"hit"`
}
