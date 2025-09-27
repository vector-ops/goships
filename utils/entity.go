package utils

import "github.com/vector-ops/goships/types"

func ValidEntityPosition(e types.Ship, gridHeight, gridWidth int) bool {
	return e.StartPosition.X <= gridWidth && e.StartPosition.Y <= gridHeight && e.EndPosition.X <= gridWidth && e.EndPosition.Y <= gridHeight && e.StartPosition.X >= 0 && e.StartPosition.Y >= 0 && e.EndPosition.X >= 0 && e.EndPosition.Y >= 0
}

func ExpectedEndPosition(position types.Position, sprite []rune, o types.Orientation) types.Position {
	if o == types.HORIZONTAL {
		return types.Position{
			X: position.X + len(sprite) - 1,
			Y: position.Y,
		}
	}
	return types.Position{
		X: position.X,
		Y: position.Y + len(sprite) - 1,
	}
}

func ExpectedEndCoordinate(start int, sprite []rune) int {
	return start + len(sprite) - 1
}

func GetEntitySprite(shipType types.ShipType) []rune {
	switch shipType {
	case types.BATTLESHIP:
		return types.BATTLESHIP_SPRITE
	case types.AIRCRAFT_CARRIER:
		return types.CARRIER_SPRITE
	case types.CRUISER:
		return types.CRUISER_SPRITE
	case types.DESTROYER:
		return types.DESTROYER_SPRITE
	case types.SUBMARINE:
		return types.SUBMARINE_SPRITE
	default:
		return []rune{}
	}
}

func CheckOverlap(grid map[types.Position]types.Cell, ship types.Ship) bool {
	for i := ship.StartPosition.X; i <= ship.EndPosition.X; i++ {
		for j := ship.StartPosition.Y; j <= ship.EndPosition.Y; j++ {
			if grid[types.Position{X: i, Y: j}].Type != types.CELL_CURSOR && grid[types.Position{X: i, Y: j}].Type != types.CELL_WATER {
				return true
			}
		}
	}
	return false
}

func GetShipType(ship types.ShipType) string {
	switch ship {
	case types.AIRCRAFT_CARRIER:
		return "Carrier"
	case types.BATTLESHIP:
		return "Battleship"
	case types.CRUISER:
		return "Cruiser"
	case types.SUBMARINE:
		return "Submarine"
	case types.DESTROYER:
		return "Destroyer"
	default:
		return "Unknown"
	}
}

func GetCellType(cell types.CellType) string {
	switch cell {
	case types.CELL_BLANK:
		return "Water"
	case types.CELL_SHIP:
		return "Ship"
	case types.CELL_MISS:
		return "Miss"
	case types.CELL_DESTROYED:
		return "Destroyed"
	case types.CELL_CURSOR:
		return "Cursor"
	case types.CELL_WATER:
		return "Water"
	default:
		return "Unknown"
	}
}
