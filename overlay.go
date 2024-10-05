package main

import (
	"image/color"
	"math/rand/v2"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// type fcoor struct {
// 	x, y float64
// }

func (m *Maze) textSize(s string) (float64, float64) {
	return text.Measure(s, m.font, float64(m.scale))
}

func (m *Maze) blockToImageCoords(x, y int) (float64, float64) {
	blockWidth, _ := m.textSize(cell{}.String())
	// fmt.Print(x, " ", y)
	edgWidth, _ := m.textSize(strings.Repeat(" ", m.edge))
	// return float64(x)*blockWidth - edgWidth, float64(y) * blockHight
	return float64(x)*blockWidth - edgWidth, float64(y) * float64(m.scale)
}

var (
	coin = item{
		"O", color.RGBA{231, 245, 37, 255},
	}
)

var grassColor = color.RGBA{12, 71, 6, 255}
var berrieColor = color.RGBA{92, 16, 102, 255}

func plants(chars string, outOf float64) string {
	ans := ""
	for i := 0; i < 3; i++ {
		if rand.Float64()*outOf < 1 {
			ans += string(chars[rand.IntN(len(chars))])
		} else {
			ans += " "
		}
	}
	return ans
}

func (m *Maze) genItems(c coor) {
	cpy := *m
	cpy.coor = c
	var hasDown, hasUp bool
	// println(c.x, ":", c.y)
	for _, pos := range m.posDir() {
		// fmt.Println(pos)
		cpy.modify(pos, func(p *path) {
			if !p.crossable && pos.y == 1 {
				hasDown = true
			}
			if !p.crossable && pos.y == -1 {
				hasUp = true
			}
		})
	}
	cell := m.area[c.y][c.x]
	var new []item
	if hasDown {
		grass := plants(`\|/`, 5)
		new = append(new, item{grass, grassColor})
	}
	if hasUp {
		new = append(new, item{plants(`"''`, 15), berrieColor})
	}
	for _, e := range new {
		if !cell.x.crossable {
			e.string = e.string[:len(e.string)-1]
		}
		cell.internal.background = append(cell.internal.background, e)
	}
	m.area[c.y][c.x] = cell
}

func (m *Maze) fillWithGrass() {
	for i, line := range m.area[1:] {
		for j := range line[1:int(m.max.x)] {
			m.genItems(coor{j + 1, i + 1})
		}
	}
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
	// fmt.Println(":", x, ",", y)
	// bounds := screen.Bounds()
	// min := bounds.Min
	// max := bounds.Max
	// area := screen.SubImage(image.Rect(min.X+int(x), min.Y+int(y), max.X, max.Y)).(*ebiten.Image)
	place := &ebiten.DrawImageOptions{}
	color := item.color
	place.ColorScale.Scale(float32(color.R)/255, float32(color.G)/255, float32(color.B)/255, 1)
	place.GeoM.Translate(x, y)
	text.Draw(
		screen,
		item.string,
		m.font,
		&text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: float64(m.scale),
			},
			DrawImageOptions: *place,
		})
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
