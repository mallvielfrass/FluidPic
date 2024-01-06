package main

type FiguresCollection struct {
	Item []Item
}
type Item struct {
	X int
	Y int
}

type Point struct {
	X, Y int
}
type ImageCenter struct {
	X int
	Y int
}
type ImageParams struct {
	width  int
	height int
}
