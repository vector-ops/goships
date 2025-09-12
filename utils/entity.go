package utils

import "github.com/vector-ops/goships/types"

func ValidateEntityPosition(e types.Entity) bool {
	return e.StartPosition.X <= 11 && e.StartPosition.Y <= 11 && e.EndPosition.X <= 11 && e.EndPosition.Y <= 11
}

func ExpectedEndPosition(start int, sprite []rune) int {
	return start + len(sprite) - 1
}
