package cmd

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lancelee2885/go-snake/internal"
	"github.com/spf13/cobra"
)

func play() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "play",
		Short: "Play a task",
		RunE: func(cmd *cobra.Command, args []string) error {
			screen, err := tcell.NewScreen()
			if err != nil {
				return err
			}

			boardConfig := internal.GameConfig{
				BoardWidth:  20,
				BoardHeight: 15,
				SnakeSpawn:  []internal.Coord{{X: 10, Y: 7}, {X: 10, Y: 8}, {X: 10, Y: 9}},
				FoodSpawn:   internal.Coord{X: 5, Y: 5},
			}
			internal.NewGame(boardConfig, screen).Start()

			return nil
		},
	}

	return cmd
}
