package game

import (
	"math/rand"
)

type RandomBot struct {
	marker string
}

func NewRandomBot() *RandomBot {
	return &RandomBot{"?"}
}
func (b *RandomBot) getName() string {
	return "RandomBot"
}
func (b *RandomBot) startGame(marker string) {
	b.marker = marker
}

func (b *RandomBot) getMarker() string {
	return b.marker
}

func (b *RandomBot) step(board Board) (row, col int, marker string) {
	for {
		row, col, marker = rand.Intn(3), rand.Intn(3), b.marker
		if board[row][col] == " " {
			return
		}
	}
}
