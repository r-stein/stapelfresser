package model

import (
	"errors"
	"fmt"
)

// Stone is stone on the board
type Stone struct {
	Board  *Board
	Player Player
	Size   int
	// Index is used to differentiate the stones of each player
	Index int
	// HasBeenEaten signals whether the stone has been eatn
	HasBeenEaten bool
	Location     *Location
	Stomach      []*Stone
}

// IsOnBoard return whether the stone is on the board
func (s *Stone) IsOnBoard() bool {
	return s.Location != nil
}

func (s *Stone) At(x, y int) bool {
	if !s.IsOnBoard() {
		return false
	}
	return s.Location.X == x && s.Location.Y == y
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
		if s.Stomach == nil {
			s.Stomach = make([]*Stone, 0, 2)
		}
		s.Stomach = append(s.Stomach, other)
	}
	s.Location = &Location{x, y}
	return nil
}
