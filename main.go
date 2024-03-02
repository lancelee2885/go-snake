/*
Copyright © 2024 Lance Lee <lancelee2885@gmail.com>
*/
package main

import (
	"os"

	"github.com/lancelee2885/go-snake/cmd"
)

func main() {
	if err := cmd.Root().Execute(); err != nil {
		os.Exit(1)
	}
}
