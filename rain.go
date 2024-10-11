package main

import (
	"image"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
)

type rainSorce struct {
	loc    coor[int]
	isRain [2]bool
	len    int
	offset float64
	data   string
}

var waterColor = color.RGBA{54, 54, 117, 255}

const underScoreOffset = 2

func (r *rainSorce) draw(screen *ebiten.Image, m *Maze) {
	dx, dy := m.blockToImageCoords(1, float64(r.len))
	x, y := m.blockToImageCoords(float64(r.loc.x), float64(r.loc.y))
	y += underScoreOffset
	place := image.Rect(int(x), int(y), int(x+dx), int(y+dy))
	droplets := screen.SubImage(place).(*ebiten.Image)

	m.drawText(droplets, coor[float64]{x, y - r.offset}, r.data, waterColor)

	// options := &ebiten.DrawImageOptions{}
	// options.GeoM.Translate(x, y)
	// screen.DrawImage(droplets, options)
}

func (m *Maze) drawRain(screen *ebiten.Image) {
	for _, drops := range m.rain {
		drops.draw(screen, m)
	}
}

func genDrop(r *rainSorce) (ans string) {
	const chars = ";.:'  "
	for _, e := range r.isRain {
		if e {
			ans += string(chars[rand.IntN(len(chars))])
		} else {
			ans += " "
		}
	}
	// println(ans)
	return
}

func (r *rainSorce) updateRainSorce(m *Maze) {
	if r.offset < 0 {
		_, yDist := m.blockToImageCoords(0, 1)
		r.data = genDrop(r) + "\n" + r.data
		r.offset = yDist
	}
	if len(r.data)/3 > r.len {
		r.data = r.data[:len(r.data)-1]
	}
	r.offset -= 1
}

func (m *Maze) updateRain() {
	for _, rain := range m.rain {
		rain.updateRainSorce(m)
	}
}

func (m *Maze) addRain() {
	_, maxOffset := m.blockToImageCoords(0, 1)
	for i, line := range m.area[1:] {
		for j := range line[1:int(m.max.x)] {
			coo := coor[int]{j + 1, i + 1}
			// no pipe above
			prev := m.area[coo.y-1][coo.x].rainSorce
			if !m.canMoveInDir(coo, coor[int]{0, -1}) {
				// mabey i could make this more styalised by only allowing rain in where there is a gap bellow as well
				// the rain would then only be in longer tunnels.
				// Hoever I think this would look less realistc for some reason
				if rand.IntN(15) == 1 {
					var isWater [2]bool
					isWater[rand.IntN(len(isWater))] = true
					if prev != nil && isWater == prev.isRain {
						isWater[0], isWater[1] = isWater[1], isWater[0]
					}
					stream := &rainSorce{coo, isWater, 1, float64(rand.Float64() * maxOffset), ""}
					m.area[coo.y][coo.x].rainSorce = stream
					m.rain = append(m.rain, stream)
				}
			} else if coo.y > 0 {
				sorce := m.area[coo.y-1][coo.x].rainSorce
				m.area[coo.y][coo.x].rainSorce = sorce
				if sorce != nil {
					sorce.len += 1
				}
			}
		}
	}
}
