package internal

import (
	"os"

	"github.com/gdamore/tcell/v2"
)

func (g *Game) processInput(snake *[]Coord, direction *Direction, quit chan struct{}) {
	ev := g.Screen.PollEvent()
	switch ev := ev.(type) {
	case *tcell.EventKey:
		switch ev.Key() {
		case tcell.KeyEscape:
			select {
			case quit <- struct{}{}:
			default:
				g.Screen.Fini()
				close(quit)
				os.Exit(0)
			}
		case tcell.KeyUp:
			if *direction != Down {
				*direction = Up
			}
		case tcell.KeyDown:
			if *direction != Up {
				*direction = Down
			}
		case tcell.KeyLeft:
			if *direction != Right {
				*direction = Left
			}
		case tcell.KeyRight:
			if *direction != Left {
				*direction = Right
			}
		}
	}
}
