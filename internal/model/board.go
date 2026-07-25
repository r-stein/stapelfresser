package model

type Board struct {
	Width, Height int
	Stones        []*Stone
}

func NewBoard() *Board {
	st := make([]*Stone, 0, 6)
	b := &Board{
		Width:  3,
		Height: 3,
		Stones: st,
	}
	for _, p := range []Player{false, true} {
		for s := range 3 {
			b.Stones = append(b.Stones, &Stone{Board: b, Player: p, Size: s, Index: s * 2})
			b.Stones = append(b.Stones, &Stone{Board: b, Player: p, Size: s, Index: s*2 + 1})
		}
	}
	return b
}

func (b *Board) Search(x, y int) *Stone {
	for _, s := range b.Stones {
		if s.At(x, y) && !s.HasBeenEaten {
			return s
		}
	}
	return nil
}

type GameResult int

const (
	InProgress GameResult = iota
	Tie
	Player1Win
	Player2Win
	Unknown
)

func samePlayer(a, b, c *Stone) (GameResult, bool) {
	if a == nil || b == nil || c == nil {
		return Unknown, false
	}
	same := a.Player == b.Player && a.Player == c.Player
	result := Unknown
	if same {
		result = Player1Win
		if a.Player == Player2 {
			result = Player2Win
		}
	}
	return result, same
}

func (b *Board) GameOver() (GameResult, bool) {
	stones := b.stones()
	for i := range 3 {
		if winner, same := samePlayer(stones[0][i], stones[1][i], stones[2][i]); same {
			return winner, true
		}
		if winner, same := samePlayer(stones[i][0], stones[i][1], stones[i][2]); same {
			return winner, true
		}
	}
	if winner, same := samePlayer(stones[0][0], stones[1][1], stones[2][2]); same {
		return winner, true
	}
	if winner, same := samePlayer(stones[0][2], stones[1][1], stones[2][0]); same {
		return winner, true
	}
	// if no player has 3 stones the game is tied, this happens if we only have 4 stones on the field
	c := 0
	for _, s := range b.Stones {
		if !s.HasBeenEaten {
			c += 1
		}
		if c > 4 {
			break
		}
	}
	if c == 4 {
		return Tie, true
	}
	return InProgress, false
}

func (b *Board) stones() [][]*Stone {
	rows := make([][]*Stone, 0, b.Height)
	for range b.Width {
		rows = append(rows, make([]*Stone, b.Height))
	}
	for _, s := range b.Stones {
		if s.IsOnBoard() {
			rows[s.Location.X][s.Location.Y] = s
		}
	}
	return rows
}

// Player is the first (false) or second player (true)
type Player = bool

const (
	Player1 Player = false
	Player2 Player = true
)

// Location is the location on the board
type Location struct {
	X, Y int
}
