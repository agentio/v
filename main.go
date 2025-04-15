package main

import (
	"os"

	"github.com/agentio/v/cmd"
)

func main() {
	if err := cmd.Cmd().Execute(); err != nil {
		os.Exit(1)
	}
}
