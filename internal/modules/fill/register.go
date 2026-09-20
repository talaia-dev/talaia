package fill

import (
	"fmt"
	"os"

	"talaia.dev/internal/cli"
)

// Register wires the fill module into the CLI app.
func Register(app *cli.App) {
	app.AddCommand(&cli.Command{
		Name:    "fill",
		Summary: "Print the prompt that asks your AI agent to populate .talaia/context/ from the code",
		Usage:   "talaia fill  (e.g. talaia fill | claude -p)",
		Run: func(args []string) error {
			if len(args) != 0 {
				return fmt.Errorf("fill takes no arguments")
			}
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			return Run(root, os.Stdout)
		},
	})
}
