package audit

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fixture mimics `git log --no-merges --shortstat --format=<RS>%H<US>%an<US>%ae<US>%s<US>%b`.
const fixture = "\x1e" + "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2" + "\x1f" + "Roger" + "\x1f" + "r@x" + "\x1f" +
	"feat: export orders as CSV (#42)" + "\x1f" +
	"Range boundaries are inclusive.\n\nAssisted-by: Claude Code\nSigned-off-by: Roger <r@x>\n" +
	"\n 3 files changed, 120 insertions(+), 4 deletions(-)\n" +
	"\x1e" + "b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3" + "\x1f" + "Ana" + "\x1f" + "a@x" + "\x1f" +
	"fix typo" + "\x1f" + "\n 1 file changed, 1 insertion(+), 1 deletion(-)\n" +
	"\x1e" + "c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4" + "\x1f" + "Roger" + "\x1f" + "r@x" + "\x1f" +
	"refactor: export module (#43)" + "\x1f" +
	"Assisted-by: Codex\nReviewed-by: Ana <a@x>\nSigned-off-by: Roger <r@x>\n" +
	"\n 6 files changed, 500 insertions(+), 210 deletions(-)\n" +
	"\x1e" + "d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5" + "\x1f" + "Ana" + "\x1f" + "a@x" + "\x1f" +
	"Revert \"feat: export orders as CSV (#42)\"" + "\x1f" +
	"This reverts commit a1b2c3d.\n\nSigned-off-by: Ana <a@x>\n" +
	"\n 3 files changed, 4 insertions(+), 120 deletions(-)\n"

func TestParseReadsCommitsAndSignals(t *testing.T) {
	commits := Parse(fixture)
	if len(commits) != 4 {
		t.Fatalf("parsed %d commits, want 4", len(commits))
	}
	c := commits[0]
	if c.Hash != "a1b2c3d4e5f6a1b2c3d4e5f6a1b2c3d4e5f6a1b2" || c.Author != "Roger" || c.Subject != "feat: export orders as CSV (#42)" {
		t.Errorf("first commit fields wrong: %+v", c)
	}
	if !c.Signed || !c.Story || len(c.AssistedBy) != 1 || c.AssistedBy[0] != "Claude Code" || c.Reviewed || c.Revert {
		t.Errorf("first commit signals wrong: %+v", c)
	}
	if c.Insertions != 120 || c.Deletions != 4 || c.Lines() != 124 {
		t.Errorf("first commit stats wrong: +%d -%d", c.Insertions, c.Deletions)
	}
	if strings.Contains(c.Body, "files changed") {
		t.Errorf("shortstat line must be stripped from the body: %q", c.Body)
	}
	if u := commits[1]; u.Signed || u.Story || len(u.AssistedBy) != 0 || u.Lines() != 2 {
		t.Errorf("unsigned commit signals wrong: %+v", u)
	}
	if r := commits[2]; !r.Reviewed || r.AssistedBy[0] != "Codex" || r.Lines() != 710 {
		t.Errorf("reviewed commit signals wrong: %+v", r)
	}
	if v := commits[3]; !v.Revert || !v.Signed || !v.Story {
		t.Errorf("revert commit signals wrong: %+v", v)
	}
	if got := Parse(""); len(got) != 0 {
		t.Errorf("empty log parsed %d commits", len(got))
	}
}

func TestSummarizeCountsAndGaps(t *testing.T) {
	r := Summarize(Parse(fixture), 300)
	got := []int{r.Total, r.Signed, r.Story, r.Assisted, r.AssistedReviewed, r.WithinBudget, r.Reverts}
	want := []int{4, 3, 3, 2, 1, 3, 1}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("summary counts = %v, want %v (total, signed, story, assisted, assistedReviewed, withinBudget, reverts)", got, want)
			break
		}
	}
	kinds := map[string]int{}
	for _, g := range r.Gaps {
		kinds[g.Kind]++
	}
	if kinds["unsigned"] != 1 || kinds["no story"] != 1 || kinds["710 lines"] != 1 {
		t.Errorf("gaps = %+v", r.Gaps)
	}
	if !r.HasGaps() {
		t.Error("a report with an unsigned commit must have gaps")
	}
	if len(r.Authors) != 2 || r.Authors[0].Name != "Roger" || r.Authors[0].Commits != 2 || r.Authors[0].Assisted != 2 || r.Authors[0].OverBudget != 1 {
		t.Errorf("authors = %+v", r.Authors)
	}
	if a := r.Authors[1]; a.Name != "Ana" || a.Signed != 1 || a.Story != 1 {
		t.Errorf("second author = %+v", a)
	}
	if e := Summarize(nil, 300); e.Total != 0 || e.HasGaps() {
		t.Errorf("empty summary = %+v", e)
	}
}

func TestWriteTextReport(t *testing.T) {
	var out bytes.Buffer
	r := Summarize(Parse(fixture), 300)
	r.Range = "HEAD"
	WriteText(&out, r, true)
	text := out.String()
	for _, want := range []string{
		"4 commits", "signed (Signed-off-by)", "3/4", "75%",
		"AI-assisted (Assisted-by)", "of which Reviewed-by", "1/2",
		"within budget", "300", "reverts", "gaps", "b2c3d4e", "unsigned", "fix typo",
		"by author", "Roger", "Ana",
	} {
		if !strings.Contains(text, want) {
			t.Errorf("text report lacks %q\n%s", want, text)
		}
	}
	out.Reset()
	WriteText(&out, r, false)
	if strings.Contains(out.String(), "by author") {
		t.Error("authors section must be omitted when disabled")
	}
	out.Reset()
	if err := WriteJSON(&out, r); err != nil {
		t.Fatal(err)
	}
	var back Report
	if err := json.Unmarshal(out.Bytes(), &back); err != nil {
		t.Fatalf("json report does not round-trip: %v", err)
	}
	if back.Total != 4 || len(back.Gaps) != 3 {
		t.Errorf("json report = %+v", back)
	}
}

func TestFlagsParseOptions(t *testing.T) {
	o, err := ParseArgs([]string{"--budget", "150", "--strict", "--json", "--no-authors", "main..feature", "--since=2026-01-01"})
	if err != nil {
		t.Fatal(err)
	}
	if o.Budget != 150 || !o.Strict || !o.JSON || o.Authors || strings.Join(o.LogArgs, " ") != "main..feature --since=2026-01-01" {
		t.Errorf("options = %+v", o)
	}
	d, err := ParseArgs(nil)
	if err != nil || d.Budget != 300 || d.Strict || d.JSON || !d.Authors || len(d.LogArgs) != 0 {
		t.Errorf("defaults = %+v (%v)", d, err)
	}
	if _, err := ParseArgs([]string{"--budget", "x"}); err == nil {
		t.Error("a non-numeric budget must be rejected")
	}
}

// TestRunOnARealRepository builds a throwaway repository and audits it.
func TestRunOnARealRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not installed")
	}
	root := t.TempDir()
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		cmd.Env = append(os.Environ(), "GIT_AUTHOR_NAME=Dev", "GIT_AUTHOR_EMAIL=dev@x", "GIT_COMMITTER_NAME=Dev", "GIT_COMMITTER_EMAIL=dev@x")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	write := func(name, content string) {
		t.Helper()
		if err := os.WriteFile(filepath.Join(root, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	git("init", "-q", ".")
	write("a.txt", "a\n")
	git("add", "a.txt")
	git("commit", "-q", "-m", "feat: a (#1)\n\nAssisted-by: Claude Code\nSigned-off-by: Dev <dev@x>")
	write("b.txt", "b\n")
	git("add", "b.txt")
	git("commit", "-q", "-m", "fix typo")
	write("big.txt", strings.Repeat("line\n", 400))
	git("add", "big.txt")
	git("commit", "-q", "-m", "feat: big (#2)\n\nSigned-off-by: Dev <dev@x>")

	var out bytes.Buffer
	rep, err := Run(root, Options{Budget: 300, Authors: true}, &out)
	if err != nil {
		t.Fatalf("Run: %v\n%s", err, out.String())
	}
	if rep.Total != 3 || rep.Signed != 2 || rep.Story != 2 || rep.Assisted != 1 || rep.WithinBudget != 2 {
		t.Errorf("report = %+v", rep)
	}
	if !strings.Contains(out.String(), "fix typo") {
		t.Errorf("gap list should name the unsigned commit:\n%s", out.String())
	}
	if _, err := Run(t.TempDir(), Options{Budget: 300}, &out); err == nil {
		t.Error("auditing a directory that is not a repository must fail")
	}
}
