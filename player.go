package main

import (
	"image"
	"image/color"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type player struct {
	coor[float64]
	dir  coor[int]
	next coor[int]

	speed float64 // max 1.0
	time  float64
	tilt  float64
	goal  coor[int]
}

type key struct {
	ebiten.Key
	image.Point
}

var keyCoords = []key{
	//  W            ^
	// A S D       < v >
	{ebiten.KeyArrowUp, image.Point{0, -1}},
	{ebiten.KeyW, image.Point{0, -1}},

	{ebiten.KeyArrowDown, image.Point{0, 1}},
	{ebiten.KeyS, image.Point{0, 1}},

	{ebiten.KeyArrowLeft, image.Point{-1, 0}},
	{ebiten.KeyA, image.Point{-1, 0}},

	{ebiten.KeyArrowRight, image.Point{1, 0}},
	{ebiten.KeyD, image.Point{1, 0}},
}

func mod(a, b int) int {
	if a < 0 {
		a += (-a/b + 1) * b
	}
	return a % b
}

// clamps the given int between -1 and 1
func forceTo1[T int | float64](x T) T {
	return min(max(-1, x), 1)
}

func keyMovement() (dx, dy int) {
	for _, e := range keyCoords {
		if ebiten.IsKeyPressed(e.Key) {
			dx += e.X
			// dx = forceTo1(dx)
			dy += e.Y
			// dy = forceTo1(dy)
		}
	}
	if dy != 0 {
		dx = 0
	}
	return //
}

var batColor = color.RGBA{110, 94, 86, 255}

const batMovmentRatio = 1.2

func (p *player) Draw(screen *ebiten.Image, m *Maze) {
	const tiltEffect = 1.7
	const nearnessEffect = 0.7
	x, y := m.blockToImageCoords(p.x, p.y)
	y += tiltEffect
	movement := (math.Sin(p.time) + 1) / 2
	y += movement * 4
	sep := movement * movement * batMovmentRatio
	tilt := p.tilt * tiltEffect
	tilt *= (batMovmentRatio - sep) * nearnessEffect
	m.drawText(screen, coor[float64]{x + sep, y - tilt}, "^", batColor)
	m.drawText(screen, coor[float64]{x - sep, y + tilt}, " ^", batColor)
}

func isWithin(a, b, dx float64) bool {
	return math.Abs(a-b) <= dx
}

func move(pos *float64, goal int, movement float64) {
	if isWithin(*pos, float64(goal), movement) {
		*pos = float64(goal)
	} else {
		if *pos < float64(goal) {
			*pos += movement
		} else {
			*pos -= movement
		}
	}
	// println(*pos)
}

func (p *player) Update(m *Maze) {
	dx, dy := keyMovement()
	// fmt.Println(dx, dy)
	// update posible direction up to decision boundery
	if dx != 0 || dy != 0 {
		p.next = coor[int]{dx, dy}
	}
	if p.dir.x != 0 {
		p.tilt += float64(p.dir.x) * 0.7
		p.tilt = forceTo1(p.tilt)
	} else {
		p.tilt *= 0.9
		// set to 1 if close enough
		if p.tilt < 0.15 && p.tilt > -0.15 {
			p.tilt = 1
		}
	}
	p.time += 0.3
	p.speed *= 1.1
	p.speed = min(p.speed, 0.1)

	move(&p.coor.x, p.goal.x, p.speed)
	move(&p.coor.y, p.goal.y, p.speed)

	if p.coor.x == math.Round(p.coor.x) || p.coor.y == math.Round(p.coor.y) {
		// m.area[(p.coor.y)][(p.coor.x)]
		prev := p.dir
		p.dir = coor[int]{}
		if m.canMoveInDir(coor[int]{int(p.coor.x), int(p.coor.y)}, p.next) {
			p.dir = p.next
			p.goal.x += p.dir.x
			p.goal.y += p.dir.y
			if p.dir != prev {
				p.speed = 0.02
			}
		}
		p.next = coor[int]{}
	}
}

func (m *Maze) canMoveInDir(coor coor[int], dir coor[int]) bool {
	cpy := *m
	cpy.coor = coor
	// println(c.x, ":", c.y)
	canMove := false
	for _, pos := range m.posDir() {
		if pos == dir {
			cpy.modify(pos, func(pa *path) {
				if pa.crossable {
					canMove = true
				}
			})
		}
	}
	return canMove
	// godump.Dump(p)

}

func (m *Maze) hasWon() bool {
	return int(m.player.coor.x) == m.max.x
}
