package main

import "github.com/samber/lo"

// delete.go
func DeletePoint(x, y int, box FiguresCollection) FiguresCollection {
	nBox := lo.Filter(box.Points, func(item Point, index int) bool {
		if item.X == x && item.Y == y {
			return false
		}
		return true
	})

	return FiguresCollection{
		Points: nBox,
	}
}
