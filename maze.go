// TODO:
// - add player and controler
// - black out areas that are not yet reached
// - If I have extra time: add coins mode and monster mode

package main

// $Env:GOOS = 'js'
// $Env:GOARCH = 'wasm'
// go build -o web/yourgame.wasm .
// Remove-Item Env:GOOS
// Remove-Item Env:GOARCH

import (
	"image"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

// - Don't include maze size option - overcomplicating navigation. And no good way to deal with screen sizeing with larger mazes withough giving up resolution or trying to resize / spawn a new window which would be confusing to the player
type Maze struct {
	area []row
	coor[int]
	edge          int
	width, height int
	scale         int
	font          *text.GoTextFace
	max           coor[int] // max position won't go past this

}

func (m *Maze) Update() error {
	return nil
}

func (g *Maze) Layout(outsideWidth, outsideHeight int) (int, int) {
	return g.width, g.height
}

func (m *Maze) Draw(screen *ebiten.Image) {
	m.DrawStuff(screen)
	text.Draw(screen, m.String(), m.font, &text.DrawOptions{LayoutOptions: text.LayoutOptions{LineSpacing: float64(m.scale)}})

}

// TODO moving back and forth just randomly tends to keep us in one corner making larger mazes more and more expensive and this also makes mazes slightly more predictable
func main() {
	go soundtrack()
	data := genMaze()
	ebiten.SetWindowIcon([]image.Image{get("icon.png")})

	ebiten.SetWindowSize(data.width, data.height)
	if err := ebiten.RunGame(data); err != nil {
		log.Fatal(err)
	}
}

// game idea:
// for every 2 moves you make; the enemy gets to make one move toward you.
// but you are constained by the walls of the maze while your enemy can move through walls
// the goal is to get to the treasure the enemey is protecting
// Enemy: }< (red)
// Player: .\/_ (brown)
// grass: \/ (green)
// hanging: ^" (these will be green)
// -O
