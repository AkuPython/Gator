package main

import (
	"fmt"
	"os"

	"github.com/AkuPython/Gator/internal/config"
)

type state struct {
	conf config.Config
}

type command struct {
	Name string
	Args []string
}

func main() {
	conf, err := config.Read()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	
	cState := state{conf: conf}
	cCommands := commands{Command: make(map[string]func(*state, command) error)}
	cCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		fmt.Println("Usage: 'gator command <additional args>'")
		os.Exit(1)
	}
	err = cCommands.run(&cState, command{Name: os.Args[1], Args: os.Args[2:]})

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
