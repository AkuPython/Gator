package main

import "fmt"

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
	user := cmd.Args[0]
	err := s.conf.SetUser(user)
	if err != nil {
		return fmt.Errorf("Error setting username '%v' - %v", user, err)
	}
	fmt.Printf("Username set to: %v\n", user)
	return nil
}
