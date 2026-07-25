package game

import (
	"bytes"
	"image/color"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"

	"github.com/r-stein/stapelfresser/internal/model"
	"github.com/r-stein/stapelfresser/internal/view"
)

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

// Game is the game structure
type Game struct {
	Board        *model.Board
	Stones       []*view.Stone
	ActivePlayer model.Player
	GameOver     bool
	GameResult   model.GameResult
}

func New(board *model.Board) *Game {
	stones := make([]*view.Stone, 0, len(board.Stones))
	for _, s := range board.Stones {
		stones = append(stones, &view.Stone{
			Inner: s,
		})
	}
	return &Game{
		Board:        board,
		Stones:       stones,
		ActivePlayer: model.Player1,
		GameResult:   model.InProgress,
	}
}

// moveStoneToFront moves a stone to the front of the view
func (g *Game) moveStoneToFront(s *view.Stone) {
	index := -1
	for i, ss := range g.Stones {
		if ss == s {
			index = i
			break
		}
	}
	g.Stones = append(g.Stones[:index], g.Stones[index+1:]...)
	g.Stones = append(g.Stones, s)
}

// field calculates the field on the board based on the view coordinates,
// last argument returns whether it is a valid field
func field(x, y int) (int, int, bool) {
	// so we have an offset of 100 and each field is the size of 100
	x -= 100
	y -= 100
	if x < 0 || y < 0 || x > 300 || y > 300 {
		return -1, -1, false
	}
	return x / 100, y / 100, true
}

// stoneAt returns the stone at a view location
// return nil if no stone exists at that location
func (g *Game) stoneAt(x, y int) *view.Stone {
	for _, s := range g.Stones {
		if s.At(x, y) {
			return s
		}
	}
	return nil
}

// draggedStone returns the current dragged stone in the view
func (g *Game) draggedStone() *view.Stone {
	for _, s := range g.Stones {
		if s.Dragged {
			return s
		}
	}
	return nil
}

func (g *Game) Update() error {
	if g.GameResult != model.InProgress {
		return nil
	}
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		s := g.draggedStone()
		if s != nil {
			s.Drag(x, y)
		} else {
			s = g.stoneAt(x, y)
			if s != nil && s.Inner.Player == g.ActivePlayer {
				s.StartDrag(x, y)
				g.moveStoneToFront(s)
			}
		}
	} else if inpututil.MouseButtonPressDuration(ebiten.MouseButtonLeft) > 0 {
		s := g.draggedStone()
		if s != nil {
			s.Drag(ebiten.CursorPosition())
		}
	} else {
		s := g.draggedStone()
		if s != nil {
			// get the field location
			s.EndDrag()
			x, y, valid := field(ebiten.CursorPosition())
			if valid {
				err := s.Inner.Move(x, y)
				if err != nil {
					slog.Info("failed to move", "err", err)
				} else {
					g.ActivePlayer = !g.ActivePlayer
				}
			}
		}
	}

	if result, over := g.Board.GameOver(); over {
		_ = result
		g.GameOver = true
		g.GameResult = result
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.Gray{Y: uint8(12)})

	offsetX, offsetY := 105, 105

	// TODO we should make the sizes variables...

	// draw the info board
	opi := &ebiten.DrawImageOptions{}
	imi := ebiten.NewImage(200, 70)
	opi.GeoM.Translate(150., 15)
	outputText := "Player 1 turn"
	if g.ActivePlayer == model.Player2 {
		outputText = "Player 2 turn"
	}
	if g.GameOver {
		outputText = "Game Over: "
		switch g.GameResult {
		case model.Player1Win:
			outputText += "Player 1 wins"
		case model.Player2Win:
			outputText += "Player 2 wins"
		case model.Tie:
			outputText += "Tie"
		}
	}
	text.Draw(imi, outputText, &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   10,
	}, &text.DrawOptions{})
	screen.DrawImage(imi, opi)

	// draw the board
	opb := &ebiten.DrawImageOptions{}
	imb := ebiten.NewImage(300, 300)
	opb.GeoM.Translate(100., 100.)
	imb.Fill(colornames.Darkgray)
	screen.DrawImage(imb, opb)
	for y := range 3 {
		for x := range 3 {
			op := &ebiten.DrawImageOptions{}
			img := ebiten.NewImage(90, 90)
			op.GeoM.Translate(float64(offsetX+100*x), float64(offsetY+100*y)) // TODO improve
			img.Fill(colornames.Darkslategrey)
			screen.DrawImage(img, op)
		}
	}

	// draw each stone
	// TODO: draw the dragged last
	for _, vs := range g.Stones {
		vs.Draw(screen)
	}

	// if true {
	// 	return
	// }
	// loc := func(x, y int) [2]float64 {
	// 	return [2]float64{
	// 		float64(offsetX + 100*x),
	// 		float64(offsetY + 100*y),
	// 	}
	// }
	// for i, s := range g.Board.Stones {
	// 	if s.IsOnBoard() {
	// 		continue
	// 	}
	// 	imgSize := 45
	// 	op := &ebiten.DrawImageOptions{}

	// 	x := float64(25)
	// 	if s.Player == model.Player2 {
	// 		x = float64(425)
	// 	}
	// 	y := float64(offsetY + (i%6)*50)
	// 	op.GeoM.Translate(x, y) // TODO improve

	// 	img := ebiten.NewImage(imgSize, imgSize)
	// 	if s.Player == model.Player1 {
	// 		img.Fill(colornames.Darkmagenta)
	// 	} else {
	// 		img.Fill(colornames.Darkgreen)
	// 	}

	// 	opI := &ebiten.DrawImageOptions{}
	// 	inSize := (s.Size + 1) * 8
	// 	t := float64(imgSize/2 - inSize/2)
	// 	opI.GeoM.Translate(t, t)
	// 	imgI := ebiten.NewImage(inSize, inSize)
	// 	imgI.Fill(color.Black)
	// 	img.DrawImage(imgI, opI)

	// 	screen.DrawImage(img, op)
	// }

	// rows := g.Board.Show()
	// for y, row := range rows {
	// 	for x, s := range row {
	// 		if s != nil {
	// 			op := &ebiten.DrawImageOptions{}
	// 			img := ebiten.NewImage(90, 90)
	// 			op.GeoM.Translate(loc(x, y)[0], loc(x, y)[1]) // TODO improve
	// 			// img.Fill(color.RGBA{R: 255})
	// 			if s.Player == model.Player1 {
	// 				img.Fill(colornames.Darkmagenta)
	// 			} else {
	// 				img.Fill(colornames.Darkgreen)
	// 			}
	// 			opI := &ebiten.DrawImageOptions{}
	// 			inSize := (s.Size + 1) * 16
	// 			t := float64(90/2 - inSize/2)
	// 			opI.GeoM.Translate(t, t)
	// 			imgI := ebiten.NewImage(inSize, inSize)
	// 			imgI.Fill(color.Black)
	// 			img.DrawImage(imgI, opI)
	// 			screen.DrawImage(img, op)
	// 		}
	// 	}
	// }
	// // if true {

	// // 	op := &ebiten.DrawImageOptions{}
	// // 	img := ebiten.NewImage(100, 100)
	// // 	img.Fill(color.RGBA{R: 255})
	// // 	screen.DrawImage(img, op)
	// // }
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 500, 500
	// return outsideWidth, outsideHeight
}
