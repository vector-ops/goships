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
	CELL_WALL
	CELL_CURSOR
	CELL_CRUISER
	CELL_DESTROYER
	CELL_BATTLESHIP
	CELL_CARRIER
	CELL_SUBMARINE
	CELL_DESTROYED
	CELL_MISS
)

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
	StartPosition Position
	EndPosition   Position
	Color         int16
	Sprite        map[Orientation][]rune
}
