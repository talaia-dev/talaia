// Package audit implements "talaia audit": it reads git history and
// reports whether the framework's signatures are there — Signed-off-by,
// a story id, Assisted-by, Reviewed-by, the size budget, reverts. Read
// only: one `git log`, parsed here; no file is written.
//
//	talaia audit                       everything reachable from HEAD
//	talaia audit main..feature         any git log range or option
//	talaia audit --budget 200 --strict --json --no-authors
package audit

import (
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Commit is one non-merge commit with the signals the framework cares
// about, derived from its message and its shortstat.
type Commit struct {
	Hash       string   `json:"hash"`
	Author     string   `json:"author"`
	Email      string   `json:"email"`
	Subject    string   `json:"subject"`
	Body       string   `json:"-"`
	Insertions int      `json:"insertions"`
	Deletions  int      `json:"deletions"`
	Signed     bool     `json:"signed"`
	Story      bool     `json:"story"`
	AssistedBy []string `json:"assistedBy,omitempty"`
	Reviewed   bool     `json:"reviewed"`
	Revert     bool     `json:"revert"`
}

// Lines is the number of changed lines: insertions plus deletions.
func (c Commit) Lines() int { return c.Insertions + c.Deletions }

// Format is the git log format Parse understands: records separated by
// RS (0x1e), fields by US (0x1f); --shortstat appends its line to the
// body.
const Format = "%x1e%H%x1f%an%x1f%ae%x1f%s%x1f%b"

var (
	statRe     = regexp.MustCompile(`(?m)^\s*(\d+) files? changed(?:, (\d+) insertions?\(\+\))?(?:, (\d+) deletions?\(-\))?\s*$`)
	storyRe    = regexp.MustCompile(`(^|[^A-Za-z0-9_/&])#\d+`)
	signedRe   = regexp.MustCompile(`(?m)^Signed-off-by:\s*\S`)
	reviewedRe = regexp.MustCompile(`(?m)^Reviewed-by:\s*\S`)
	assistedRe = regexp.MustCompile(`(?m)^Assisted-by:\s*(.+?)\s*$`)
)

// Parse turns the output of `git log --shortstat --format=Format` into
// commits, classifying each one.
func Parse(log string) []Commit {
	var commits []Commit
	for _, rec := range strings.Split(log, "\x1e") {
		if strings.TrimSpace(rec) == "" {
			continue
		}
		f := strings.SplitN(rec, "\x1f", 5)
		if len(f) < 4 {
			continue
		}
		c := Commit{Hash: strings.TrimSpace(f[0]), Author: f[1], Email: f[2], Subject: strings.TrimSpace(f[3])}
		body := ""
		if len(f) == 5 {
			body = f[4]
		}
		if m := statRe.FindStringSubmatch(body); m != nil {
			c.Insertions, _ = strconv.Atoi(m[2])
			c.Deletions, _ = strconv.Atoi(m[3])
			body = statRe.ReplaceAllString(body, "")
		}
		c.Body = strings.TrimSpace(body)
		text := c.Subject + "\n" + c.Body
		c.Signed = signedRe.MatchString(c.Body)
		c.Reviewed = reviewedRe.MatchString(c.Body)
		c.Story = storyRe.MatchString(text)
		c.Revert = strings.HasPrefix(c.Subject, "Revert ")
		for _, m := range assistedRe.FindAllStringSubmatch(c.Body, -1) {
			c.AssistedBy = append(c.AssistedBy, m[1])
		}
		commits = append(commits, c)
	}
	return commits
}

// Gap is a commit that misses a signal: "unsigned", "no story" or
// "<n> lines" when over the budget.
type Gap struct {
	Hash    string `json:"hash"`
	Kind    string `json:"kind"`
	Subject string `json:"subject"`
}

// AuthorStat is the per-person breakdown.
type AuthorStat struct {
	Name       string `json:"name"`
	Commits    int    `json:"commits"`
	Signed     int    `json:"signed"`
	Story      int    `json:"story"`
	Assisted   int    `json:"assisted"`
	OverBudget int    `json:"overBudget"`
}

// Report is what audit prints.
type Report struct {
	Range            string       `json:"range"`
	Budget           int          `json:"budget"`
	Total            int          `json:"total"`
	Signed           int          `json:"signed"`
	Story            int          `json:"story"`
	Assisted         int          `json:"assisted"`
	AssistedReviewed int          `json:"assistedReviewed"`
	WithinBudget     int          `json:"withinBudget"`
	Reverts          int          `json:"reverts"`
	Gaps             []Gap        `json:"gaps"`
	Authors          []AuthorStat `json:"authors"`
}

// HasGaps reports whether any commit is unsigned or has no story — the
// two things --strict refuses. The budget is advisory.
func (r Report) HasGaps() bool {
	return r.Signed < r.Total || r.Story < r.Total
}

// Summarize computes the report for the given commits and line budget.
func Summarize(commits []Commit, budget int) Report {
	r := Report{Budget: budget, Total: len(commits)}
	byAuthor := map[string]*AuthorStat{}
	var order []string
	for _, c := range commits {
		a := byAuthor[c.Author]
		if a == nil {
			a = &AuthorStat{Name: c.Author}
			byAuthor[c.Author] = a
			order = append(order, c.Author)
		}
		a.Commits++
		if c.Signed {
			r.Signed++
			a.Signed++
		} else {
			r.Gaps = append(r.Gaps, Gap{c.Hash, "unsigned", c.Subject})
		}
		if c.Story {
			r.Story++
			a.Story++
		} else {
			r.Gaps = append(r.Gaps, Gap{c.Hash, "no story", c.Subject})
		}
		if len(c.AssistedBy) > 0 {
			r.Assisted++
			a.Assisted++
			if c.Reviewed {
				r.AssistedReviewed++
			}
		}
		if c.Lines() <= budget {
			r.WithinBudget++
		} else {
			a.OverBudget++
			r.Gaps = append(r.Gaps, Gap{c.Hash, fmt.Sprintf("%d lines", c.Lines()), c.Subject})
		}
		if c.Revert {
			r.Reverts++
		}
	}
	for _, name := range order {
		r.Authors = append(r.Authors, *byAuthor[name])
	}
	sort.SliceStable(r.Authors, func(i, j int) bool { return r.Authors[i].Commits > r.Authors[j].Commits })
	return r
}

func pct(n, of int) string {
	if of == 0 {
		return "  –"
	}
	return fmt.Sprintf("%3d%%", n*100/of)
}

// WriteText prints the human-readable report.
func WriteText(w io.Writer, r Report, authors bool) {
	rng := r.Range
	if rng == "" {
		rng = "HEAD"
	}
	fmt.Fprintf(w, "talaia audit — %d commits (%s, no merges)\n\n", r.Total, rng)
	row := func(label string, n, of int) {
		fmt.Fprintf(w, "  %-30s %4d/%-4d %s\n", label, n, of, pct(n, of))
	}
	row("signed (Signed-off-by)", r.Signed, r.Total)
	row("story referenced (#id)", r.Story, r.Total)
	row("AI-assisted (Assisted-by)", r.Assisted, r.Total)
	row("  …of which Reviewed-by", r.AssistedReviewed, r.Assisted)
	row(fmt.Sprintf("within budget (≤%d lines)", r.Budget), r.WithinBudget, r.Total)
	row("reverts", r.Reverts, r.Total)
	if len(r.Gaps) > 0 {
		fmt.Fprintf(w, "\n  gaps\n")
		for _, g := range r.Gaps {
			fmt.Fprintf(w, "  %-8s %-20s %s\n", short(g.Hash), g.Kind, g.Subject)
		}
	}
	if authors && len(r.Authors) > 0 {
		fmt.Fprintf(w, "\n  by author\n")
		for _, a := range r.Authors {
			fmt.Fprintf(w, "  %-24s %4d   signed %-4d story %-4d assisted %-4d over budget %d\n",
				a.Name, a.Commits, a.Signed, a.Story, a.Assisted, a.OverBudget)
		}
	}
}

// WriteJSON prints the report as JSON.
func WriteJSON(w io.Writer, r Report) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(r)
}

func short(hash string) string {
	if len(hash) > 7 {
		return hash[:7]
	}
	return hash
}

// Options are the command's settings.
type Options struct {
	Budget  int
	Strict  bool
	JSON    bool
	Authors bool
	LogArgs []string // passed to git log verbatim: ranges, --since, paths…
}

// ParseArgs reads the command line. Unknown arguments go to git log.
func ParseArgs(args []string) (Options, error) {
	o := Options{Budget: 300, Authors: true}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--budget":
			if i+1 >= len(args) {
				return o, fmt.Errorf("--budget needs a number")
			}
			i++
			n, err := strconv.Atoi(args[i])
			if err != nil || n <= 0 {
				return o, fmt.Errorf("--budget: %q is not a positive number", args[i])
			}
			o.Budget = n
		case strings.HasPrefix(a, "--budget="):
			n, err := strconv.Atoi(strings.TrimPrefix(a, "--budget="))
			if err != nil || n <= 0 {
				return o, fmt.Errorf("--budget: %q is not a positive number", a)
			}
			o.Budget = n
		case a == "--strict":
			o.Strict = true
		case a == "--json":
			o.JSON = true
		case a == "--no-authors":
			o.Authors = false
		default:
			o.LogArgs = append(o.LogArgs, a)
		}
	}
	return o, nil
}

// Run audits the repository at root and writes the report to w. It
// returns the report so callers can decide on --strict.
func Run(root string, o Options, w io.Writer) (Report, error) {
	args := append([]string{"-C", root, "log", "--no-merges", "--shortstat", "--format=" + Format}, o.LogArgs...)
	out, err := exec.Command("git", args...).Output()
	if err != nil {
		msg := err.Error()
		if ee, ok := err.(*exec.ExitError); ok && len(ee.Stderr) > 0 {
			msg = strings.TrimSpace(string(ee.Stderr))
		}
		return Report{}, fmt.Errorf("git log: %s", msg)
	}
	r := Summarize(Parse(string(out)), o.Budget)
	r.Range = strings.Join(o.LogArgs, " ")
	if r.Range == "" {
		r.Range = "HEAD"
	}
	if o.JSON {
		if err := WriteJSON(w, r); err != nil {
			return r, err
		}
	} else {
		WriteText(w, r, o.Authors)
	}
	return r, nil
}
