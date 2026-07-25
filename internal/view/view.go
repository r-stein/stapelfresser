package view

// IsVertical indicated whether the screen is horizontal or vertical
var IsVertical = true

var BoardX int = 100
var BoardY int = 100

const BoardTileSize int = 100
const BoardSize int = 300

func SetIsVertical(vertical bool, xBuffer, yBuffer int) {
	IsVertical = vertical
	if vertical {
		BoardX = 0 + xBuffer
		BoardY = 200 + yBuffer
	} else {
		BoardX = 100 + xBuffer
		BoardY = 100 + yBuffer
	}
}
