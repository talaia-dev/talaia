// Package scaffold implements "talaia init": it generates the harness of
// the talaia framework in a project — the AGENTS.md rules that make an
// AI agent work inside Brief → Build → Watch → Sign, and the .talaia/
// context memory.
//
// Layout created (existing files are never overwritten):
//
//	AGENTS.md                        framework rules for AI agents (repo root)
//	CLAUDE.md                        redirect to AGENTS.md for Claude Code
//	GEMINI.md                        redirect to AGENTS.md for Gemini CLI
//	.github/PULL_REQUEST_TEMPLATE.md the explanation's shape and the Certificate of Understanding
//	.github/workflows/talaia.yml     the audit check on every pull request (the enforcement)
//	.talaia/context/architecture.md  project overview, structure, principles
//	.talaia/context/conventions.md   coding and development conventions
//	.talaia/context/constraints.md   hard limits, incl. decisions that bind future work
//	.talaia/.gitignore               keeps the agent's session note out of history
//
// context/ is the only memory: the story lives on the team's board, the
// change in git, the explanation in the pull request. Nothing here needs
// configuring after init: the agent reminds the certificate at sign-off
// (AGENTS.md) and the pull request template is applied by the forge.
// Gates — hooks, CI, branch protection — belong to the team.
//
// The content of every generated file lives as an editable template in the
// templates/ folder of this package, embedded into the binary at build time.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// templatesFS embeds the templates/ folder. Edit the files there to change
// what "talaia init" generates; rebuild to pick up the changes.
//
//go:embed templates
var templatesFS embed.FS

// tpl returns the content of a template file, panicking on a missing name:
// templates are embedded, so a failure here is a build defect, not a
// runtime condition.
func tpl(name string) string {
	data, err := fs.ReadFile(templatesFS, "templates/"+name)
	if err != nil {
		panic(fmt.Sprintf("scaffold: embedded template %q missing: %v", name, err))
	}
	return string(data)
}

// Report lists what Run did: files it created and files it left untouched
// because they already existed. Paths are relative to the project root.
type Report struct {
	Created []string
	Kept    []string
}

// Run generates the scaffolding under root. Existing files are never
// overwritten: they are reported in Report.Kept and left exactly as found.
func Run(root string) (*Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}

	rep := &Report{}
	create := func(relPath, content string) error {
		path := filepath.Join(root, relPath)
		if _, err := os.Lstat(path); err == nil {
			rep.Kept = append(rep.Kept, relPath)
			return nil
		} else if !os.IsNotExist(err) {
			return err
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			return err
		}
		rep.Created = append(rep.Created, relPath)
		return nil
	}

	for _, name := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"} {
		if err := create(name, tpl(name)); err != nil {
			return nil, err
		}
	}
	if err := os.MkdirAll(filepath.Join(root, ".github", "workflows"), 0o755); err != nil {
		return nil, err
	}
	if err := create(filepath.Join(".github", "PULL_REQUEST_TEMPLATE.md"), tpl("pull-request-template.md")); err != nil {
		return nil, err
	}
	if err := create(filepath.Join(".github", "workflows", "talaia.yml"), tpl("talaia-workflow.yml")); err != nil {
		return nil, err
	}

	if err := os.MkdirAll(filepath.Join(root, ".talaia", "context"), 0o755); err != nil {
		return nil, err
	}
	if err := create(filepath.Join(".talaia", ".gitignore"), tpl("talaia-gitignore")); err != nil {
		return nil, err
	}
	for _, name := range []string{"architecture.md", "conventions.md", "constraints.md"} {
		if err := create(filepath.Join(".talaia", "context", name), tpl("context/"+name)); err != nil {
			return nil, err
		}
	}
	return rep, nil
}

// Print writes a human-readable summary of the report to stdout.
func (r *Report) Print() {
	for _, p := range r.Created {
		fmt.Printf("  created  %s\n", p)
	}
	for _, p := range r.Kept {
		fmt.Printf("  kept     %s (already exists)\n", p)
	}
	fmt.Printf("%d created, %d kept\n", len(r.Created), len(r.Kept))
}
