package hooks

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func readJSON(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var m map[string]any
	if err := json.Unmarshal(data, &m); err != nil {
		t.Fatalf("%s is not valid JSON: %v\n%s", path, err, data)
	}
	return m
}

func list(t *testing.T, m map[string]any, keys ...string) []string {
	t.Helper()
	var cur any = m
	for _, k := range keys {
		obj, ok := cur.(map[string]any)
		if !ok {
			t.Fatalf("no object at %v", keys)
		}
		cur = obj[k]
	}
	raw, ok := cur.([]any)
	if !ok {
		t.Fatalf("no list at %v: %v", keys, cur)
	}
	var out []string
	for _, v := range raw {
		out = append(out, v.(string))
	}
	return out
}

func TestClaudeWritesDenyRulesForEveryGitWrite(t *testing.T) {
	root := t.TempDir()
	rep, err := Run(root, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Created) != 1 || filepath.ToSlash(rep.Created[0]) != ".claude/settings.json" {
		t.Fatalf("created = %v", rep.Created)
	}
	deny := list(t, readJSON(t, filepath.Join(root, ".claude", "settings.json")), "permissions", "deny")
	for _, sub := range GitWrites {
		if !contains(deny, "Bash(git "+sub+" *)") {
			t.Errorf("deny lacks git %s: %v", sub, deny)
		}
	}
	if contains(deny, "Bash(git status *)") || contains(deny, "Bash(git *)") {
		t.Errorf("read-only git must stay allowed: %v", deny)
	}
}

func TestClaudeMergesIntoExistingSettings(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude")
	os.MkdirAll(dir, 0o755)
	existing := `{"model":"opus","permissions":{"allow":["Bash(npm test *)"],"deny":["Bash(git push *)","WebFetch"]}}`
	os.WriteFile(filepath.Join(dir, "settings.json"), []byte(existing), 0o644)

	rep, err := Run(root, "claude")
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Updated) != 1 {
		t.Fatalf("updated = %v, created = %v", rep.Updated, rep.Created)
	}
	m := readJSON(t, filepath.Join(dir, "settings.json"))
	if m["model"] != "opus" {
		t.Errorf("other keys must survive: %v", m)
	}
	if a := list(t, m, "permissions", "allow"); len(a) != 1 || a[0] != "Bash(npm test *)" {
		t.Errorf("allow must survive: %v", a)
	}
	deny := list(t, m, "permissions", "deny")
	if !contains(deny, "WebFetch") || count(deny, "Bash(git push *)") != 1 {
		t.Errorf("existing deny entries must survive, once: %v", deny)
	}
	if len(deny) != 2+len(GitWrites)-1 {
		t.Errorf("deny has %d entries, want %d: %v", len(deny), 2+len(GitWrites)-1, deny)
	}
}

func TestRunIsIdempotent(t *testing.T) {
	for _, agent := range []string{"claude", "gemini"} {
		root := t.TempDir()
		if _, err := Run(root, agent); err != nil {
			t.Fatal(err)
		}
		path := filepath.Join(root, "."+agent, "settings.json")
		before, _ := os.ReadFile(path)
		rep, err := Run(root, agent)
		if err != nil {
			t.Fatal(err)
		}
		after, _ := os.ReadFile(path)
		if len(rep.Kept) != 1 || len(rep.Created)+len(rep.Updated) != 0 || string(before) != string(after) {
			t.Errorf("%s: second run must change nothing: %+v", agent, rep)
		}
	}
}

func TestInvalidSettingsAreNotTouched(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, ".claude")
	os.MkdirAll(dir, 0o755)
	broken := []byte("{ not json")
	path := filepath.Join(dir, "settings.json")
	os.WriteFile(path, broken, 0o644)
	if _, err := Run(root, "claude"); err == nil {
		t.Fatal("invalid JSON must be an error")
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(broken) {
		t.Error("invalid file must be left exactly as found")
	}
}

func TestGeminiWritesExcludeForEveryGitWrite(t *testing.T) {
	root := t.TempDir()
	if _, err := Run(root, "gemini"); err != nil {
		t.Fatal(err)
	}
	m := readJSON(t, filepath.Join(root, ".gemini", "settings.json"))
	exclude := list(t, m, "tools", "exclude")
	for _, sub := range GitWrites {
		if !contains(exclude, "run_shell_command(git "+sub+")") {
			t.Errorf("exclude lacks git %s: %v", sub, exclude)
		}
	}
	if _, ok := m["tools"].(map[string]any)["core"]; ok {
		t.Error("tools.core must never be set: it turns every tool into an allowlist")
	}
}

func TestGitInstallsHooksInTheClone(t *testing.T) {
	root := t.TempDir()
	if _, err := Run(root, "git"); err == nil {
		t.Fatal("outside a git repository, hooks git must fail")
	}
	hooksDir := filepath.Join(root, ".git", "hooks")
	os.MkdirAll(hooksDir, 0o755)
	os.WriteFile(filepath.Join(hooksDir, "pre-push"), []byte("#!/bin/sh\n# husky\n"), 0o755)

	rep, err := Run(root, "git")
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Created) != 1 || len(rep.Kept) != 1 {
		t.Fatalf("report = %+v", rep)
	}
	info, err := os.Stat(filepath.Join(hooksDir, "reference-transaction"))
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode()&0o111 == 0 {
		t.Error("hook must be executable")
	}
	data, _ := os.ReadFile(filepath.Join(hooksDir, "reference-transaction"))
	for _, want := range []string{"CLAUDECODE", "GEMINI_CLI", "CODEX_SANDBOX", "TALAIA_AGENT", "refs/remotes/", "AGENTS.md"} {
		if !strings.Contains(string(data), want) {
			t.Errorf("reference-transaction lacks %q", want)
		}
	}
	kept, _ := os.ReadFile(filepath.Join(hooksDir, "pre-push"))
	if !strings.Contains(string(kept), "husky") {
		t.Error("an existing hook must never be overwritten")
	}
	if _, err := os.Stat(filepath.Join(root, ".talaia", "hooks")); !os.IsNotExist(err) {
		t.Error(".talaia/hooks must not be created")
	}
}

func TestGitHookBlocksAgentCommitsOnly(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	git := func(env []string, args ...string) error {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		cmd.Env = append(humanEnv(), env...)
		return cmd.Run()
	}
	for _, args := range [][]string{{"init", "-q"}, {"config", "user.email", "t@x"}, {"config", "user.name", "T"}} {
		if err := git(nil, args...); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Run(root, "git"); err != nil {
		t.Fatal(err)
	}
	os.WriteFile(filepath.Join(root, "a"), []byte("a\n"), 0o644)
	if err := git(nil, "add", "a"); err != nil {
		t.Fatal(err)
	}
	if err := git([]string{"CLAUDECODE=1"}, "commit", "-q", "-m", "by the agent"); err == nil {
		t.Error("an agent commit must be refused")
	}
	if err := git([]string{"TALAIA_AGENT=1"}, "commit", "-q", "--no-verify", "-m", "by the agent"); err == nil {
		t.Error("--no-verify must not skip the refusal")
	}
	if err := git(nil, "commit", "-q", "-m", "by the person"); err != nil {
		t.Errorf("the person's commit must pass: %v", err)
	}
	if err := git([]string{"CLAUDECODE=1"}, "log", "--oneline"); err != nil {
		t.Errorf("read-only git must pass for the agent: %v", err)
	}
}

func TestUnknownAgent(t *testing.T) {
	if _, err := Run(t.TempDir(), "cursor"); err == nil || !strings.Contains(err.Error(), "cursor") {
		t.Errorf("unknown agent must be a named error, got %v", err)
	}
}

// humanEnv is the test process environment without any agent marker,
// so the hook sees a person unless a test adds one back.
func humanEnv() []string {
	var env []string
	for _, kv := range os.Environ() {
		k := strings.SplitN(kv, "=", 2)[0]
		switch k {
		case "CLAUDECODE", "GEMINI_CLI", "CODEX_SANDBOX", "TALAIA_AGENT":
			continue
		}
		env = append(env, kv)
	}
	return env
}

func contains(l []string, s string) bool { return count(l, s) > 0 }

func count(l []string, s string) int {
	n := 0
	for _, v := range l {
		if v == s {
			n++
		}
	}
	return n
}
