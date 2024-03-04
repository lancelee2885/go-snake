package internal

import (
	"fmt"
	"math/rand"
	"os"

	// "strings"
	"time"

	"github.com/gdamore/tcell/v2"
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Coord struct {
	X, Y int
}

type Game struct {
	Config GameConfig
	Screen tcell.Screen
	Dead   bool
}

type GameConfig struct {
	BoardWidth  int
	BoardHeight int
	SnakeSpawn  []Coord
	FoodSpawn   Coord
}

func NewGame(config GameConfig, screen tcell.Screen) *Game {
	return &Game{
		Config: config,
		Screen: screen,
	}
}

func (g *Game) Start() {
	if err := g.Screen.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing tcell: %v\n", err)
		return
	}
	defer g.Screen.Fini()

	// Game setup
	boardWidth, boardHeight := g.Config.BoardWidth, g.Config.BoardHeight
	board := make([][]rune, boardHeight)
	for i := range board {
		board[i] = make([]rune, boardWidth)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}

	snake := g.Config.SnakeSpawn
	direction := Up
	food := g.Config.FoodSpawn
	score := 0

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-quit:
				return
			default:
				g.processInput(&snake, &direction, quit)
			}
		}
	}()

	// Main game loop
	movementTicker := time.NewTicker(200 * time.Millisecond)
	gameOver := false
	for !gameOver {
		select {
		case <-movementTicker.C:
			moveSnake(board, &snake, direction, food, quit)
			if g.isGameOver(board, snake) {
				g.Screen.Clear()
				g.Screen.Show()
				gameOver = true
				break
			}
			if snake[0] == food {
				score++
				food = g.generateFood(board)
			}
			g.renderBoard(board, snake, food, score)
		case <-quit:
			movementTicker.Stop()
			return
		}
	}
	g.renderGameOver(score, quit)
}

func (g *Game) renderBoard(board [][]rune, snakeBody []Coord, food Coord, score int) {
	// Clear the screen
	g.Screen.Clear()

	// Board rows
	for i, row := range board {
		for j, cell := range row {
			style := tcell.StyleDefault // Default style
			if containsCoord(snakeBody, Coord{X: j, Y: i}) {
				style = style.Foreground(tcell.ColorGreen)
				body := 'x'
				if snakeBody[0].X == j && snakeBody[0].Y == i {
					style = style.Bold(true)
					body = 'X'
				}
				g.Screen.SetContent(j, i, body, nil, style)
			} else if food.X == j && food.Y == i {
				style = style.Foreground(tcell.ColorRed)
				g.Screen.SetContent(j, i, '*', nil, style)
			} else {
				g.Screen.SetContent(j, i, cell, nil, style)
			}
		}
	}

	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	width := len(board[0])
	for x := 0; x < width+2; x++ {
		g.Screen.SetContent(x, 0, tcell.RuneBlock, nil, style)            // Top border
		g.Screen.SetContent(x, len(board)+1, tcell.RuneBlock, nil, style) // Bottom border
	}

	for y := 1; y < len(board)+1; y++ {
		g.Screen.SetContent(0, y, tcell.RuneVLine, nil, style)       // Left border
		g.Screen.SetContent(width+1, y, tcell.RuneVLine, nil, style) // Right border
	}

	// Render score
	scoreText := fmt.Sprintf("Score: %d", score)
	for i, r := range scoreText {
		g.Screen.SetContent(i+1, len(board)+2, r, nil, style)
	}

	g.Screen.SetStyle(tcell.StyleDefault)

	g.Screen.Show()
}

func (g *Game) renderGameOver(score int, quit chan struct{}) {

	width, height := g.Config.BoardWidth, g.Config.BoardHeight

	gameOverText := "GAME OVER"
	scoreText := fmt.Sprintf("Final Score: %d", score)

	// Calculate positions for centered text
	textX := width/2 - len(gameOverText)/2
	textY := height/2 - 2

	// Render text with styling
	style := tcell.StyleDefault.Bold(true).Foreground(tcell.ColorRed)
	for _, r := range gameOverText {
		g.Screen.SetContent(textX, textY, r, nil, style)
		textX++
	}
	textX = width/2 - len(scoreText)/2
	g.Screen.SetContent(textX, textY+2, ' ', nil, tcell.StyleDefault) // Reset to default style
	for _, r := range scoreText {
		g.Screen.SetContent(textX, textY+2, r, nil, tcell.StyleDefault)
		textX++
	}
	pressAnyKeyText := "Press any key to restart"
	textX = width/2 - len(pressAnyKeyText)/2
	g.Screen.SetContent(textX, textY+4, ' ', nil, tcell.StyleDefault) // Reset to default style
	for _, r := range pressAnyKeyText {
		g.Screen.SetContent(textX, textY+4, r, nil, tcell.StyleDefault)
		textX++
	}

	g.Screen.Show()

	// Wait for key press to restart
	ev := g.Screen.PollEvent()
	for ev != nil {
		if ev, ok := ev.(*tcell.EventKey); ok {
			if ev.Key() != tcell.KeyEscape { // Restart on any key except Escape
				g.restartGame(quit)
				return
			}
		}
		ev = g.Screen.PollEvent()
	}
}

func (g *Game) restartGame(quit chan struct{}) {
	// Signal all goroutines to stop
	close(quit)

	// Properly stop and clear the screen
	g.Screen.Clear()
	g.Screen.Fini() // Properly finalize the screen

	// Reinitialize the game state and start a new game
	newGame := NewGame(g.Config, g.Screen)
	newGame.Start()
}

func (g *Game) resetBoard() [][]rune {
	boardWidth, boardHeight := g.Config.BoardWidth, g.Config.BoardHeight
	board := make([][]rune, boardHeight)
	for i := range board {
		board[i] = make([]rune, boardWidth)
		for j := range board[i] {
			board[i][j] = ' '
		}
	}
	return board
}

func containsCoord(coords []Coord, c Coord) bool {
	for _, coord := range coords {
		if coord.X == c.X && coord.Y == c.Y {
			return true
		}
	}
	return false
}

func (g *Game) generateFood(board [][]rune) Coord {
	boardWidth, boardHeight := len(board[0]), len(board)
	for {
		food := Coord{X: randInt(1, boardWidth-2), Y: randInt(1, boardHeight-2)}
		if board[food.Y][food.X] == ' ' {
			return food
		}
	}
}

func randInt(min, max int) int {
	return min + rand.Intn(max-min)
}
