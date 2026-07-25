package main

import (
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/r-stein/stapelfresser/internal/game"
	"github.com/r-stein/stapelfresser/internal/model"
)

func main() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
	slog.Info("Starting game...")

	// g := &game.Game{
	// 	Board: model.NewBoard(),
	// }

	g := game.New(model.NewBoard())
	// g.Board.Stones = append(g.Board.Stones, &model.Stone{
	// 	Board:        g.Board,
	// 	Size:         1,
	// 	Player:       false,
	// 	HasBeenEaten: false,
	// }, &model.Stone{
	// 	Board:  g.Board,
	// 	Player: true,
	// 	Size:   2,
	// })
	g.Board.Stones[3].Move(1, 2)
	g.Board.Stones[8].Move(0, 0)

	ebiten.SetWindowSize(1200, 1080)

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("", "error", err)
		os.Exit(1)
	}
}
