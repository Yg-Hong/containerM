package cmd

import (
	"fmt"
	"os"

	"containerM/container"
)

func Execute() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: containerM run <command>")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "run":
		container.Run(os.Args[2:])
	case "child":
		container.Child(os.Args[2:])
	default:
		panic("unknown command")
	}
}
