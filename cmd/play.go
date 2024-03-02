package cmd

import (
	"github.com/lancelee2885/go-snake/internal/basic"
	"github.com/spf13/cobra"
)

func play() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "play",
		Short: "Play a task",
		Run: func(cmd *cobra.Command, args []string) {
			basic.NewGame()
		},
	}

	return cmd
}
