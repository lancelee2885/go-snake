package internal

import (
// "fmt"
)

func (g *Game) moveSnake(board [][]rune, snake *[]Coord, direction Direction, food Coord, quit chan struct{}) bool {
	oldHead := (*snake)[0]
	var newHead Coord
	switch direction {
	case Up:
		newHead = Coord{X: oldHead.X, Y: oldHead.Y - 1}
	case Down:
		newHead = Coord{X: oldHead.X, Y: oldHead.Y + 1}
	case Left:
		newHead = Coord{X: oldHead.X - 1, Y: oldHead.Y}
	case Right:
		newHead = Coord{X: oldHead.X + 1, Y: oldHead.Y}
	}

	if newHead.X < 0 || newHead.X >= len(board[0]) || newHead.Y < 0 || newHead.Y >= len(board) {
		close(quit)
		return true
	}

	// Check if the new head collides with the snake's body
	for _, segment := range (*snake)[1:] {
		if newHead.X == segment.X && newHead.Y == segment.Y {
			close(quit)
			return true
		}
	}

	// Update snake's body
	*snake = append([]Coord{newHead}, (*snake)...)
	if newHead != food {
		*snake = (*snake)[:len(*snake)-1]
	}

	return false
}

func (g *Game) isGameOver(board [][]rune, snake []Coord) bool {
	if len(snake) == 0 {
		return false
	}

	head := snake[0]
	if head.X < 0 || head.X >= len(board[0]) || head.Y < 0 || head.Y >= len(board) {
		return true
	}
	for _, bodySegment := range snake[1:] {
		if head.X == bodySegment.X && head.Y == bodySegment.Y {
			return true
		}
	}
	return false
}
