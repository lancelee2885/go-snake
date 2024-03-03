package internal

import ()

func moveSnake(snake *[]Coord, direction Direction, food Coord) {
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

	// Check if food was eaten
	if newHead.X != food.X || newHead.Y != food.Y {
		// Remove tail if no food consumed (in-place modification)
		*snake = (*snake)[:len(*snake)-1]
	}
}
