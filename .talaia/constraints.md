# constraints.md — Technical, business and architectural constraints

<!-- Design (.talaia/). Written and approved by a person: the AI reads
     it and never changes it on its own. Hard limits that bound
     implementation choices: an approach that violates one of these is
     wrong even if it works. -->

- **The design is the person's; the AI never changes it on its own**
  (the Developer, 2026-10-09). These three files are the truths of
  the system. The agent stops — and starts no Brief — when the design
  is missing, when an ask contradicts it, or when work shows it wrong
  or incomplete; it goes on only when the person has decided, by
  editing the files or by dictating the change. The agent writes here
  only at the person's dictation: wording and format, never the
  content. Signing the design (`git commit -s`) is good practice,
  with no check. There is no agent memory: `talaia fill` and the
  MEMORY section of AGENTS.md were removed for this reason — do not
  bring back any way for the AI to fill or update these files by
  itself.

- **Standard library only.** `go.mod` declares zero dependencies;
  `README.md` states "No external dependencies: standard library
  only, fully offline". Adding a dependency violates a stated design
  point (2026-08: recorded as ADR-002 in the legacy `decisions/`).
- **Go 1.23+** (`go.mod`, `README.md`). Build-time requirement only;
  the produced binary is standalone.
- **`talaia init` must never overwrite or delete existing files** —
  idempotent and non-destructive, enforced by tests
  (`scaffold_test.go`: `TestRunIsIdempotent`,
  `TestRunPreservesEditedFiles`). It only writes at the target repo
  root, never in subdirectories (`TestRunTouchesOnlyTheRoot`).
- **Talaia executes no AI and makes no network calls.** `audit` runs
  git **read-only** (`git log`), the one external process the tool
  spawns; it must never grow a write. (An "executes AI via workers"
  revocation lived for a few hours on 2026-09-01 and died with the
  worktable — see below.)
- **Rules live in AGENTS.md only.** `CLAUDE.md`/`GEMINI.md` are thin
  redirects (`@AGENTS.md` import + hyperlink); tests forbid drift
  (`TestRedirectFilesPointToAgentsMD`).
- **`Signed-off-by` is deliberately additive to the DCO** (2026-09-02):
  talaia reuses the trailer, `-s` and the required-check model; where
  a project also uses the DCO, one sign-off certifies both —
  understanding on top of origin, never replacing it. Stated in the
  README's Sign section and in AGENTS.md; the PR template's explicit
  clauses disambiguate sign-offs made out of DCO habit. Never mint a
  custom trailer: it would lose `-s` automation and twenty years of
  ecosystem weight.
- **No planning/approval workflow may be reintroduced.**
  `TestAgentsMDDescribesTheFramework` explicitly fails if the
  generated AGENTS.md mentions `journal.json` or "approved plan"
  (legacy ADR-001, 2026-08-11: the v0 approval workflow added
  ceremony that real AI-assisted development routed around).
- **The framework generates no files a human must read** (2026-08-28; hardened
  2026-09-01). `init` creates no `decisions/`, `features/` or `context/`
  (`TestRunCreatesOnlyTheDesign`); the generated AGENTS.md bans
  documents outright ("Write no documents") and must mention neither
  `delivery-note` nor `.talaia/work`
  (`TestAgentsMDDescribesTheFramework`). This lesson has now been
  paid **four times** — feature files/ADRs (Aug), the v0 approval
  plans, the MWI worktable + delivery note (built and deleted
  within 2026-09-01), and the session note `.talaia/delivery-note.md`
  (restored 09-02, deleted for good 2026-09-29 with its
  `.talaia/.gitignore`: nobody read it and it did not help
  understanding): an artifact that requires reading to have value
  has no value; even talaia's own creator would not read them.
  Visibility = checks that run on git + surfaces the human already
  uses (IDE diff, commit, PR). Never generate documentation as a
  feature again; tests forbid the session note in AGENTS.md and the
  PR template and forbid `init` from creating `.talaia/.gitignore`.
- **Rules live in AGENTS.md, and the repo's own AGENTS.md is a copy of
  the template** (`internal/modules/scaffold/templates/AGENTS.md`):
  change the template, copy it over — talaia is developed with its
  own framework.
- **The README is universal: no explicit religious references**
  (2026-08-28). The foundation — the opening of `## What we hold` (under the `# Manifesto` part),
  there is no separate section — ("the human person comes first";
  Babel vs the coastal *talaies*, which stand in for Nehemiah's wall;
  machines cannot answer; restraint) is
  stated as talaia's own conviction, without naming a church, a pope,
  a faith or a document, and without quotations or paragraph numbers
  **in the body**. Since 2026-09-20 (the Developer's ask) the Babel
  paragraph carries one footnote naming the source: Leo XIV's
  encyclical *Magnifica Humanitas* (2026; Babel/Nehemiah in ¶7–9,
  the "not yes or no to technology" choice in ¶9; also ¶97–111,
  ¶140, ¶148–164) with the vatican.va link — the only religious
  reference in the README, and it stays in the footnote. Since
  2026-10-01 `## What AI is` restates ¶98–100 and ¶104
  under the same rule (no name, no marker). Changes to
  the manifesto are measured against the foundation.
- **The agent writes no files outside code.** The
  Developer writes the commit message and the PR body in their own
  words (a note-as-commit-message experiment on 2026-09-01 was
  reverted the same day). The certificate reminder is zero-config:
  the agent prints the Sign block at the end of finished work; the
  forge applies the PR template (What / Why / How it was tested with
  evidence + the certificate). `talaia init` never touches `.git/` or
  the root `.gitignore`; only the explicit `talaia hooks git` writes
  into `.git/hooks/` (2026-09-30), and it says it is per clone.
- **Enforcing the git rule is layered, and none of it is a security
  boundary** (2026-09-30, verified in a lab repo, git 2.43):
  `reference-transaction` stops every ref update even under
  `--no-verify` but not the index/working tree (`add`, `reset
  --hard` still resets, `checkout -- f`, `clean`, `switch`);
  `pre-push` is skipped by `--no-verify`; env-var detection is
  bypassable by `env -u`. TTY detection was rejected: GUI git clients
  have no TTY and would block the person. Do not present any layer
  as a guarantee; the guarantee is the person's name on the commit.
- **Talaia is a standard plus its check, not a review product**
  (decided 2026-09-01, after studying CodeRabbit's "agentic change
  management"): the value is the convention (Signed-off-by as
  Certificate of Understanding + Assisted-by + story id + PR
  template), the installer (`init`), and enforcement. **The check
  exists** (2026-09-02): `action.yml` at the repo root — composite
  GitHub Action, `uses: talaia-dev/talaia@<ref>`, builds the binary
  from `github.action_path` with setup-go and runs `audit --strict`
  on the PR's base..head — plus `.github/workflows/ci.yml` (Go test/
  vet/fmt on push+PR; the audit check on PRs, dogfooding). As a local
  report audit stays thin on purpose. Do NOT build review UIs,
  boards, portals or approval steps — review tools (CodeRabbit et
  al.) scale the reading and are named as complements in the README;
  talaia certifies that the reading happened. The Watch's surface is
  the IDE diff; staging is the first acceptance.
- **The README's `# Framework` part is medium-agnostic** (2026-09-14;
  tool-agnostic since 2026-08-29): it speaks of "the instrument" and
  "the one who asks / the one who answers", names git only as the
  software example inside a rule, and nothing else — no `talaia
  init`/`hooks`/`audit`, no `.talaia/` paths, no `AGENTS.md`, no `gh`.
  Everything git-shaped lives in `# In practice › Software` and
  `# Using talaia`.
- **The Certificate of Understanding is the centre and is versioned**
  (2026-09-14, reshaped 2026-09-16; since 2026-10-05 the centre of
  the *Understanding* part only — the manifesto and its values come
  first, see `architecture.md`): **two clauses**, by the
  Developer's decision — (a) I made it with the help of an AI
  (disclosure; `Assisted-by:` names the tool) and (b) I have read it
  in full, I understand it, I take responsibility for it, and I make
  it my own (certification). Directed / explain / checked survive as
  the *reading* of the clauses in `What each clause asks`, never as
  clauses of their own — do not re-split them. The text keeps the
  **DCO's shape** (version header, `I certify that:`, lettered
  clauses, DCO register: "in whole or in part"; "to a brief", the
  native idiom, not "on a brief" — 2026-09-17) by the Developer's
  ask — but **no copyright line, no verbatim-copies rule, no
  disclaimers** (removed 2026-09-17 on the Developer's ask: "no te
  pases con el copyright y otros disclaimers"). Keep that shape and
  that restraint on any rewording.
  Text fixed as the ```txt fence directly under `# The Certificate of
  Understanding` (no `## The text` heading since 2026-09-17), canonical
  URL `talaia.dev/certificate` (a real page,
  `web/src/pages/certificate.astro`, that shows the text in the DCO's
  format and nothing else — it reads the README's only ```txt fence at
  build time, so the README is the one place the text is written and
  a second ```txt fence would break the page; the earlier Astro
  redirect to the part's anchor is gone). The certificate is invoked
  by **one line of two parts: the address and a name** —
  `talaia.dev/certificate — John Doe`; in git, `-s` writes the name
  and the address goes in the message. Never a text block or a
  second address line. **Which AI is evidence, not certification**
  (Developer's decision, 2026-09-17, after two tries): clause (a)
  says "an AI", the address already carries it, and a signature that
  names a product dates itself — so the tool goes only where the
  medium has room (`Assisted-by:` trailer in git, the cover note of
  a report), as a courtesy. Ruled out the same day: `?model=` on the
  address (reads as a parameter, not a signature) and a third part
  ", assisted by Claude Code" in the line (redundant with (a)).
  Example names in the README are placeholders (John Doe, Jane Roe);
  the Developer's name appears only in the site footer. Changing the
  wording means a new version, never a silent edit — one exception,
  the Developer's on 2026-10-05, while nobody had signed yet: (b)
  gained "its opinions and assumptions are mine" (from `What AI is`:
  AI is not neutral) in place of "or, where there is nothing to
  check, that it says what I meant", still Version 1. A third clause (c) about the
  work itself (serves people; AI not presented as alive; no machine
  in the place of a person's judgment) was written and removed the
  same day: too long, and it made the certificate three clauses.
  Instead the text ends with one unlettered closing line, "I sign
  this in keeping with the talaia manifesto: talaia.dev" — the manifesto's
  values enter by reference, never as clauses; a clause "follows
  the manifesto" was ruled out (it would certify against an
  unversioned text). The closing line lives only in the README
  text, not in the Sign block or the PR template.
  The README is the
  standard for three media (published content, documents and
  reports, software); only software has a machine check. Aligned
  2026-09-18: the scaffold `AGENTS.md` (and this repo's, kept
  identical), the PR template and their tests carry the two-clause
  certificate, the address in the commit message, and the
  framework's Brief/Build rules (read the ask against the project,
  no design decisions, seven Build rules incl. hand-made content).
- **Never `Co-Authored-By:` for an AI** (2026-09-29, README "Against
  Co-Authorship", names Claude Code's default trailer): a machine
  cannot be held to its work, so it is recorded as `Assisted-by:`,
  never as an author.
- **Git writes are reserved to the developer** in any
  Talaia-managed project, including this one (`AGENTS.md`).
- Performance, security or compliance requirements: no evidence in
  the repository beyond the offline/no-network posture above.
- Platform targets: none stated beyond Go itself; tests use
  `filepath.ToSlash`, suggesting portability across OSes is intended
  but no explicit platform matrix exists.
