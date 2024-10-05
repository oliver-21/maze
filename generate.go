package main

import (
	"bytes"
	_ "embed"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type path struct {
	crossable bool
	direction int8 // +1 for forwards -1 for backwards
}

type cell struct {
	//    | (x)>
	// ---
	// (y)
	x, y path
	internal
}

type row []cell

type coor struct {
	x, y int
}

func fillerCell(x, y bool) cell {
	return cell{path{x, 1}, path{y, 1}, internal{false, genItems()}}
}

// directions will always by > or V
func rowOnly(x int) row {
	line := make(row, x)
	for i := range line {
		line[i] = fillerCell(true, false)
	}
	return line
}

// directions will always by > or V
var downOnly = fillerCell(false, true)

func (m *Maze) UpdateSize() {
	width, height := text.Measure(m.String(), m.font, float64(m.scale))
	m.width, m.height = int(width), int(height)
}

var (
	//go:embed JetBrainsMono_Regular.ttf
	fontSorce       []byte
	cascidiaMono, _ = text.NewGoTextFaceSource(bytes.NewReader(fontSorce))
)

func basicMaze(width, height int) Maze {
	firstLine := rowOnly(width + 1)
	firstLine[0] = fillerCell(true, true)

	var res = []row{
		firstLine,
	}

	for i := height; i != 0; i-- {
		line := rowOnly(width + 1)
		line[0] = downOnly
		line[len(line)-1] = downOnly
		res = append(res, line)
	}

	res[height][width] = fillerCell(false, false)
	m := Maze{
		area: res,
		coor: coor{width, height},
		font: &text.GoTextFace{
			Source: cascidiaMono,
			Size:   20,
		},
		edge:  2,
		scale: 22,
	}
	m.UpdateSize()
	return m
}

func (m Maze) posDir() []coor {
	var (
		res []coor
		x   = m.x
		y   = m.y
	)
	if x > 1 {
		res = append(res, coor{-1, 0})
	}
	if y > 1 {
		res = append(res, coor{0, -1})
	}
	if x < len(m.area[0])-1 {
		res = append(res, coor{1, 0})
	}
	if y < len(m.area)-1 {
		res = append(res, coor{0, 1})
	}
	return res
}

func (m *Maze) add(c coor) {
	m.x += c.x
	m.y += c.y
}

func (m *Maze) modify(move coor, callback func(p *path)) {
	curr := m.coor
	switch {
	case move.x == -1:
		curr.x--
	case move.y == -1:
		curr.y--
	}
	cell := &m.area[curr.y][curr.x]
	if move.x != 0 {
		callback(&cell.x)
	} else {
		callback(&cell.y)
	}
}

func (m *Maze) deletePointing() {
	for _, pos := range m.posDir() {
		correct := pos.x + pos.y
		m.modify(pos, func(p *path) {
			if p.crossable && correct == int(p.direction) {
				p.crossable = false
			}
		})
	}
}

func (m *Maze) moveCenter() {
	possibilities := m.posDir()
	next := possibilities[rand.Intn(len(possibilities))]
	// fmt.Println(possibilities, ":", m.coor)
	correct := next.x + next.y
	m.modify(next, func(p *path) {
		p.crossable = true
		p.direction = int8(correct)
	})
	m.add(next)
	m.deletePointing()
}

func (m *Maze) exitIn(col int) {
	row := rand.Intn(len(m.area) - 1)
	m.area[row+1][col].x.crossable = true
}

func (m *Maze) addExits() {
	m.exitIn(0)
	m.exitIn(len(m.area[0]) - 1)
}

func genMaze() *Maze {
	m := basicMaze(20, 20)
	// m := basicMaze(40, 40)
	for i := 0; i < 100000; i++ {
		m.moveCenter()
		// time.Sleep(time.Second / 5)
		// fmt.Println(m)
	}
	m.addExits()
	return &m
}
