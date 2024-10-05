package main

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type key struct {
	ebiten.Key
	image.Point
}

var keyCoords = []key{
	{ebiten.KeyArrowUp, image.Point{0, -1}},
	{ebiten.KeyW, image.Point{0, -1}},

	{ebiten.KeyArrowDown, image.Point{0, 1}},
	{ebiten.KeyS, image.Point{0, 1}},

	{ebiten.KeyArrowLeft, image.Point{-1, 0}},
	{ebiten.KeyA, image.Point{-1, 0}},

	{ebiten.KeyArrowRight, image.Point{1, 0}},
	{ebiten.KeyD, image.Point{1, 0}},
}

const screenHeight, screenWidth = 16, 16

func mod(a, b int) int {
	if a < 0 {
		a += (-a/b + 1) * b
	}
	return a % b
}

// clamps the given int between -1 and 1
func forceToSurrounding(x int) int {
	return min(max(-1, x), 1)
}

func keyMovement() (dx, dy int) {
	for _, e := range keyCoords {
		if ebiten.IsKeyPressed(e.Key) {
			dx += e.X
			dx = forceToSurrounding(dx)
			dy += e.Y
			dy = forceToSurrounding(dy)
		}
	}
	if dy != 0 {
		dx = 0
	}
	return
}
