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
	food := g.generateFood(board)
	score := 0

	quit := make(chan struct{})
	done := make(chan struct{})
	go func() {
		defer close(done)
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
			gameOver = g.moveSnake(board, &snake, direction, food, quit)
			if gameOver {
				break
			}
			if snake[0] == food {
				score++
				food = g.generateFood(board)
			}
			g.renderBoard(board, snake, food, score)
		case <-quit:
			gameOver = true
		}
	}
	movementTicker.Stop()
	<-done

	// Render game over screen
	g.renderGameOver(score)
}

func (g *Game) renderBoard(board [][]rune, snakeBody []Coord, food Coord, score int) {
	g.Screen.Clear()

	// Render snake
	for i, segment := range snakeBody {
		style := tcell.StyleDefault.Foreground(tcell.ColorGreen)
		body := 'x'
		if i == 0 {
			style = style.Bold(true)
			body = 'X'
		}
		g.Screen.SetContent(segment.X, segment.Y, body, nil, style)
	}

	// Render food
	g.Screen.SetContent(food.X, food.Y, '*', nil, tcell.StyleDefault.Foreground(tcell.ColorRed))

	// Render borders
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	width := len(board[0])
	for x := 0; x < width+2; x++ {
		g.Screen.SetContent(x, 0, tcell.RuneBlock, nil, style)
		g.Screen.SetContent(x, len(board)+1, tcell.RuneBlock, nil, style)
	}
	for y := 1; y < len(board)+1; y++ {
		g.Screen.SetContent(0, y, tcell.RuneVLine, nil, style)
		g.Screen.SetContent(width+1, y, tcell.RuneVLine, nil, style)
	}

	// Render score
	scoreText := fmt.Sprintf("Score: %d", score)
	for i, r := range scoreText {
		g.Screen.SetContent(i+1, len(board)+2, r, nil, style)
	}

	g.Screen.Show()
}

func (g *Game) renderGameOver(score int) {
	width, height := g.Config.BoardWidth, g.Config.BoardHeight

	gameOverText := "GAME OVER"
	scoreText := fmt.Sprintf("Final Score: %d", score)
	pressAnyKeyText := "Press any key to restart (or Esc to exit)"

	// Calculate positions for centered text
	gameOverX := width/2 - len(gameOverText)/2
	gameOverY := height/2 - 2
	scoreX := width/2 - len(scoreText)/2
	scoreY := gameOverY + 2
	pressAnyKeyX := width/2 - len(pressAnyKeyText)/2
	pressAnyKeyY := scoreY + 2

	// Render game over text
	style := tcell.StyleDefault.Bold(true).Foreground(tcell.ColorRed)
	for _, r := range gameOverText {
		g.Screen.SetContent(gameOverX, gameOverY, r, nil, style)
		gameOverX++
	}

	// Render score text
	scoreStyle := tcell.StyleDefault
	for _, r := range scoreText {
		g.Screen.SetContent(scoreX, scoreY, r, nil, scoreStyle)
		scoreX++
	}

	// Render "Press any key to restart" text
	pressAnyKeyStyle := tcell.StyleDefault
	for _, r := range pressAnyKeyText {
		g.Screen.SetContent(pressAnyKeyX, pressAnyKeyY, r, nil, pressAnyKeyStyle)
		pressAnyKeyX++
	}

	g.Screen.Show()

	// Wait for a key press to restart or exit the game
	for {
		ev := g.Screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape {
				g.Screen.Fini()
				os.Exit(0)
			} else {
				g.restartGame()
				return
			}
		}
	}
}

func (g *Game) restartGame() {
	g.Screen.Fini()
	g.Dead = true

	newScreen, err := tcell.NewScreen()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating new screen: %v\n", err)
		return
	}

	newGame := NewGame(g.Config, newScreen)
	newGame.Start()
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
