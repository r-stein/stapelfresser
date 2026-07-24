package model

import (
	"errors"
	"fmt"
)

type Board struct {
	Width, Height int
	stones        []Stone
}

func NewBoard() *Board {
	return &Board{
		Width:  3,
		Height: 3,
		stones: []Stone{},
	}
}

func (b *Board) Search(x, y int) *Stone {
	for _, s := range b.stones {
		if s.At(x, y) && !s.HasBeenEaten {
			return &s
		}
	}
	return nil
}

// Stone is stone on the board
type Stone struct {
	Board  *Board
	Player Player
	Size   int
	// HasBeenEaten signals whether the stone has been eatn
	HasBeenEaten bool
	Location     *Location
}

// IsOnBoard return whether the stone is on the board
func (s *Stone) IsOnBoard() bool {
	return s.Location != nil
}

func (s *Stone) At(x, y int) bool {
	if !s.IsOnBoard() {
		return false
	}
	return s.Location.x == x && s.Location.y == y
}

func (s *Stone) GetEaten() error {
	if !s.IsOnBoard() {
		return errors.New("you can not eat a stone that is not on the board")
	}
	s.Location = nil
	s.HasBeenEaten = true
	return nil
}

func (s *Stone) Move(x, y int) error {
	if x > s.Board.Width || y > s.Board.Width {
		return errors.New("out of board size")
	}
	other := s.Board.Search(x, y)
	if other != nil {
		if other.Size >= s.Size {
			return fmt.Errorf("can not eat other stone, size: %d, other size: %d", s.Size, other.Size)
		}
		err := other.GetEaten()
		if err != nil {
			return fmt.Errorf("eating other stone: %w", err)
		}
	}
	s.Location = &Location{x, y}
	return nil
}

// Player is the first (false) or second player (true)
type Player = bool

const (
	Player1 Player = false
	Player2 Player = true
)

// Location is the location on the board
type Location struct {
	x, y int
}
