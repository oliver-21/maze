package main

import (
	"fmt"
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// type fcoor struct {
// 	x, y float64
// }

func (m *Maze) blockToImageCoords(x, y int) (float64, float64) {
	blockWidth, blockHight := text.Measure(cell{}.String(), m.font, float64(m.scale))
	fmt.Print(x, " ", y)
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

func (m *Maze) DrawItem(screen *ebiten.Image, pos coor, item item) {
	x, y := m.blockToImageCoords(pos.x, pos.y)
	fmt.Println(":", x, ",", y)
	bounds := screen.Bounds()
	min := bounds.Min
	max := bounds.Max
	area := screen.SubImage(image.Rect(min.X+int(x), min.Y+int(y), max.X, max.Y)).(*ebiten.Image)
	text.Draw(area, item.string, m.font, &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: float64(m.scale)}})
}

func (m *Maze) DrawItems(screen *ebiten.Image, pos coor, int internal) {
	for _, item := range int.background {
		m.DrawItem(screen, pos, item)
	}
	if int.isCoin {
		m.DrawItem(screen, pos, coin)
	}
}

func (m *Maze) DrawStuff(screen *ebiten.Image) {
	for y, row := range m.area {
		for x, cell := range row {
			m.DrawItems(screen, coor{x, y}, cell.internal)
		}
	}
}
