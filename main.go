package main

import (
	"log/slog"
	"os"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/r-stein/stapelfresser/internal/game"
	"github.com/r-stein/stapelfresser/internal/model"
)

func main() {
	slog.Info("Starting game...")

	g := &game.Game{
		Board: model.NewBoard(),
	}

	if err := ebiten.RunGame(g); err != nil {
		slog.Error("", "error", err)
		os.Exit(1)
	}
}
