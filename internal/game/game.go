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
	Board          *model.Board
	Stones         []*view.Stone
	ActivePlayer   model.Player
	GameOver       bool
	GameResult     model.GameResult
	touchIDs       []ebiten.TouchID
	mouseX, mouseY int
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
		touchIDs:     make([]ebiten.TouchID, 0),
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
	x -= view.BoardX
	y -= view.BoardY
	if x < 0 || y < 0 || x > view.BoardSize || y > view.BoardSize {
		return -1, -1, false
	}
	return x / view.BoardTileSize, y / view.BoardTileSize, true
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

func (g *Game) PressAt(x, y int) {
	g.mouseX, g.mouseY = x, y
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
}

func (g *Game) Release() {
	s := g.draggedStone()
	if s != nil {
		// get the field location
		s.EndDrag()
		x, y, valid := field(g.mouseX, g.mouseY)
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

func (g *Game) Update() error {
	if g.GameResult != model.InProgress {
		return nil
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		g.PressAt(ebiten.CursorPosition())
	} else if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		g.Release()
	}

	wasTouched := len(g.touchIDs) > 0
	g.touchIDs = ebiten.AppendTouchIDs(g.touchIDs[:0])
	if len(g.touchIDs) > 0 {
		touch := g.touchIDs[0]
		g.PressAt(ebiten.TouchPosition(touch))
	} else if wasTouched {
		g.Release()
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

	// offsetX, offsetY := 100, 100
	// if view.IsVertical {
	// 	offsetX = 0
	// 	offsetY = 200
	// }

	// TODO we should make the sizes variables...

	// draw the info board
	opi := &ebiten.DrawImageOptions{}
	imi := ebiten.NewImage(200, 70)
	opi.GeoM.Translate(float64(view.BoardX+50), 15)
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
	view.DrawBoard(screen)

	// draw each stone
	for _, vs := range g.Stones {
		vs.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	view.SetIsVertical(outsideWidth < outsideHeight)
	if outsideWidth < outsideHeight {
		// view.SetIsVertical(true) // TODO check whether we should improve this
		// 200 at top:
		// 		100 at top and bottom for stuff
		// 		100 at top for stones
		// 100 at bottom for stones
		// 50 at bottom for symmetry
		return 330, 600
	}
	// 100 on top for stuff and 100 bottom for symetry
	// 100 left and right for the stones
	return 500, 500
	// return outsideWidth, outsideHeight
}
