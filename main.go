package main

import (
	"fmt"
	"os"

	"github.com/Bieri2004/gator/internal/config"
)

// State struct - hält die Anwendungsdaten
type state struct {
	cfg *config.Config
}

// Command struct - repräsentiert einen CLI-Befehl
type command struct {
	name string
	args []string
}

// Commands struct - hält alle registrierten Command-Handler
type commands struct {
	handlers map[string]func(*state, command) error
}

// register Methode - registriert einen neuen Command-Handler
func (c *commands) register(name string, f func(*state, command) error) {
	c.handlers[name] = f
}

// run Methode - führt einen Command aus
func (c *commands) run(s *state, cmd command) error {
	handler, exists := c.handlers[cmd.name]
	if !exists {
		return fmt.Errorf("unknown command: %s", cmd.name)
	}
	return handler(s, cmd)
}

// handlerLogin - Handler für den login Command
func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("the login handler expects a single argument, the username")
	}

	username := cmd.args[0]

	err := s.cfg.SetUser(username)
	if err != nil {
		return err
	}

	fmt.Printf("User has been set to: %s\n", username)
	return nil
}

func main() {
	// Config lesen
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}

	// State erstellen
	appState := &state{
		cfg: &cfg,
	}

	// Commands-Struct mit initialisierter Map erstellen
	cmds := &commands{
		handlers: make(map[string]func(*state, command) error),
	}

	// Login-Command registrieren
	cmds.register("login", handlerLogin)

	// Command-Line-Argumente holen
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Error: not enough arguments provided")
		os.Exit(1)
	}

	// Command-Struct erstellen
	cmd := command{
		name: args[1],
		args: args[2:],
	}

	// Command ausführen
	err = cmds.run(appState, cmd)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
