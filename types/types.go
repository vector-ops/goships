package types

type GameType string

const (
	PvE        GameType = "PVE"
	PvP        GameType = "PVP"
	Unselected GameType = "unselected"
	Quit       GameType = "quit"
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
	None ShipType = iota
	Cruiser
	Destroyer
	Submarine
	AircraftCarrier
	Battleship
)

type Orientation int

const (
	Vertical Orientation = iota
	Horizontal
)

type WindowType int

const (
	Player WindowType = iota
	Enemy
	Score
	Guide
	Menu
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
	CellCursor CellType = iota
	CellShip
	CellDestroyed
	CellMiss
	CellWater
	CellBlank
	CellWallHorizontal
	CellWallVertical
	CellWallCorner
	CellWallTopLeft
	CellWallTopRight
	CellWallBottomLeft
	CellWallBottomRight
	CellWallTeeLeft
	CellWallTeeRight
	CellWallTeeUp
	CellWallTeeDown
)

var WallsBox = map[CellType]rune{
	CellWallHorizontal:  0x2500, // ─
	CellWallVertical:    0x2502, // │
	CellWallCorner:      0x253C, // ┼
	CellWallTopLeft:     0x250C, // ┌
	CellWallTopRight:    0x2510, // ┐
	CellWallBottomLeft:  0x2514, // └
	CellWallBottomRight: 0x2518, // ┘
	CellWallTeeLeft:     0x2524, // ┬
	CellWallTeeRight:    0x251C, // ┴
	CellWallTeeUp:       0x2534, // ├
	CellWallTeeDown:     0x252C, // ┤
}

var WallsASCII = map[CellType]rune{
	CellWallHorizontal:  '-', // ─
	CellWallVertical:    '|', // │
	CellWallCorner:      '+', // ┼
	CellWallTopLeft:     '+', // ┌
	CellWallTopRight:    '+', // ┐
	CellWallBottomLeft:  '+', // └
	CellWallBottomRight: '+', // ┘
	CellWallTeeLeft:     '+', // ┬
	CellWallTeeRight:    '+', // ┴
	CellWallTeeUp:       '+', // ├
	CellWallTeeDown:     '+', // ┤
}

const (
	ColorWater int16 = iota + 1
	ColorWall
	ColorShip
	ColorHit
	ColorMiss
	ColorCursor
	WhiteBlack
	RedBlack
	GreenBlack
	BlueBlack
	YellowBlack
	MagentaBlack
	CyanBlack
	WhiteRed
	WhiteGreen
	WhiteBlue
	BlackYellow
	BlackCyan
	BlackMagenta
	BlackRed
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
	CruiserSprite    = []rune{'C', 'R', 'U', 'Z'}
	BattleshipSprite = []rune{'B', 'B', 'S', 'H', 'P'}
	DestroyerSprite  = []rune{'D', 'D', 'S'}
	SubmarineSprite  = []rune{'S', 'S'}
	CarrierSprite    = []rune{'C', 'V', 'S', 'H', 'P'}
)

type Cell struct {
	Type     CellType `json:"type"`
	ShipType ShipType `json:"ship_type"`
	Color    int16    `json:"color"`
	Content  rune     `json:"content"`
	Hit      bool     `json:"hit"`
}
