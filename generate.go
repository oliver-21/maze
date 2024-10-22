package main

import (
	"bytes"
	_ "embed"
	"image/color"
	"math"
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
	width, height := m.textSize(m.String())
	m.width, m.height = int(width), int(height)
}

var (
	//go:embed fonts/CamingoCode_Regular.ttf
	fontSorce       []byte
	cascidiaMono, _ = text.NewGoTextFaceSource(bytes.NewReader(fontSorce))
)

func (m *Maze) addMessage(mes string, timout int) {
	m.messages = append(m.messages, message{timout * 60, mes})
}
func (m *Maze) setMessage(mes string, timout float32) {
	m.message = message{int(timout * 60), mes}
	m.messages = []message{}
}
func (m *Maze) getMessage() string {
	return m.message.string
}

func basicMaze(width, height int) Maze {
	firstLine := rowOnly(width + 2)
	firstLine[0] = fillerCell(true, true)
	firstLine[len(firstLine)-1] = fillerCell(true, true)

	m := Maze{
		area:  []row{firstLine},
		coor:  coor[int]{width, height},
		theme: randColorRange(),
		font:  cascidiaMono,
		ledge: 1,
		scale: 40,
		gameState: gameState{
			level: 1, // default
		},
		entry: rand.Intn(height),
		exit:  rand.Intn(height),
		max:   coor[int]{width + 1, height + 1},
		min:   coor[int]{1, 1},
		player: player{
			speed: 1,
		},
	}
	m.addMessage("Maze Bat - By Oliver Day for Ludum Dare 56", 2)
	m.addMessage("Hit Enter/R for new level; WASD/Arrow keys to move", 4)
	m.player.set(0, m.entry+1)
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
	if x > m.min.x {
		res = append(res, coor[int]{-1, 0})
	}
	if y > m.min.y {
		res = append(res, coor[int]{0, -1})
	}
	if x < m.max.x-1 {
		res = append(res, coor[int]{1, 0})
	}
	if y < m.max.y-1 {
		res = append(res, coor[int]{0, 1})
	}
	// if len(res) != 4 {
	// 	fmt.Println(res)
	// }
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
	if curr.y < 0 || curr.y >= len(m.area) {
		return
	}
	row := m.area[curr.y]
	if curr.x < 0 || curr.x >= len(row) {
		return
	}
	cell := &row[curr.x]
	if move.x != 0 {
		callback(&cell.x)
	} else if move.y != 0 {
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

func (m *Maze) exitIn(col int, row int) {
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
	m.exitIn(0, m.entry)
	m.exitIn(len(m.area[0])-2, m.exit)
	m.area[m.exit+1][len(m.area[0])-1].isCoin[0] = true
	m.area[m.exit+1][len(m.area[0])-1].coinColor = color.RGBA{19, 4, 217, 255}
	m.area[m.exit][len(m.area[0])-1].x.crossable = false
	m.numCoins++
}

func (m *Maze) extraItems() {
	if rand.Intn(3) == 0 {
		m.spawnItems("vix", // cursed grass
			colorTheme{
				color.RGBA{100, 7, 102, 255},
				color.RGBA{29, 7, 102, 255},
			},
			40,
		)
	} else {
		m.spawnItems("nm", // bolders
			colorTheme{
				color.RGBA{68, 69, 68, 255},
				color.RGBA{61, 47, 44, 255},
			},
			40,
		)
	}
}

func (m *Maze) widenRange() {
	m.min = coor[int]{}
	m.max = coor[int]{len(m.area[0]), len(m.area)}
}

func genMaze() *Maze {
	m := basicMaze(20, 20)
	for i := 0; i < 50000; i++ {
		m.moveCenter()
		// time.Sleep(time.Second / 5)
		// fmt.Println(m)
	}
	m.addExits()
	m.addRain()
	m.extraItems()
	m.fillWithGrass()
	m.AddCoins()
	m.widenRange()
	return &m
}

func wave(x float64) float64 {
	return (math.Sin(x) + 1) / 2
}
