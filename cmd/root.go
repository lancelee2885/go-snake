/*
Copyright © 2024 Lance Lee <lancelee2885@gmail.com>
*/
package cmd

import (
	"github.com/spf13/cobra"
)

func Root() *cobra.Command {

	root := &cobra.Command{
		Use: "snake-game",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			cmd.SilenceUsage = true
		},
	}

	root.AddCommand(
		play(),
	)

	return root
}
