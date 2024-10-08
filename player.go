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

	speed     float64 // max 1.0
	time      float64
	tilt      float64
	goal      coor[int]
	wingtime  float64
	lastMoved int
}

func (p *player) set(x, y int) {
	p.x = float64(x)
	p.y = float64(y)
	p.goal = coor[int]{x, y}
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
			dx = forceTo1(dx)
			dy += e.Y
			dy = forceTo1(dy)
		}
	}
	if dy != 0 {
		dx = 0
	}
	return //
}

var batColor = color.RGBA{145, 138, 138, 255}

const batMovmentRatio = 1.2

func (p *player) Draw(screen *ebiten.Image, m *Maze) {
	const tiltEffect = 1.7
	const nearnessEffect = 0.7
	x, y := m.blockToImageCoords(p.x, p.y)
	y += tiltEffect
	movement := (math.Sin(p.wingtime) + 1) / 2
	y += movement * 4
	sep := movement * movement * batMovmentRatio
	tilt := p.tilt * tiltEffect
	tilt *= (batMovmentRatio - sep) * nearnessEffect
	m.drawText(screen, coor[float64]{x + sep, y - tilt}, "^", batColor)
	m.drawText(screen, coor[float64]{x - sep, y + tilt}, " ^", batColor)
}

func isWithin(a, b, dx float64) bool {
	return math.Abs(a-b) <= math.Abs(dx)
}

const maxSpeed = 0.1

func (m *Maze) coinFromLoc(place coor[int]) *cell {
	// fmt.Println(m.area[place.y][place.x/3])
	return &m.area[place.y][place.x/3]
}

func (m *Maze) hasCoin(place coor[int]) bool {
	cell := m.coinFromLoc(place)
	// fmt.Println(place.x, place.y, cell.isCoin)
	loc := place.x % 3 // cellLen
	return loc < len(cell.isCoin) && cell.isCoin[loc]
}

func (m *Maze) deleteCoin(place coor[int]) {
	cell := m.coinFromLoc(place)
	loc := place.x % cellLen
	cell.isCoin[loc] = false
}

// DO NOT MODIFY
var cellLen = len(cell{}.String())

func (p *player) HandleCoins(m *Maze) {
	// print(cellLen)
	curr := p.coor
	curr.x *= float64(cellLen)
	end := math.Ceil(curr.x + 2)
	for i := int(math.Floor(curr.x)); i < int(end); i++ {
		spot := coor[int]{i, int(curr.y)}
		if m.hasCoin(spot) {
			m.deleteCoin(spot)
			m.score++
		}
	}
}

func move(pos *float64, goal int, movement float64) {
	if isWithin(*pos, float64(goal), movement+0.1) {
		*pos = float64(goal)
	} else {
		if *pos < float64(goal) {
			*pos += movement
		} else if *pos > float64(goal) {
			*pos -= movement
		}
	}
	// println(*pos)
}

func (p *player) Update(m *Maze) {
	dx, dy := keyMovement()
	if dx != 0 || dy != 0 {
		m.playState = playing
	}
	// fmt.Println(dx, dy)
	// update posible direction up to decision boundery
	wingDx := p.speed * 5
	const maxWingSpeed = maxSpeed * 5
	// if dx != 0 || dy != 0 { //TODO: add this back in but stop sliding
	p.next = coor[int]{dx, dy}
	// }

	if p.dir.x != 0 {
		p.tilt += float64(p.dir.x) * 0.3
		p.tilt = forceTo1(p.tilt)
		wingDx += maxWingSpeed / 3
	} else {
		p.tilt *= 0.9
		// set to 1 if close enough
		if p.tilt < 0.15 && p.tilt > -0.15 {
			p.tilt = 0
		}
		if p.dir.y < 0 {
			wingDx = maxWingSpeed + wingDx
		} else if p.dir.y > 0 {
			wingDx = maxWingSpeed - wingDx
		}
	}
	p.wingtime += wingDx
	p.speed *= 1.025
	p.speed = min(p.speed, maxSpeed)

	move(&p.coor.x, p.goal.x, p.speed)
	move(&p.coor.y, p.goal.y, p.speed)

	if p.coor.x == float64(p.goal.x) && p.coor.y == float64(p.goal.y) {
		// m.area[(p.coor.y)][(p.coor.x)]
		// prev := p.dir
		p.dir = coor[int]{}
		// if p.goal.x+p.dir.x <= 0 {
		// 	return
		// }
		if p.coor.x >= float64(len(m.area[0]))-1.5 && p.next.x > 0 {
			m.playState = end
			p.next.x = 0
		}
		if m.canMoveInDir(coor[int]{int(p.coor.x), int(p.coor.y)}, p.next) {
			// if p.next != p.dir {
			// 	godump.Dump(p.goal)
			// }
			p.dir = p.next
			p.goal.x += p.dir.x
			p.goal.y += p.dir.y
			// if p.dir != prev {
			// 	p.speed = 0.05
			// }
			p.lastMoved = 3
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
	return m.score == m.numCoins
}
