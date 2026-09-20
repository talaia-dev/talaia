package fill

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunRequiresTalaiaFolder(t *testing.T) {
	root := t.TempDir()
	var out strings.Builder
	err := Run(root, &out)
	if err == nil {
		t.Fatal("Run succeeded without .talaia/, want error")
	}
	if !strings.Contains(err.Error(), "talaia init") {
		t.Errorf("error %q does not suggest 'talaia init'", err)
	}
	if out.Len() != 0 {
		t.Errorf("Run wrote output despite the error: %q", out.String())
	}
}

func TestRunEmitsThePrompt(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, ".talaia"), 0o755); err != nil {
		t.Fatal(err)
	}
	var out strings.Builder
	if err := Run(root, &out); err != nil {
		t.Fatal(err)
	}
	prompt := out.String()
	for _, want := range []string{
		"file by file",
		".talaia/context/architecture.md",
		".talaia/context/conventions.md",
		".talaia/context/constraints.md",
		"Do not modify anything outside",
		"do not invent",
	} {
		if !strings.Contains(prompt, want) {
			t.Errorf("prompt lacks %q", want)
		}
	}
	for _, banned := range []string{".talaia/decisions/", ".talaia/features/", "ADR-"} {
		if strings.Contains(prompt, banned) {
			t.Errorf("prompt must not ask for %q: context/ is the only memory", banned)
		}
	}
}
