package utils

import "strconv"

func To2DSlice[T any](slice1 []T, slice2 []T) [][]T {
	if len(slice1) == 0 && len(slice2) == 0 {
		return nil
	}
	return [][]T{slice1, slice2}
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}
