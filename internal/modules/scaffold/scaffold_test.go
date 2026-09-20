package scaffold

import (
	"crypto/sha256"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// newProject builds a fixture tree and returns its root. The subdirectories
// verify that init only touches the root.
func newProject(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	for _, d := range []string{"cmd/app", "internal/core", ".git/objects"} {
		if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	return root
}

func TestRunCreatesScaffolding(t *testing.T) {
	root := newProject(t)
	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}

	wantCreated := []string{
		".github/PULL_REQUEST_TEMPLATE.md",
		".github/workflows/talaia.yml",
		".talaia/.gitignore",
		".talaia/context/architecture.md",
		".talaia/context/constraints.md",
		".talaia/context/conventions.md",
		"AGENTS.md",
		"CLAUDE.md",
		"GEMINI.md",
	}
	got := append([]string(nil), rep.Created...)
	for i, p := range got {
		got[i] = filepath.ToSlash(p)
	}
	sort.Strings(got)
	if len(got) != len(wantCreated) {
		t.Fatalf("created = %v, want %v", got, wantCreated)
	}
	for i := range got {
		if got[i] != wantCreated[i] {
			t.Fatalf("created = %v, want %v", got, wantCreated)
		}
	}
	if len(rep.Kept) != 0 {
		t.Fatalf("kept = %v, want none", rep.Kept)
	}
}

func TestAgentsMDDescribesTheFramework(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)
	for _, want := range []string{
		"Brief",
		"Build",
		"Watch",
		"Sign",
		"Never draft the commit message",
		"Test first",
		"Write no documents",
		"Two strikes",
		"Assisted-by",
		".talaia/context/",
		"If `.talaia/disabled` exists, ignore this file",
		"Sign only if you can certify",
		"--trailer \"Assisted-by:",
		"Never run a state-changing git command",
		"## ALWAYS",
		"talaia.dev/certificate",
		"Nothing more, nothing less",
		"What a person wrote stays as they wrote it",
		"No questions left: build",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("AGENTS.md lacks %q", want)
		}
	}
	for _, banned := range []string{
		"journal.json", "approved plan", "decisions/", "features/",
		".talaia/work", "human Developer", "the Developer's", "Developer says",
		"delivery-note", "session note", "everything is clear", "Three questions",
	} {
		if strings.Contains(content, banned) {
			t.Errorf("AGENTS.md must not reference %q", banned)
		}
	}
}

func TestPullRequestTemplateCarriesTheNoteAndTheCertificate(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".github", "PULL_REQUEST_TEMPLATE.md"))
	if err != nil {
		t.Fatal(err)
	}
	tpl := string(data)
	for _, want := range []string{
		"## What changed", "## Why", "## How it was tested",
		"## Three questions",
		"Assisted by:",
		"talaia.dev/certificate",
		"- [ ] This change was created with the assistance of an AI, at my direction and to a brief that was mine",
		"- [ ] I have read it in full; I can explain what it does and why without the help of an AI; I checked what it claims; I take responsibility for it and make it my own",
	} {
		if !strings.Contains(tpl, want) {
			t.Errorf("PULL_REQUEST_TEMPLATE.md lacks %q", want)
		}
	}
	ign, err := os.ReadFile(filepath.Join(root, ".talaia", ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(ign), "delivery-note.md") {
		t.Error(".talaia/.gitignore must keep the session note out of history")
	}
	for _, gone := range []string{".gitmessage", ".talaia/hooks", ".talaia/work"} {
		if _, err := os.Stat(filepath.Join(root, gone)); !os.IsNotExist(err) {
			t.Errorf("%s must not be created: gates belong to the team, not the framework", gone)
		}
	}
}

func TestWorkflowInstallsTheAuditCheck(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(filepath.Join(root, ".github", "workflows", "talaia.yml"))
	if err != nil {
		t.Fatal(err)
	}
	wf := string(data)
	for _, want := range []string{
		"on: pull_request",
		"uses: talaia-dev/talaia@",
		"fetch-depth: 0",
	} {
		if !strings.Contains(wf, want) {
			t.Errorf("talaia.yml lacks %q", want)
		}
	}
}

func TestRunCreatesOnlyContextMemory(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"decisions", "features"} {
		if _, err := os.Stat(filepath.Join(root, ".talaia", d)); !os.IsNotExist(err) {
			t.Errorf(".talaia/%s must not be created: the story lives on the board, the change in git", d)
		}
	}
}

func TestRunTouchesOnlyTheRoot(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	for _, d := range []string{"cmd", "cmd/app", "internal", "internal/core", ".git"} {
		for _, name := range []string{"AGENTS.md", "CLAUDE.md", "GEMINI.md"} {
			p := filepath.Join(root, d, name)
			if _, err := os.Stat(p); !os.IsNotExist(err) {
				t.Errorf("%s must not exist in %s", name, d)
			}
		}
	}
}

func TestRedirectFilesPointToAgentsMD(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"CLAUDE.md", "GEMINI.md"} {
		data, err := os.ReadFile(filepath.Join(root, name))
		if err != nil {
			t.Fatal(err)
		}
		content := string(data)
		if !strings.Contains(content, "@AGENTS.md") {
			t.Errorf("%s lacks the @AGENTS.md import line", name)
		}
		if !strings.Contains(content, "[AGENTS.md](AGENTS.md)") {
			t.Errorf("%s lacks a hyperlink to AGENTS.md", name)
		}
	}
}

func TestRunIsIdempotent(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	before := hashTree(t, root)

	rep, err := Run(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Created) != 0 {
		t.Fatalf("second run created %v, want none", rep.Created)
	}
	if len(rep.Kept) == 0 {
		t.Fatal("second run kept nothing, want all existing files reported")
	}
	after := hashTree(t, root)
	for p, h := range before {
		if after[p] != h {
			t.Errorf("file %s changed on second run", p)
		}
	}
	if len(after) != len(before) {
		t.Errorf("file count changed: %d -> %d", len(before), len(after))
	}
}

func TestRunPreservesEditedFiles(t *testing.T) {
	root := newProject(t)
	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	edited := []byte("# My own rules\n\n- never touch prod\n")
	target := filepath.Join(root, "AGENTS.md")
	if err := os.WriteFile(target, edited, 0o644); err != nil {
		t.Fatal(err)
	}

	if _, err := Run(root); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(edited) {
		t.Fatal("edited AGENTS.md was overwritten by init")
	}
}

// hashTree returns sha256 digests for every regular file under root.
func hashTree(t *testing.T, root string) map[string][32]byte {
	t.Helper()
	hashes := map[string][32]byte{}
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		hashes[path] = sha256.Sum256(data)
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return hashes
}
