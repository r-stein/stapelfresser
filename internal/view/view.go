package view

// IsVertical indicated whether the screen is horizontal or vertical
var IsVertical = true

var BoardX int = 100
var BoardY int = 100

const BoardTileSize int = 100
const BoardSize int = 300

func SetIsVertical(vertical bool) {
	IsVertical = vertical
	if vertical {
		BoardX = 0
		BoardY = 200
	} else {
		BoardX = 100
		BoardY = 100
	}
}
