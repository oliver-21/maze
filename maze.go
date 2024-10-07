// TODO:
// - add player and controler
// - black out areas that are not yet reached
// - If I have extra time: add coins mode and monster mode

// Enemy: }<
// Player: @<, \., ^'^, ^^', maby animate up and down with sine wave (delay for wing opposite direction going) move together at end of down strok as well?

package main

// $Env:GOOS = 'js'
// $Env:GOARCH = 'wasm'
// go build -o web/yourgame.wasm .
// Remove-Item Env:GOOS
// Remove-Item Env:GOARCH

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
	won
	paused
	alreadyWon
)

type message struct {
	timout int
	string
}

// - Don't include maze size option - overcomplicating navigation. And no good way to deal with screen sizeing with larger mazes withough giving up resolution or trying to resize / spawn a new window which would be confusing to the player
type Maze struct {
	theme colorTheme
	area  []row
	coor[int]
	fedge, ledge    int
	width, height   int
	scale           int
	font            *text.GoTextFace
	max             coor[int] // max position won't go past this
	numCoins, score int

	player      player
	entry, exit int
	playState
	messages []message
	message
	offset float64
}

func (m *Maze) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEnter) {
		if m.playState != won {
			m.setMessage("Are you sure you want to restart? (Y/N)", 1000000)
			m.playState = paused
		} else {
			*m = *genMaze()
			m.setMessage("New Game", 2)
		}
	}
	if m.playState == paused {
		if ebiten.IsKeyPressed(ebiten.KeyY) {
			*m = *genMaze()
			m.setMessage("New Game", 2)
			m.playState = playing
		}
		if ebiten.IsKeyPressed(ebiten.KeyN) {
			m.playState = playing
			m.setMessage("", 0)
		}
	}
	m.offset += 1
	m.player.HandleCoins(m)
	m.player.Update(m)
	if m.timout <= 0 {
		if len(m.messages) != 0 {
			m.message = m.messages[0]
			m.messages = m.messages[1:]
		} else {
			m.message.string = ""
		}
	}

	if m.playState == playing && m.hasWon() {
		m.playState = won
		m.setMessage("You Won", 5)
	}
	m.timout--
	return nil
}

func (g *Maze) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (m *Maze) Draw(screen *ebiten.Image) {
	m.DrawStuff(screen)
	text.Draw(screen, m.String(), m.font, &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: float64(m.scale)}})
	m.player.Draw(screen, m)
	// pos := m.player.coor
	// cx, cy := m.textSize("^^")
	// cx /= 2
	// cy /= 2
	// darkenOutside(screen, coor[float64]{pos.x + cx, pos.y + cy}, m.max)

	x, _ := m.textSize(cell{}.String())
	m.drawText(screen, coor[float64]{x, 0}, m.getMessage(), color.RGBA{255, 255, 255, 255})
	x, _ = m.blockToImageCoords(float64(m.max.x-1), 0)
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
