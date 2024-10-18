package main

type enemy struct {
	points []coor[int]
	offset float64
}

// blockToImageCoords()

func (e *enemy) move(m *Maze) {
	e.offset -= 1
}
