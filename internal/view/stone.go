package view

import (
	"bytes"
	"fmt"
	"log"
	"log/slog"

	"github.com/hajimehoshi/ebiten/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"golang.org/x/image/colornames"

	"github.com/r-stein/stapelfresser/internal/model"
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

type Stone struct {
	Inner      *model.Stone
	Dragged    bool
	dragCenter [2]int
}

type location struct {
	/*
		  ...
		-- x1 --
		|      |
		y1     y2
		|      |
		-- x2 --
	*/
	// ...
	// -- x1 --
	// |
	// y1
	x1, x2, y1, y2 float64
}

func (l *location) unpack() (float64, float64, float64, float64) {
	return l.x1, l.x2, l.y1, l.y2
}

func (vs *Stone) location() *location {
	s := vs.Inner
	offsetX, offsetY := 105., 105.

	imgSize := 90.
	if vs.Dragged {
		imgSize = 60.
	} else if !s.IsOnBoard() {
		imgSize /= 2
	}
	var x1, x2, y1, y2 float64
	switch {
	case vs.Dragged:
		cx, cy := float64(vs.dragCenter[0]), float64(vs.dragCenter[1]) // TODO can we use unpack?
		x1 = cx - imgSize/2
		y1 = cy - imgSize/2
	case vs.Inner.HasBeenEaten:
		return nil
	case vs.Inner.IsOnBoard():
		x1, y1 = offsetX+float64(100*s.Location.X), offsetY+float64(100*s.Location.Y)
	default:
		// by default it is at it starting position
		// XXX(r-stein): We may use this to switch the location for mobile modus
		x1 = 5.
		if s.Player == model.Player2 {
			x1 = offsetX + 300.
		}
		y1 = offsetY + float64(50*s.Index)
	}

	x2, y2 = x1+imgSize, y1+imgSize
	return &location{x1, x2, y1, y2}
}

func (vs *Stone) At(x, y int) bool {
	loc := vs.location()
	if loc == nil {
		return false
	}
	xf, yf := float64(x), float64(y)
	x1, x2, y1, y2 := loc.unpack()
	return x1 <= xf && xf <= x2 && y1 <= yf && yf <= y2
	// return false // TODO
}

func (vs *Stone) StartDrag(x, y int) {
	slog.Debug("start srag", "x", x, "y", y)
	vs.Dragged = true
	vs.dragCenter = [2]int{x, y}
}

func (vs *Stone) Drag(x, y int) {
	slog.Debug("drag", "x", x, "y", y)
	vs.Dragged = true
	vs.dragCenter = [2]int{x, y}
}

func (vs *Stone) EndDrag() {
	if vs == nil {
		return
	}
	vs.Dragged = false
	// vs.dragCenter = nil
}

func (vs *Stone) Draw(screen *ebiten.Image) {
	// TODO more static?
	// offsetX, offsetY := 105, 105

	s := vs.Inner
	col := colornames.Darkmagenta
	if s.Player == model.Player2 {
		col = colornames.Darkgreen
	}
	loc := vs.location()
	if loc == nil {
		return
	}
	x1, x2, y1, y2 := loc.unpack()
	// slog.Debug("draw stone", "s", s, "loc", loc)
	op := &ebiten.DrawImageOptions{}
	img := ebiten.NewImage(int(x2-x1), int(y2-y1))
	img.Fill(col)
	text.Draw(img, fmt.Sprintf("s: %d", s.Size), &text.GoTextFace{
		Source: mplusFaceSource,
		Size:   (x2 - x1) / 6,
	}, &text.DrawOptions{})
	op.GeoM.Translate(x1, y1)
	screen.DrawImage(img, op)

	// imgSize := 90
	// if !s.IsOnBoard() {
	// 	imgSize /= 2
	// 	return // TODO
	// }

	// op := &ebiten.DrawImageOptions{}
	// img := ebiten.NewImage(imgSize, imgSize)
	// img.Fill(col)
	// text.Draw(img, fmt.Sprintf("s: %d", s.Size), &text.GoTextFace{
	// 	Source: mplusFaceSource,
	// 	Size:   float64(imgSize) / 6,
	// }, &text.DrawOptions{})
	// // TODO check on board
	// op.GeoM.Translate(float64(offsetX+100*s.Location.X), float64(offsetY+100*s.Location.Y))

	// screen.DrawImage(img, op)

	_ = col
}
