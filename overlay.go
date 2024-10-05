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

type colorTheme struct {
	a, b color.RGBA
}

// TODO 3 mazes 3 themes
func randColorRange() colorTheme {
	var colors = []colorTheme{
		{ // cave
			color.RGBA{69, 6, 66, 255},
			color.RGBA{18, 42, 66, 255},
		},
		{ // grass
			color.RGBA{12, 71, 6, 255},
			color.RGBA{44, 46, 34, 255},
		},
		{ // autumn
			color.RGBA{74, 33, 25, 255},
			color.RGBA{82, 87, 16, 255},
		},
	}
	return colors[rand.IntN(len(colors))]
}

func LinearPoint(a, b color.Color, p float64) color.RGBA {
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}

	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()

	red := uint8(int(float64(int(r2>>8)-int(r1>>8))*p) + int(r1>>8))
	green := uint8(int(float64(int(g2>>8)-int(g1>>8))*p) + int(g1>>8))
	blue := uint8(int(float64(int(b2>>8)-int(b1>>8))*p) + int(b1>>8))
	alpha := uint8(int(float64(int(a2>>8)-int(a1>>8))*p) + int(a1>>8))

	c := color.RGBA{red, green, blue, alpha}
	return c
}

func randLerpColor(a, b color.Color) color.RGBA {
	return LinearPoint(a, b, rand.Float64())
}

func (m *Maze) textSize(s string) (float64, float64) {
	return text.Measure(s, m.font, float64(m.scale))
}

func (m *Maze) blockToImageCoords(x, y int) (float64, float64) {
	blockWidth, _ := m.textSize(cell{}.String())
	// fmt.Print(x, " ", y)
	edgWidth, _ := m.textSize(strings.Repeat(" ", m.fedge))
	// return float64(x)*blockWidth - edgWidth, float64(y) * blockHight
	return float64(x)*blockWidth - edgWidth, float64(y) * float64(m.scale)
}

var (
	coin = item{
		"O", color.RGBA{231, 245, 37, 255},
	}
)

func plants(chars string, outOf float64) string {
	ans := ""
	for i := 0; i < len(cell{}.String()); i++ {
		if rand.Float64()*outOf < 1 {
			ans += string(chars[rand.IntN(len(chars))])
		} else {
			ans += " "
		}
	}
	return ans
}

var berrieColor = color.RGBA{92, 16, 102, 255}

func (m *Maze) genItems(c coor[int]) {
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
		//TODO: randomly decide new color scheme
		new = append(new, item{grass, randLerpColor(m.theme.a, m.theme.b)})
	}
	if hasUp {
		new = append(new, item{plants(`"''`+"`", 15), berrieColor})
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
			m.genItems(coor[int]{j + 1, i + 1})
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

func (m *Maze) DrawItem(screen *ebiten.Image, pos coor[int], item item) {
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

func (m *Maze) DrawItems(screen *ebiten.Image, pos coor[int], int internal) {
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
			m.DrawItems(screen, coor[int]{x, y}, cell.internal)
		}
	}
}
