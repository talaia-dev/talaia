package hooks

import (
	"fmt"
	"os"
	"strings"

	"talaia.dev/internal/cli"
)

// Register wires the hooks module into the CLI app.
func Register(app *cli.App) {
	app.AddCommand(&cli.Command{
		Name:    "hooks",
		Summary: "Deny git writes in the agent's settings (claude, gemini) or in this clone's git hooks (git)",
		Usage:   "talaia hooks <agent>",
		Run: func(args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("hooks takes one agent: %s", strings.Join(Agents, ", "))
			}
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			rep, err := Run(root, args[0])
			if err != nil {
				return err
			}
			rep.Print()
			if args[0] == "git" {
				fmt.Println("hooks live in this clone only: run it again after every fresh clone")
			}
			return nil
		},
	})
}
