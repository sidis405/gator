package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/sidis405/gator/internal/config"
)

type state struct {
	c *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	list map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	command, ok := c.list[cmd.name]
	if !ok {
		return errors.New("command does not exist")
	}
	return command(s, cmd)
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.list[name] = f
}

func main() {
	c, err := config.Read()
	if err != nil {
		fmt.Printf("cannot read config file: %q", err)
		return
	}
	s := state{c: &c}
	cmds := commands{list: map[string]func(*state, command) error{
		"login": handlerLogin,
	}}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("not enough arguments provided")
		os.Exit(1)
	}
	cmdName := args[1]
	otherParams := args[2:]
	err = cmds.run(&s, command{
		name:      cmdName,
		arguments: otherParams,
	})

	if err != nil {
		fmt.Printf("%q", err)
		os.Exit(1)
	}
}

func handlerLogin(s *state, cmd command) error {
	if len(cmd.arguments) == 0 {
		return errors.New("username is required")
	}
	userName := cmd.arguments[0]
	err := s.c.SetUser(userName)
	if err != nil {
		return err
	}
	fmt.Println("The user has been set to", userName)
	return nil
}
