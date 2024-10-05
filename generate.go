package main

import (
	"bytes"
	_ "embed"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

func shuffle[T any](arr []T) {
	for i := range arr {
		next := rand.Intn(i + 1)
		arr[i], arr[next] = arr[next], arr[i]
	}
}

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

type coor[T any] struct {
	x, y T
}

func fillerCell(x, y bool) cell {
	return cell{path{x, 1}, path{y, 1}, internal{}}
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
	//go:embed fonts/CamingoCode_Regular.ttf
	fontSorce       []byte
	cascidiaMono, _ = text.NewGoTextFaceSource(bytes.NewReader(fontSorce))
)

func basicMaze(width, height int) Maze {
	firstLine := rowOnly(width + 2)
	firstLine[0] = fillerCell(true, true)
	firstLine[len(firstLine)-1] = fillerCell(true, true)

	m := Maze{
		area:  []row{firstLine},
		coor:  coor[int]{width, height},
		theme: randColorRange(),
		font: &text.GoTextFace{
			Source: cascidiaMono,
			Size:   20,
		},
		ledge: 1,
		scale: 22,
		max:   coor[int]{width + 1, height + 1},
	}

	for i := 1; i <= height; i++ {
		line := rowOnly(width + 2)
		m.area = append(m.area, line)
		line[0] = downOnly
		line[len(line)-2] = downOnly
		line[len(line)-1] = fillerCell(true, true)
	}

	m.area[height][width] = fillerCell(false, false)
	m.UpdateSize()
	return m
}

func (m Maze) posDir() []coor[int] {
	var (
		res []coor[int]
		x   = m.x
		y   = m.y
	)
	if x > 1 {
		res = append(res, coor[int]{-1, 0})
	}
	if y > 1 {
		res = append(res, coor[int]{0, -1})
	}
	if x < m.max.x-1 {
		res = append(res, coor[int]{1, 0})
	}
	if y < m.max.y-1 {
		res = append(res, coor[int]{0, 1})
	}
	return res
}

func (m *Maze) add(c coor[int]) {
	m.x += c.x
	m.y += c.y
}

func (m *Maze) modify(move coor[int], callback func(p *path)) {
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
	cell := &m.area[row+1][col]
	cell.x.crossable = true
	var bridgeCol int
	if col == 0 {
		bridgeCol = 0
	} else {
		bridgeCol = min(len(m.area[row])-1, col+1)
	}
	m.area[row][bridgeCol].y.crossable = false
	m.area[row+1][bridgeCol].y.crossable = false
}

func (m *Maze) addExits() {
	m.exitIn(0)
	m.exitIn(len(m.area[0]) - 2)
}

func genMaze() *Maze {
	m := basicMaze(20, 20)
	for i := 0; i < 100000; i++ {
		m.moveCenter()
		// time.Sleep(time.Second / 5)
		// fmt.Println(m)
	}
	m.addExits()
	m.fillWithGrass()
	return &m
}
