package scaffold

import (
	"fmt"
	"os"

	"talaia.dev/internal/cli"
)

// Register wires the scaffold module into the CLI app.
func Register(app *cli.App) {
	app.AddCommand(&cli.Command{
		Name:    "init",
		Summary: "Set up the talaia harness: AGENTS.md, .talaia/context/, the pull request template and the check (never overwrites)",
		Usage:   "talaia init",
		Run: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("init takes no arguments")
			}
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			rep, err := Run(root)
			if err != nil {
				return err
			}
			rep.Print()
			return nil
		},
	})
}
