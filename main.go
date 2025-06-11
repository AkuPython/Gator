package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/AkuPython/Gator/internal/config"
	"github.com/AkuPython/Gator/internal/database"
	_ "github.com/lib/pq"
)

type state struct {
	db *database.Queries
	cfg *config.Config
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
	
	db, err := sql.Open("postgres", conf.DbURL)
	if err != nil {
		fmt.Println("Could not connect to DB:", err)
		os.Exit(1)
	}
	
	dbQueries := database.New(db)
	
	cState := state{cfg: &conf, db: dbQueries}
	cCommands := commands{Command: make(map[string]func(*state, command) error)}

	// REGISTER COMMANDS
	cCommands.register("login", handlerLogin)
	cCommands.register("register", handlerRegister)
	cCommands.register("reset", handlerReset)
	cCommands.register("users", handlerGetUsers)
	cCommands.register("agg", handlerAgg)
	cCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	cCommands.register("feeds", handlerGetFeeds)
	cCommands.register("follow", middlewareLoggedIn(handlerCreateFeedFollow))
	cCommands.register("following", middlewareLoggedIn(handlerFeedFollowsForUser))


	// -----------------
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
