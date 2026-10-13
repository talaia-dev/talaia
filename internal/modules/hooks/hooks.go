// Package hooks implements "talaia hooks <agent>": it makes the
// framework's git rule — no state-changing git command by the agent —
// mechanical, in the place each agent obeys.
//
//	talaia hooks claude   .claude/settings.json   permissions.deny, one Bash rule per git write
//	talaia hooks gemini   .gemini/settings.json   tools.exclude, one run_shell_command per git write
//	talaia hooks git      .git/hooks/             reference-transaction + pre-push, refuse when an
//	                                              agent's environment variable is set; per clone
//
// The settings files are versioned and shared; hooks edits them in place
// — only the list it owns, keeping everything else — and refuses to
// touch a file it cannot parse. The git hooks live in the clone: never
// overwritten, never versioned, and they hold only for the agents they
// recognise. None of this is a security boundary; it is the rule made
// mechanical for the agent that forgets it.
package hooks

import (
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
)

// GitWrites are the git subcommands the framework reserves to the
// person: everything AGENTS.md lists, plus their equivalents.
var GitWrites = []string{
	"add", "am", "apply", "branch", "checkout", "cherry-pick", "clean",
	"commit", "merge", "mv", "pull", "push", "rebase", "reset", "restore",
	"revert", "rm", "stash", "switch", "tag",
}

// Agents lists the supported targets, in the order help prints them.
var Agents = []string{"claude", "gemini", "git"}

//go:embed templates
var templatesFS embed.FS

// Report lists what Run did, with paths relative to the project root.
type Report struct {
	Created []string
	Updated []string
	Kept    []string
}

// Run applies the git rule for one agent to the project at root.
func Run(root, agent string) (*Report, error) {
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	rep := &Report{}
	switch agent {
	case "claude":
		rules := make([]string, len(GitWrites))
		for i, sub := range GitWrites {
			rules[i] = "Bash(git " + sub + " *)"
		}
		err = rep.mergeList(root, filepath.Join(".claude", "settings.json"), []string{"permissions", "deny"}, rules)
	case "gemini":
		rules := make([]string, len(GitWrites))
		for i, sub := range GitWrites {
			rules[i] = "run_shell_command(git " + sub + ")"
		}
		err = rep.mergeList(root, filepath.Join(".gemini", "settings.json"), []string{"tools", "exclude"}, rules)
	case "git":
		err = rep.installGitHooks(root)
	default:
		return nil, fmt.Errorf("hooks: unknown agent %q — supported: claude, gemini, git (no other agent takes a per-command deny list from the repository; AGENTS.md still applies)", agent)
	}
	if err != nil {
		return nil, err
	}
	return rep, nil
}

// mergeList adds entries to the string list at keys inside the JSON
// object file relPath, creating the file or the path as needed and
// keeping everything else as found. Duplicates are not added; an
// unparseable file is an error and is left untouched.
func (r *Report) mergeList(root, relPath string, keys []string, entries []string) error {
	path := filepath.Join(root, relPath)
	doc := map[string]any{}
	exists := true
	data, err := os.ReadFile(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		exists = false
	case err != nil:
		return err
	default:
		if err := json.Unmarshal(data, &doc); err != nil {
			return fmt.Errorf("%s: not a JSON object, left untouched: %v", relPath, err)
		}
	}

	cur := doc
	for _, k := range keys[:len(keys)-1] {
		next, ok := cur[k].(map[string]any)
		if !ok {
			if cur[k] != nil {
				return fmt.Errorf("%s: %q is not an object, left untouched", relPath, k)
			}
			next = map[string]any{}
			cur[k] = next
		}
		cur = next
	}
	last := keys[len(keys)-1]
	list, ok := cur[last].([]any)
	if !ok && cur[last] != nil {
		return fmt.Errorf("%s: %q is not a list, left untouched", relPath, last)
	}
	present := map[string]bool{}
	for _, v := range list {
		if s, ok := v.(string); ok {
			present[s] = true
		}
	}
	added := 0
	for _, e := range entries {
		if !present[e] {
			list = append(list, e)
			present[e] = true
			added++
		}
	}
	if exists && added == 0 {
		r.Kept = append(r.Kept, relPath)
		return nil
	}
	cur[last] = list

	out, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, append(out, '\n'), 0o644); err != nil {
		return err
	}
	if exists {
		r.Updated = append(r.Updated, relPath)
	} else {
		r.Created = append(r.Created, relPath)
	}
	return nil
}

// installGitHooks writes the two hook scripts into the clone's
// .git/hooks, leaving any existing hook exactly as found.
func (r *Report) installGitHooks(root string) error {
	hooksDir := filepath.Join(root, ".git", "hooks")
	if info, err := os.Stat(filepath.Join(root, ".git")); err != nil || !info.IsDir() {
		return fmt.Errorf("hooks git: %s is not the root of a git repository (.git/ not found)", root)
	}
	if err := os.MkdirAll(hooksDir, 0o755); err != nil {
		return err
	}
	names, err := fs.ReadDir(templatesFS, "templates")
	if err != nil {
		return err
	}
	sort.Slice(names, func(i, j int) bool { return names[i].Name() < names[j].Name() })
	for _, entry := range names {
		relPath := filepath.Join(".git", "hooks", entry.Name())
		path := filepath.Join(root, relPath)
		if _, err := os.Lstat(path); err == nil {
			r.Kept = append(r.Kept, relPath)
			continue
		} else if !errors.Is(err, fs.ErrNotExist) {
			return err
		}
		content, err := fs.ReadFile(templatesFS, "templates/"+entry.Name())
		if err != nil {
			return err
		}
		if err := os.WriteFile(path, content, 0o755); err != nil {
			return err
		}
		r.Created = append(r.Created, relPath)
	}
	return nil
}

// Print writes a human-readable summary of the report to stdout.
func (r *Report) Print() {
	for _, p := range r.Created {
		fmt.Printf("  created  %s\n", p)
	}
	for _, p := range r.Updated {
		fmt.Printf("  updated  %s\n", p)
	}
	for _, p := range r.Kept {
		fmt.Printf("  kept     %s (already there)\n", p)
	}
}
