package main

import (
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func (m *Maze) blockToImageCoords(x, y int) (float64, float64) {
	blockWidth, blockHight := text.Measure(cell{}.String(), m.font, float64(m.scale))
	return float64(x) * blockWidth, float64(y) * blockHight
}

var (
	coin = item{
		"O", color.RGBA{231, 245, 37, 255},
	}
	grass = item{
		`\/`, color.RGBA{19, 158, 42, 255},
	}
)

func genItems() []item {
	if rand.IntN(10) == 0 {
		return []item{grass}
	}
	return []item{}
}

type item struct {
	string
	color color.RGBA
}

type internal struct {
	isCoin     bool
	background []item
}

func DrawItem(screen *ebiten.Image, pos coor, item item)

func DrawItems(screen *ebiten.Image, pos coor, int internal) {
	for _, item := range int.background {
		DrawItem(screen, pos, item)
	}
	if int.isCoin {
		DrawItem(screen, pos, coin)
	}
}

func (m *Maze) DrawStuff(screen *ebiten.Image) {
	for y, row := range m.area {
		for x, cell := range row {
			DrawItems(screen, coor{x, y}, cell.internal)
		}
	}
}
