package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "%s\n", "linkview: expected single argument")
		return
	}

	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		return
	}

	links, err := FindLinks(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		return
	}

	terminal := NewTerminal(&links)

	err = terminal.Start()
	if err != nil {
		panic(err)
	}
	defer terminal.Close()

	terminal.Render()

	for {
		go PollEvent()
		quit, err := terminal.HandleEvent(<-EventChan)
		if err != nil {
			fmt.Fprintf(os.Stderr, "linkview: %s\n", err)
		}

		if quit {
			return
		}
	}
}
