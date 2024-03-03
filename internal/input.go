package internal

import (
	"github.com/gdamore/tcell/v2"
)

func (g *Game) processInput(snake *[]Coord, direction *Direction, quit chan struct{}) {
	ev := g.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			close(quit)
		case tcell.KeyUp:
			*direction = Up
		case tcell.KeyDown:
			*direction = Down
		case tcell.KeyLeft:
			*direction = Left
		case tcell.KeyRight:
			*direction = Right
		}
	}
}
