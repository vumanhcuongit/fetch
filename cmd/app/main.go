package main

import (
	"fetch-go/cmd/commands"
	"os"
)

func main() {
	err := commands.Execute()
	if err != nil {
		os.Exit(1)
	}
}
