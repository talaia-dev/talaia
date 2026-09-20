// Package cli implements a small extensible command dispatcher.
//
// Feature modules register Commands (e.g. "talaia init") against the App;
// the dispatcher routes invocations accordingly.
package cli

import (
	"fmt"
	"sort"
)

// Command is a subcommand, e.g. "talaia init".
type Command struct {
	Name    string
	Summary string
	Usage   string
	Run     func(args []string) error
}

// App holds the registered commands and dispatches invocations.
type App struct {
	name     string
	commands map[string]*Command
}

// New creates an empty App. Modules add functionality via AddCommand before
// Run is called.
func New(name string) *App {
	a := &App{
		name:     name,
		commands: map[string]*Command{},
	}
	a.AddCommand(&Command{
		Name:    "help",
		Summary: "Show this help",
		Usage:   name + " help",
		Run: func([]string) error {
			a.printHelp()
			return nil
		},
	})
	return a
}

// AddCommand registers a subcommand.
func (a *App) AddCommand(c *Command) {
	a.commands[c.Name] = c
}

// Run dispatches the given arguments to a command.
func (a *App) Run(args []string) error {
	if len(args) == 0 {
		a.printHelp()
		return nil
	}
	if cmd, ok := a.commands[args[0]]; ok {
		return cmd.Run(args[1:])
	}
	return fmt.Errorf("unknown command %q, run '%s help'", args[0], a.name)
}

func (a *App) printHelp() {
	fmt.Printf("%s — human-in-command standard for work done with AI: the tool for software\n\n", a.name)

	fmt.Println("Commands:")
	names := make([]string, 0, len(a.commands))
	for n := range a.commands {
		names = append(names, n)
	}
	sort.Strings(names)
	for _, n := range names {
		fmt.Printf("  %-28s %s\n", a.commands[n].Usage, a.commands[n].Summary)
	}
}
