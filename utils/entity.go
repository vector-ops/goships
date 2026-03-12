package utils

import "github.com/vector-ops/goships/types"

func ValidEntityPosition(e types.Ship, gridHeight, gridWidth int) bool {
	return e.StartPosition.X <= gridWidth && e.StartPosition.Y <= gridHeight && e.EndPosition.X <= gridWidth && e.EndPosition.Y <= gridHeight && e.StartPosition.X >= 0 && e.StartPosition.Y >= 0 && e.EndPosition.X >= 0 && e.EndPosition.Y >= 0
}

func ExpectedEndPosition(position types.Position, sprite []rune, o types.Orientation) types.Position {
	if o == types.Horizontal {
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
	case types.Battleship:
		return types.BattleshipSprite
	case types.AircraftCarrier:
		return types.CarrierSprite
	case types.Cruiser:
		return types.CruiserSprite
	case types.Destroyer:
		return types.DestroyerSprite
	case types.Submarine:
		return types.SubmarineSprite
	default:
		return []rune{}
	}
}

func CheckOverlap(grid map[types.Position]types.Cell, ship types.Ship) bool {
	for i := ship.StartPosition.X; i <= ship.EndPosition.X; i++ {
		for j := ship.StartPosition.Y; j <= ship.EndPosition.Y; j++ {
			if grid[types.Position{X: i, Y: j}].Type != types.CellCursor && grid[types.Position{X: i, Y: j}].Type != types.CellWater {
				return true
			}
		}
	}
	return false
}

// GetShipType returns the string representation of ShipType
//
// Deprecated: GetShipType is deprecated. Use [ShipType.String] instead
func GetShipType(ship types.ShipType) string {
	switch ship {
	case types.AircraftCarrier:
		return "Carrier"
	case types.Battleship:
		return "Battleship"
	case types.Cruiser:
		return "Cruiser"
	case types.Submarine:
		return "Submarine"
	case types.Destroyer:
		return "Destroyer"
	default:
		return "Unknown"
	}
}

// GetCellType returns the string representation of [CellType]
//
// Deprecated: GetCellType is deprecated. Use [CellType.String] instead
func GetCellType(cell types.CellType) string {
	switch cell {
	case types.CellBlank:
		return "Water"
	case types.CellShip:
		return "Ship"
	case types.CellMiss:
		return "Miss"
	case types.CellDestroyed:
		return "Destroyed"
	case types.CellCursor:
		return "Cursor"
	case types.CellWater:
		return "Water"
	default:
		return "Unknown"
	}
}
