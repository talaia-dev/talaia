// Command talaia is the entry point of the Talaia tool.
//
// Talaia sets up the human-in-charge framework in a software project:
// an AGENTS.md at the repo root with the rules that make AI coding agents
// work inside Brief → Build → Watch → Sign (plus CLAUDE.md/GEMINI.md
// redirects), a ".talaia" folder that holds the project's design — the
// files a person writes before any Brief — and the pull request
// template; and it audits git
// history for the framework's signatures. Future capabilities can be
// plugged in as modules.
package main

import (
	"fmt"
	"os"

	"talaia.dev/internal/cli"
	"talaia.dev/internal/modules/audit"
	"talaia.dev/internal/modules/hooks"
	"talaia.dev/internal/modules/scaffold"
)

func main() {
	app := cli.New("talaia")

	// Register feature modules. Future capabilities register themselves
	// here the same way.
	scaffold.Register(app)
	audit.Register(app)
	hooks.Register(app)

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "talaia: %v\n", err)
		os.Exit(1)
	}
}
