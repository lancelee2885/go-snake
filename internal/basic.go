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
	for {
		select {
		case <-movementTicker.C:
			moveSnake(&snake, direction, food)
			if snake[0] == food {
				score++
				food = generateFood(board)
			}
			g.renderBoard(board, snake, food, score)
			g.Screen.Show()
		case <-quit:
			movementTicker.Stop()
			return
		}
	}
}

func (g *Game) renderBoard(board [][]rune, snakeBody []Coord, food Coord, score int) {
	// Clear the screen
	g.Screen.Clear()

	// Board rows
	for i, row := range board {
		for j, cell := range row {
			style := tcell.StyleDefault // Default style
			if containsCoord(snakeBody, Coord{X: j, Y: i}) {
				style = style.Foreground(tcell.ColorGreen) // Example: Snake in green
				g.Screen.SetContent(j, i, 'X', nil, style)
			} else if food.X == j && food.Y == i {
				style = style.Foreground(tcell.ColorRed) // Example: Food in red
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

	g.Screen.SetStyle(tcell.StyleDefault)
}

func containsCoord(coords []Coord, c Coord) bool {
	for _, coord := range coords {
		if coord.X == c.X && coord.Y == c.Y {
			return true
		}
	}
	return false
}

func generateFood(board [][]rune) Coord {
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
