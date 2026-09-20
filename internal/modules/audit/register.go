package audit

import (
	"fmt"
	"os"

	"talaia.dev/internal/cli"
)

// Register wires the audit module into the CLI app.
func Register(app *cli.App) {
	app.AddCommand(&cli.Command{
		Name:    "audit",
		Summary: "Report signatures, story ids, Assisted-by and size budget from git history (read-only)",
		Usage:   "talaia audit [range] [--budget N] [--strict] [--json] [--no-authors]",
		Run: func(args []string) error {
			o, err := ParseArgs(args)
			if err != nil {
				return err
			}
			root, err := os.Getwd()
			if err != nil {
				return err
			}
			rep, err := Run(root, o, os.Stdout)
			if err != nil {
				return err
			}
			if o.Strict && rep.HasGaps() {
				return fmt.Errorf("audit: %d unsigned, %d without a story (--strict)", rep.Total-rep.Signed, rep.Total-rep.Story)
			}
			return nil
		},
	})
}
