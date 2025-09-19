package utils

import "github.com/vector-ops/goships/types"

func ValidateEntityPosition(e types.Ship, gridHeight, gridWidth int) bool {
	return e.StartPosition.X <= gridWidth && e.StartPosition.Y <= gridHeight && e.EndPosition.X <= gridWidth && e.EndPosition.Y <= gridHeight && e.StartPosition.X >= 0 && e.StartPosition.Y >= 0 && e.EndPosition.X >= 0 && e.EndPosition.Y >= 0
}

func GetCellType(shipType types.ShipType) types.CellType {
	switch shipType {
	case types.BATTLESHIP:
		return types.CELL_BATTLESHIP
	case types.AIRCRAFT_CARRIER:
		return types.CELL_CARRIER
	case types.CRUISER:
		return types.CELL_CRUISER
	case types.DESTROYER:
		return types.CELL_DESTROYER
	case types.SUBMARINE:
		return types.CELL_SUBMARINE
	default:
		return types.CELL_BLANK
	}
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
