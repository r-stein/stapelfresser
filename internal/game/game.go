package game

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/r-stein/stapelfresser/internal/model"
)

// Game is the game structure
type Game struct {
	Board *model.Board
}

func New(board *model.Board) *Game {
	return &Game{
		Board: board,
	}
}
func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Black)
	// draw 3 images

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
	// return outsideWidth, outsideHeight
}
