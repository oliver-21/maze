package main

import (
	"fmt"
	"math/rand"
	"strings"
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
}

func (s cell) String() string {
	var res []byte
	if s.y.crossable {
		res = []byte("   ")
	} else {
		res = []byte("___")
	}

	if !s.x.crossable {
		res[2] = '|'
	}

	return string(res)
}

type row []cell

func (r row) Array() []byte {
	var line []byte
	for _, c := range r {
		line = append(line, []byte(c.String())...)
	}
	return line
}

func (r row) String() string {
	return string(r.Array())
}

type coor struct {
	x, y int
}

type maze struct {
	area []row
	coor
}

func (m maze) String() string {
	var lines []string
	for _, row := range m.area {
		lines = append(lines, row.String())
	}
	// prev := lines[0]
	// var res = []string{string(prev)}
	// for _, row := range
	return strings.Join(lines, "\n")
}

func fillerCell(x, y bool) cell {
	return cell{path{x, 1}, path{y, 1}}
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

func basicMaze(width, height int) maze {
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
	return maze{res, coor{width, height}}
}

func (m maze) posDir() []coor {
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

func (m *maze) add(c coor) {
	m.x += c.x
	m.y += c.y
}

func (m *maze) modify(move coor, update func(p *path)) {
	curr := m.coor
	switch {
	case move.x == -1:
		curr.x--
	case move.y == -1:
		curr.y--
	}
	cell := &m.area[curr.y][curr.x]
	if move.x != 0 {
		update(&cell.x)
	} else {
		update(&cell.y)
	}
}

func (m *maze) deletePointing() {
	for _, pos := range m.posDir() {
		correct := pos.x + pos.y
		m.modify(pos, func(p *path) {
			if p.crossable && correct == int(p.direction) {
				p.crossable = false
			}
		})
	}
}

func (m *maze) update() {
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

func (m *maze) exitIn(col int) {
	row := rand.Intn(len(m.area))
	m.area[row][col].x.crossable = true
}

func (m *maze) addExits() {
	m.exitIn(0)
	m.exitIn(len(m.area[0]) - 1)
}

// TODO moving back and forth just randomly tends to keep us in one corner making larger mazes more and more expensive and this also makes mazes slightly more predictable
func main() {
	m := basicMaze(30, 30)
	for i := 0; i < 100000; i++ {
		m.update()
		// time.Sleep(time.Second)
	}
	m.addExits()
	fmt.Println(m)
}
