// Package fill implements "talaia fill": it emits the prompt that asks
// an AI agent to populate the context memory (.talaia/context/) by
// analyzing the repository file by file.
//
// Talaia executes no AI itself: pipe or paste the prompt into the agent
// of your choice, e.g.
//
//	talaia fill | claude -p
//
// The prompt lives as an editable template in the templates/ folder of
// this package, embedded into the binary at build time.
package fill

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

// templatesFS embeds the templates/ folder. Edit the files there to
// change the emitted prompt; rebuild to pick up the changes.
//
//go:embed templates
var templatesFS embed.FS

// tpl returns the content of a template file, panicking on a missing
// name: templates are embedded, so a failure here is a build defect.
func tpl(name string) string {
	data, err := fs.ReadFile(templatesFS, "templates/"+name)
	if err != nil {
		panic(fmt.Sprintf("fill: embedded template %q missing: %v", name, err))
	}
	return string(data)
}

// Run writes the fill prompt to w. It refuses to run when root has no
// .talaia/ folder: the project memory must be scaffolded first.
func Run(root string, w io.Writer) error {
	if _, err := os.Stat(filepath.Join(root, ".talaia")); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf(".talaia/ not found: run 'talaia init' first")
		}
		return err
	}
	_, err := io.WriteString(w, tpl("fill.prompt.md"))
	return err
}
