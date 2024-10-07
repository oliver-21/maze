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

func (m *Maze) blockToImageCoords(x, y float64) (float64, float64) {
	blockWidth, _ := m.textSize(cell{}.String())
	// fmt.Print(x, " ", y)
	edgWidth, _ := m.textSize(strings.Repeat(" ", m.fedge))
	// return float64(x)*blockWidth - edgWidth, float64(y) * blockHight
	return x*blockWidth - edgWidth, y * float64(m.scale)
}

func plants(chars string, outOf float64) string {
	ans := ""
	for i := 0; i < cellLen; i++ {
		if rand.Float64()*outOf < 1 {
			ans += string(chars[rand.IntN(len(chars))])
		} else {
			ans += " "
		}
	}
	return ans
}

var berrieColor = color.RGBA{92, 16, 102, 255}

// var batColor = color.RGBA{48, 44, 44, 255}

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
	// add in ^
	if hasUp {
		// if rand.IntN(6) == 0 { // TODO finde some other decoration
		// 	new = append(new, item{plants("^", 8), batColor}) // yes, these are bats
		// } else {
		new = append(new, item{plants(`"''`+"`", 15), berrieColor})
		// }
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

func (m *Maze) spawnItems(str string, col colorTheme, factor float64) {
	for i, line := range m.area[1:] {
		for j := range line[1:int(m.max.x)] {
			// println(c.x, ":", c.y)
			coo := coor[int]{j + 1, i + 1}
			if !m.canMoveInDir(coo, coor[int]{0, 1}) {
				cell := m.area[coo.y][coo.x]
				cellStr := plants(str, factor)
				if !cell.x.crossable {
					cellStr = cellStr[:len(str)-1]

				}
				cell.internal.background = append(
					cell.internal.background,
					item{
						cellStr,
						randLerpColor(col.a, col.b),
					},
				)
				m.area[coo.y][coo.x] = cell
				// godump.Dump(cell)
			}
		}
	}
}

type item struct {
	string
	color color.RGBA
}

type internal struct {
	isCoin     [2]bool
	coinColor  color.RGBA
	background []item
}

func (m *Maze) drawText(screen *ebiten.Image, pos coor[float64], str string, colour color.RGBA) {
	place := &ebiten.DrawImageOptions{}
	color := colour
	place.ColorScale.Scale(float32(color.R)/255, float32(color.G)/255, float32(color.B)/255, 1)
	place.GeoM.Translate(pos.x, pos.y)
	text.Draw(
		screen,
		str,
		m.font,
		&text.DrawOptions{
			LayoutOptions: text.LayoutOptions{
				LineSpacing: float64(m.scale),
			},
			DrawImageOptions: *place,
		})

}

func (m *Maze) DrawItem(screen *ebiten.Image, pos coor[int], item item) {
	x, y := m.blockToImageCoords(float64(pos.x), float64(pos.y))
	m.drawText(screen, coor[float64]{x, y}, item.string, item.color)
}

func (m *Maze) DrawItems(screen *ebiten.Image, pos coor[int], int internal) {
	// godump.Dump(pos)
	for _, item := range int.background {
		m.DrawItem(screen, pos, item)
	}
	var coinColor = m.area[pos.y][pos.x].coinColor
	x, y := m.blockToImageCoords(float64(pos.x), float64(pos.y))
	if int.isCoin[0] {
		y := y - wave(float64(pos.x*3)+y+m.offset/15)*4
		m.drawText(screen, coor[float64]{x, y}, "o", coinColor)
	}
	if int.isCoin[1] {
		y := y - wave(float64(pos.x*3+1)+y+m.offset/15)*4
		m.drawText(screen, coor[float64]{x, y}, " o", coinColor)
	}
}

func (m *Maze) DrawStuff(screen *ebiten.Image) {
	for y, row := range m.area {
		for x, cell := range row {
			m.DrawItems(screen, coor[int]{x, y}, cell.internal)
		}
	}
}

func (m *Maze) spawnPoints() []coor[int] {
	var res []coor[int]
	for y, row := range m.area[1:m.max.y] {
		for x := range row[1:m.max.x] {
			res = append(res, coor[int]{x + 1, y + 1})
		}
	}
	return res
}

const coinFactor = 0.3

func (m *Maze) AddCoins() {
	points := m.spawnPoints()
	shuffle(points)
	numCoins := 2 // int(math.Sqrt(float64(len(points))) * coinFactor)
	for _, e := range points[:numCoins] {
		coin := &m.area[e.y][e.x]
		coin.isCoin[rand.IntN(2)] = true
		coin.coinColor = color.RGBA{231, 245, 37, 255}
		// fmt.Println(e.x, " ", e.y)
	}
	m.numCoins += numCoins
}
