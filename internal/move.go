package internal

import (
	"fmt"
)

func moveSnake(board [][]rune, snake *[]Coord, direction Direction, food Coord, quit chan struct{}) {
	oldHead := (*snake)[0]

	// Calculate new head coordinates
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

	// Insert new head (in-place modification)
	*snake = append(*snake, newHead) // Expand the slice
	copy((*snake)[1:], (*snake)[:])  // Shift elements down
	(*snake)[0] = newHead            // Set the new head

	if newHead.X < 0 || newHead.X >= len(board[0]) ||
		newHead.Y < 0 || newHead.Y >= len(board) {
		fmt.Println("Game Over! You hit the wall.")
		close(quit) // Signal to end the game
		return
	}

	// Check if food was eaten
	if newHead.X != food.X || newHead.Y != food.Y {
		// Remove tail if no food consumed (in-place modification)
		*snake = (*snake)[:len(*snake)-1]
	}
}

func (g *Game) isGameOver(board [][]rune, snake []Coord) bool {
	head := snake[0]

	// Wall collision (same check as in the updated moveSnake)
	if head.X < 0 || head.X >= len(board[0]) ||
		head.Y < 0 || head.Y >= len(board) {
		return true
	}

	// Self-collision (check if the head overlaps with the rest of the snake's body)
	for _, bodySegment := range snake[1:] {
		if head.X == bodySegment.X && head.Y == bodySegment.Y {
			return true
		}
	}

	return false // Game continues
}
