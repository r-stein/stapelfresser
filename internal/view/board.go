package view

import (
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/colornames"
)

func DrawBoard(screen *ebiten.Image) {
	// draw the board
	opb := &ebiten.DrawImageOptions{}
	imb := ebiten.NewImage(300, 300)
	opb.GeoM.Translate(float64(BoardX), float64(BoardY))
	imb.Fill(colornames.Darkgray)
	screen.DrawImage(imb, opb)
	for y := range 3 {
		for x := range 3 {
			op := &ebiten.DrawImageOptions{}
			img := ebiten.NewImage(90, 90)
			op.GeoM.Translate(float64(BoardX+5+BoardTileSize*x), float64(BoardY+5+BoardTileSize*y)) // TODO improve
			img.Fill(colornames.Darkslategrey)
			screen.DrawImage(img, op)
		}
	}
}
