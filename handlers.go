package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/AkuPython/Gator/internal/database"
	"github.com/google/uuid"
)

type commands struct {
	Command map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	f, ok := c.Command[cmd.Name]
	if !ok {
		return fmt.Errorf("Command Does not Exist")
	}
	return f(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.Command[name] = f
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide (only) Username")
	}
	username := cmd.Args[0]

	user, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		log.Fatalf("Non-Existant Username '%v' - %v", username, err)
	}

	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting username '%v' - %v", username, err)
	}
	fmt.Printf("Username set to: %v\n", username)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("Must provide (only) Username")
	}
	username := cmd.Args[0]
	user, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: username,
	})
	if err != nil {
		log.Fatalf("Error creating user with username '%v' - %v", username, err)
	}
	fmt.Printf("Username: '%v' created \n", username)
	err = s.cfg.SetUser(user.Name)
	if err != nil {
		return fmt.Errorf("Error setting username '%v' - %v", username, err)
	}
	fmt.Printf("Username set to: %v\n", username)
	log.Print(user)
	return nil
}

func handlerReset(s *state, cmd command) error {
	return s.db.DeleteUsers(context.Background())
}

func handlerGetUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		log.Fatalf("Could not get users from DB - %v", err)
	}
	for _, user := range users {
		if user.Name == s.cfg.CurrentUserName {
			fmt.Println("*", user.Name, "(current)")
		} else {
			fmt.Println("*", user.Name)
		}
	}
	return nil
}
