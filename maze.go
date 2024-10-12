// TODO:
// - black out areas that are not yet reached
package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type playState int

const (
	playing playState = iota
	end
	options
)

type message struct {
	timout int
	string
}

type gameState struct {
	prevRestart bool
	level       int
}

// - Don't include maze size option - overcomplicating navigation. And no good way to deal with screen sizeing with larger mazes withough giving up resolution or trying to resize / spawn a new window which would be confusing to the player
type Maze struct {
	// mu    sync.Mutex
	theme colorTheme
	area  []row
	coor[int]
	fedge, ledge    int
	width, height   int
	scale           int
	font            *text.GoTextFace
	max             coor[int] // max position won't go past this
	min             coor[int]
	numCoins, score int

	player      player
	entry, exit int
	playState
	messages []message
	message
	offset float64
	gameState
	givenMesage bool
	rain        []*rainSorce
}

func (m *Maze) allowEscape() {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) || ebiten.IsKeyPressed(ebiten.KeyR) {
		if !m.prevRestart {
			m.playState = options
		}
		m.prevRestart = true
	} else {
		m.prevRestart = false
	}
}

func (m *Maze) messageUpdate() {
	if m.timout <= 0 {
		if len(m.messages) != 0 {
			m.message = m.messages[0]
			m.messages = m.messages[1:]
		} else {
			m.message.string = ""
		}
	}
	m.timout--
}

func (m *Maze) Regenerate() {
	message := "New Game"
	if m.hasWon() {
		m.level++
		message = fmt.Sprintf("Reached level %v", m.level)
	}
	prev := m.gameState
	*m = *genMaze()
	m.gameState = prev
	m.setMessage(message, 3)
}

func (m *Maze) Update() error {
	// m.mu.Lock()
	// defer m.mu.Unlock()
	m.updateRain()
	m.player.HandleCoins(m)
	m.player.Update(m)
	m.allowEscape()

	if m.playState == options {
		if m.hasWon() {
			m.Regenerate()
		} else {
			m.setMessage("Are you sure you want to start a new game? (Y/N)", 0.1)
			if ebiten.IsKeyPressed(ebiten.KeyY) {
				m.Regenerate()
				m.playState = playing
			}
			if ebiten.IsKeyPressed(ebiten.KeyN) {
				m.setMessage("", 0)
				m.playState = playing
			}
		}
	}
	if m.hasWon() && !m.givenMesage {
		m.setMessage("You Won", 3)
		m.addMessage("Exit Maze or Hit Enter Key for next level", 4)
		m.givenMesage = true
	}
	if m.playState == end && m.hasWon() {
		m.Regenerate()
	}
	m.messageUpdate()
	m.offset += 1
	return nil
}

func (g *Maze) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (m *Maze) Draw(screen *ebiten.Image) {
	m.drawRain(screen)
	m.DrawStuff(screen) // We do this before otherthings so grass doesn't overlay on top of pipes
	text.Draw(screen, m.String(), m.font, &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: float64(m.scale)}})
	m.player.Draw(screen, m)
	// pos := m.player.coor
	// cx, cy := m.textSize("^^")
	// cx /= 2
	// cy /= 2
	// darkenOutside(screen, coor[float64]{pos.x + cx, pos.y + cy}, m.max)
	x, _ := m.textSize(cell{}.String())
	m.drawText(screen, coor[float64]{x, 0}, m.getMessage(), color.RGBA{255, 255, 255, 255})
	x, _ = m.blockToImageCoords(float64(len(m.area[0])-2), 0)
	m.drawText(screen, coor[float64]{x, 0}, fmt.Sprintf("%3v", m.score), color.RGBA{255, 255, 255, 255})
}

// TODO moving back and forth just randomly tends to keep us in one corner making larger mazes more and more expensive and this also makes mazes slightly more predictable
func main() {
	go soundtrack()
	data := genMaze()
	ebiten.SetWindowTitle("Maze Bat")
	// 16x16, 32x32 and 48x48 , get("icon48.png")
	ebiten.SetWindowIcon([]image.Image{get("icon.png")})
	// fmt.Println(data.width, data.height) // TODO: remove

	ebiten.SetWindowSize(data.width, data.height)
	if err := ebiten.RunGame(data); err != nil {
		log.Fatal(err)
	}
}
