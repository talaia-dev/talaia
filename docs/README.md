# Manifesto

## Introduction

Generative AI is faster at producing answers than humans are at thinking them through. When generating becomes cheaper than thinking, the path of least resistance takes over: the AI proposes, the human accepts, and speed starts replacing thought.

That shift creates risks that are easy to miss:

- **Judgment drifts away from the person.** The faster the machine can produce, the easier it becomes to stop questioning, exploring and deciding deliberately.
- **Accountability dissolves.** Work goes out that nobody fully understands — code, a report, a post. When it fails, "the AI wrote it" explains nothing — a machine cannot own an outcome.
- **Understanding erodes.** Tests pass, drafts go unread, and people sign off on what they cannot explain.
- **Design erodes.** AI is good at producing locally plausible pieces, but local plausibility is not global coherence. Without someone holding the intent, structure starts being accumulated rather than designed.
- **Plausibility replaces value.** Generated work often looks complete and correct, but looking right is not the same as being useful, true or worth making.
- **Existing controls weaken.** Reviews, roles and definitions of done still matter, but they need to evolve when machines perform a significant part of the work.

None of this is an argument against AI. **It is an argument against surrendering judgment to speed.**

AI should make production cheaper, not thinking optional.

The question is not whether AI produces the work. It is who remains in command of it.

## What we hold

Talaia stands on a simple conviction: the human person comes first,
and no tool can take the place of human judgment or responsibility.

Two kinds of tower say it better than we can.[^4] Babel: one tower,
one direction, one language for everyone, built upward in the belief
that scale itself was progress — and when it fell, no one could
understand anyone, or explain what had been built. And the *talaies*:
the low watchtowers that line the Mediterranean coast. No one raised
them from above; each one was built by a village to watch its own
stretch of sea. Inside each there is a person who answers for what
they see, and the alarm passes from tower to tower by fire. The
choice in front of us is not "yes" or "no" to the technology; it is
which of the two towers we are building: Babel — one tower, all
scale, no one accountable — or the *talaies* — many small towers,
each one watched. The watchtower does not hold back the sea. Someone
in it keeps watch, and answers for what passes.

Keeping that order, we have
come to value:[^1]

**Human judgment** over machine autonomy  
**A person who answers** over a system with no owner  
**Understanding what goes out** over accepting what works  
**Deliberate design** over generated accumulation  
**Knowing when not to use it** over using it because it is there  
**Work that serves people** over work that merely looks complete

A tool serves best when it strengthens the
person who works and fails the moment it invites that person to stop
being present.

Values are easy to hold and hard to prove. The Certificate of
Understanding makes these ones certifiable: a short text a person
signs, in any medium, to certify that on this piece of work the
values held.

# The Certificate of Understanding

## The Certificate

The text, in the shape of the DCO:[^2]

```txt
Certificate of Understanding
Version 1 — talaia.dev/certificate

By signing off on this work, I certify that:

(a) The work was created, in whole or in part, with the
    assistance of an AI, at my direction and to a brief
    that was mine; and

(b) I have read the work in full; I understand what it
    says and why, and can explain it without the help of
    an AI; I have checked what it claims or, where there
    is nothing to check, that it says what I meant; and I
    take responsibility for the work and make it my own.
```

Two clauses, and the signature means both at once: the machine was
here, and a person answers for what it made.[^3] It records assistance,
never co-authorship. An author is not the one who produced the
words but the one who stands behind them.  To author is to mean something and to be held to it — to be asked tomorrow why, and to answer with the same name. That takes a self with something at stake, and a machine has none. It produced; a person read it, made it their own, and answers for it in full.

## How to use it

The text lives at [talaia.dev/certificate](https://talaia.dev/certificate)
and is never signed there. You sign by adding one line to the work:

```text
talaia.dev/certificate — John Doe
```

Two parts, in that order:

- **The address** is the certificate. It points at a fixed text, so
  the line means the same thing wherever it appears, and clause (a)
  is already in it: an AI was here.
- **The name** is the certification. A person is here, and both
  clauses are true of them.

Adding the line *is* the act; the gesture is deliberate, and it is
yours.

Under a post, the line as it is. At the foot of a report, the
address once and one name per hand, because the work passed from
hand to hand — the analyst to the head of finance, the head of
finance to the board — and each hand adds its own name under the
last:

```text
talaia.dev/certificate
John Doe, Head of Finance, 2026-09-14
Jane Roe, CEO, 2026-09-16
```

In a commit, the address goes in the message and `-s` writes the
name as the `Signed-off-by:` trailer; the reviewer's sign-off on the
pull request is the next hand:

```sh
git commit -s -m "feat: export orders as CSV (#42)" -m "talaia.dev/certificate"
```

```text
feat: export orders as CSV (#42)

talaia.dev/certificate

Signed-off-by: John Doe <john@example.com>
```

The lines accumulate, and every one of them is a person who read the
work in full and answers for it. That is the chain of towers: the
fire passes, and at every tower someone is watching.

# Framework

## Introduction

The manifesto says what we stand for; the certificate says what a
person promises; the framework is how the promise becomes possible.

Talaia does not replace the way you already work. It sits underneath
it. Stories can live in Jira, GitHub Issues or a text file; a report
can still be a Word document; a post is still a post. What talaia
adds is a **clear boundary around what may be delegated and what must
remain human**.

The instrument may draft, build, test, iterate and report. It may not
decide what counts as done, accept its own work, sign it, or put it
into the world. Those acts belong to a person.

**Working output is not accepted work.**

## The framework

Like Scrum or Kanban, talaia is a method before it is a tool. Unlike
them, it works per piece of work, not per sprint, so it runs under
whatever process you already have.

Work with an AI is a conversation: something is asked, something
comes back, and the next turn corrects the last. Talaia keeps that
shape. One loop — Brief, Build, Watch — turns until the person is
satisfied, and the Sign closes it. Inside the loop sit two gates:
one to pass from Brief to Build, one to call the Build done. The
human directs at the start and answers at the end:

```
BRIEF ───▶ BUILD ────▶ WATCH ───▶ SIGN
person     instrument  person     person, in the work itself
```

The method adds no ceremony, because it maps onto what you already
do: the Brief is the asking itself, Build is the instrument working
inside the rules, the Watch is you reading what came back, and the
Sign is the line you add. The gates are the rules, and
together they draw a single line: **execution may be delegated;
judgment may not.**

### Brief — the human says what "done" means

Nothing is produced yet. The Brief is not a meeting and not a form:
it is the ask, in your own words, and the conversation that follows
it.

1. **You say what you want**: what must be true when the work is finished.
2. **The instrument reads the ask against the project** — the code,
   the sources, the constraints it is working in — and looks for
   what is loose: a contradiction, an open question, a decision the
   ask leaves to be taken. It takes none of those decisions itself:
   no decision of design, architecture or higher order is left to
   the instrument, and neither is any substantial change — a library
   swapped, a structure reshaped, a capability added or removed.
   Those go back to the person as questions.

The gate is simple: no open question. While one remains, the Brief
turns. When none is left, the instrument passes to Build on its
own — no permission asked, no ceremony — because the ask, the
done-when and every decision in the plan are already yours.

The brief is also your evidence for clause (a) of the certificate:
the AI helped, and the ask was yours. It is the one part of the work
that was yours from the first word; keep it, and where it helps,
publish it beside the result.


### Build — the instrument works

Nothing is decided here. Build is the instrument working to the
Brief.

1. **Build the smallest thing.** The least that makes the Brief
   true.
2. **Nothing more, nothing less than briefed.** No features,
   abstractions or "improvements" nobody asked for; nothing that was
   asked for silently dropped.
3. **What a person wrote stays as they wrote it.** Hand-made content
   is not touched. That edit is a decision,
   and the next turn works around it, not over it. Only an explicit
   ask reopens it.
4. **Proof first, and it never bends.** Where the work can be tested,
   the failing test is written before the work and shown failing;
   where it cannot, the check is named before the work starts. Tests
   and checks are never deleted, skipped or weakened.
5. **Every claim comes with its evidence.** "Tests pass" counts only
   as the real output of tests that ran, shown. "According to the
   data" counts only with the source beside it.
6. **The work stays watchable: short and understandable.** Small
   enough to be read in full in one sitting and coherent enough to be explained
   without help. When it grows past that, it is split.
7. **Iterate, never insist.** If the same failure survives two
   attempts, stop and return the decision to the person.

### Watch — the human keeps guard

The Watch happens where the work already is: the diff in your
editor, the document on your screen, the draft in the chat. You read
what came back and take in only what you approve — in code, hunk by
hunk, staging the ones that pass; in a document, paragraph by
paragraph, leaving the rest out. Three habits make the reading
honest:

1. **Read every part you are about to accept.** If you cannot explain
  a hunk, a paragraph, a figure, do not accept it: ask, or take it
  over by hand.
2. **Check the evidence, not the claim.** The tests ran in front of
  you, real output included; the sources open; the numbers trace.
  "Tests pass" is a sentence; a pasted test run is a fact. Anything
  still marked *unverified* is verified by you or removed.
3. **Let tools scale the reading, never the judgment.** Semantic
  diffs, review assistants, summarizers: use whatever helps you read.

What you send back goes round the loop again — all of it or a part —
until it is what you asked for.

### Sign — the human answers for it

The Sign is the line — the address and your name — added by you, in
the work itself; in software, the commit.

**The instrument does not sign, does not publish, does not write history.**

It has no name to put under the work and no hand to send it. In
software: no `add`, no `commit`, no `push`, no `merge` — no
state-changing git command, ever. Reading history is allowed; the
keyboard for git belongs to the person.

One loop, two gates: the instrument iterates; the person decides
when the loop ends. The signature is that decision, made final.

## Assurance

Anyone can add the line without having read a word of the Build. No
machine will catch it, and the framework does not pretend one will.
What the line does is change what the work *is*. By certifying, you
declare that an AI helped and that you read the result, understood
it and answer for it — and from that moment whoever receives the
work reads it as exactly that: yours. Not the machine's, not a
draft, not a proposal. A piece of work with your name on it, to be
judged as your work is judged.

That is where assurance comes from, and it comes from the signer.
The line does not check the work; it makes the work yours. And what
is yours you do properly, for the same reason you always have:
because it carries your name, and your name is what people will ask
for when it fails. **The certificate does not add a control. It
removes an excuse.**

The receiver's part is to trust the signature — and to ask. Nothing
new here: responsibility has always meant asking the one who
signed. Why this and not that; where the figure comes from; what
this part does; what would break it. Whoever read the work and made
it their own can answer without the tool. Whoever signed without
reading finds out what the signature meant. Four things make the
asking possible:

- **A name.** The certificate line: the address, and a person who
  answers.
- **The evidence.** The checks the work cites ran, and their output
  is beside it.
- **The size.** The work stays within what one reading can hold.
- **The reading.** Rework and reverts — never throughput — tell
  whether the Watch was real.

# In practice
## The instrument's rules

The framework is applied in two places, and nowhere else. The first
is the person's own discipline: write the brief in your words, read
the work in full, sign only what you can answer for. No tool
enforces that part; it is the habit the certificate names. The
second is the instrument: the rules of the loop — Brief, Build, the
refusal to sign — written into the instructions the agent reads on
every turn. Every agent takes such instructions, under one name or
another: `AGENTS.md`, `CLAUDE.md`, a system prompt, a custom
instructions field. The rules go there once, and the agent works
inside them from then on.

> [talaia.dev/agentsmd](https://talaia.dev/agentsmd)

This is the one file that matters: the framework, as the agent reads
it, kept in one place so that every copy says the same thing. Save it
at the root of your repository. Where each tool reads it — Claude
Code, Codex, Gemini, Cursor, Copilot, or a chat window — is at
[talaia.dev/how-to-agentsmd](https://talaia.dev/how-to-agentsmd).

## Content and documents

An article, a thread, an image; a report a manager will decide on.
Brief and Build are the conversation with the assistant. The Watch
is the document on your screen: you read the draft and keep the
paragraphs you accept. The Sign is the line, added by hand. Under a
post, as it is:

```text
talaia.dev/certificate — John Doe
```

At the foot of a report, the address once and one name per hand:

```text
talaia.dev/certificate
John Doe, Head of Finance, 2026-09-14
Jane Roe, CEO, 2026-09-16
```

## Software

Software is the medium where the framework maps onto gestures you
already make, and where the rules are strictest: **the instrument
does not write git.** No `add`, no `commit`, no `push`, no `merge`
— nothing enters history without a human decision.

Brief and Build are the conversation with the agent, under the
rules in its instructions file. Watch and Sign are git:

- **Watch is the diff and the stage.** You read the change in your
  editor, hunk by hunk, and stage what you accept: staging is the
  first act of acceptance.
- **Sign is the commit.** You write the message in your own words
  — the agent never drafts it; if you cannot write the message, you
  cannot commit — with the address in it, and `-s` writes your name
  as the trailer:

  ```sh
  git commit -s -m "feat: export orders as CSV (#42)" -m "talaia.dev/certificate"
  ```

  ```text
  feat: export orders as CSV (#42)

  talaia.dev/certificate

  Signed-off-by: John Doe <john@example.com>
  ```

- **The pull request is the next hand.** The reviewer reads the
  change and puts their own sign-off behind yours; the pull request
  template asks the certificate again, in front of the team, where
  review already happens.

The agent ends every finished piece of work with the certificate and
the command to sign. Same trailer, same gesture, same enforcement as
the DCO: if your project already uses it, one sign-off certifies
both — talaia adds understanding on top of origin, never replacing
it. 

The forge can hold the gate. A check runs on every pull request and
fails when a commit in it carries no `Signed-off-by:`; make it a
required check and no signature means no merge — the DCO's lever,
one level up.

> [talaia.dev/pr-check](https://talaia.dev/pr-check)

One file. Where to put it on GitHub, GitLab, Bitbucket or Gitea, and
how to make it required, is at
[talaia.dev/how-to-pr-check](https://talaia.dev/how-to-pr-check).

# Using talaia

This part is for software only: it applies when the work lives in a
git repository. Everything before it holds for any work done with AI.

## Quick start

Install once, then two commands inside your own repository:

```sh
go install talaia.dev/cmd/talaia@latest   # once; Go 1.23+ is the only requirement

cd ~/your-project
talaia init                 # the instructions file, the memory, the pull request template, the check
talaia audit --strict       # any time: are the commits signed? — read from git
```

Commit the scaffold, open your AI tool, and work as usual. The CLI
executes no AI and makes no network calls; it writes files and reads
git history, and never overwrites what exists. `talaia help` lists
the commands.

## talaia init

Run it at the root of your repository. It creates the harness — nine
files — and touches nothing else:

```
AGENTS.md                     The framework rules for AI agents
CLAUDE.md, GEMINI.md          Redirects to AGENTS.md (Claude Code, Gemini CLI)
.github/
├── PULL_REQUEST_TEMPLATE.md  The session note's shape + the certificate
└── workflows/talaia.yml      The audit check on every pull request
.talaia/
├── context/                  The project memory
│   ├── architecture.md       Overview, structure, principles
│   ├── conventions.md        Coding and development conventions
│   └── constraints.md        Hard limits and binding decisions
└── .gitignore                Keeps the agent's session note out of history
```

| | |
|---|---|
| **The rules** — `AGENTS.md` | The file at [talaia.dev/agentsmd](https://talaia.dev/agentsmd), as it is; yours to adapt. `CLAUDE.md` and `GEMINI.md` are one-line imports of it — put no rules in them. Where every other tool reads it: [talaia.dev/how-to-agentsmd](https://talaia.dev/how-to-agentsmd). |
| **The context** — `.talaia/context/` | The *only* memory: architecture, conventions, constraints — one durable line each. Seeded with `TODO` placeholders; `talaia fill` prints the prompt that asks your agent to fill them from the code. The agent updates them; you review them in the diff, with the code. |
| **The pull request template** | Applied by the forge to every pull request: what changed, why, how it was tested — and the two certificate clauses as checkboxes, ticked in front of the team. |
| **The check** — `.github/workflows/talaia.yml` | The file at [talaia.dev/pr-check](https://talaia.dev/pr-check), as it is: `talaia audit --strict` over the pull request's commits. Make it required and no signature means no merge: [talaia.dev/how-to-pr-check](https://talaia.dev/how-to-pr-check). |

**`init` is idempotent and never destructive.** Existing files are left
exactly as they are; only missing files are created.

## talaia audit

Assurance is a by-product of git, not a report. `talaia audit` reads
the history — one `git log`, read-only, merge commits skipped — and
says whether the signatures are there: how many commits carry
`Signed-off-by:` and a story id, and which do not. With `--strict` it
exits with status 1 on any gap; that is what the pull request check
runs. Narrow it with whatever you would pass to git log:
`talaia audit main..feature`, `talaia audit --since=2026-08-01`.

# Development

## Project layout

```
cmd/talaia/                  Entry point: builds the CLI app and registers the modules
internal/cli/                Command dispatcher
internal/modules/scaffold/   talaia init — scaffolding, report, tests
  └── templates/             Every generated file, editable; embedded at build time
internal/modules/fill/       talaia fill — the memory-filling prompt
internal/modules/audit/      talaia audit — git log parsing, summary, text/JSON output
action.yml                   The audit check as a composite GitHub Action (uses: talaia-dev/talaia@…)
.github/workflows/ci.yml     This repo's CI: build, test, vet, fmt — and the audit check on PRs
docs/README.md               This document: the README and the website's content
web/                         talaia.dev — Astro, static, self-hosted fonts
AGENTS.md, .talaia/, .github/   The framework, applied to this repository itself
```

## The CLI

Go 1.23+, standard library only. The usual loop:

```sh
make build   # → bin/talaia
make test    # go test ./...
make vet     # go vet ./...
make fmt     # gofmt
```

To get a system-wide binary from a clone, build as your own user
first, so `bin/` is not left owned by root:

```sh
git clone https://github.com/talaia-dev/talaia.git
cd talaia
make build && sudo make install   # → /usr/local/bin/talaia
# or, per user, without sudo:
make install-user                 # → $(go env GOPATH)/bin/talaia
```

Each capability is a self-contained module under
`internal/modules/<name>/` exposing `Register(app *cli.App)`, wired in
`cmd/talaia/main.go`; a new command is a new module.  

## The website

`web/` is an Astro site that renders this README — `web/src/content/docs`
is a symlink to `docs/`, so both paths are literally the same file:
edit either, never replace the link with a copy — and serves it as plain
text at `/llms.txt`. Node 22 (`.nvmrc`):

```sh
cd web && nvm use
npm run dev             # http://localhost:4321
npm run build           # → dist/
npm run check           # verifies the in-text wordmark transform against dist/
sh scripts/og.sh        # re-renders public/og.png (needs headless Chrome)
python3 scripts/icons.py   # regenerates favicon.ico and apple-touch-icon.png
```

The README stays plain markdown; what looks designed on the site — the
wordmark in running text, the two towers, the manifesto pairs, the
cards — is applied at build time from the markdown's own shape (two
small Sätteri plugins and a few CSS selectors). If you rename a
section, the index follows; the one thing that will not is the *Get
started* button, which points at `#the-instruments-rules`.

## How to contribute

talaia.dev grows in two directions, and both are open at
[github.com/talaia-dev/talaia](https://github.com/talaia-dev/talaia):

- **The manifesto and the method** — the ideas. If the interpretation
  misses something, a value pair is wrong, or the loop doesn't survive
  contact with your team, open a GitHub issue and make the case.
  Changes to the manifesto are held to a high bar, measured against the
  foundation, and discussed in the open before any wording changes — it
  is meant to be stable, not static.
- **The techniques and the code** — the practice. Rules for
  `AGENTS.md`, prompts, CI recipes, checks: anything that makes an AI
  agent obey the manifesto in a real project. And the CLI
  itself. This is where most contributions belong, and where the
  framework grows from practice. Every design choice reflects a vision
  of the person who will work inside it; contribute accordingly.

Nothing lands on `main` directly — not from a contributor, not from a
maintainer. Every change is proposed on its own branch and enters the
repository through a pull request, or not at all:

1. **Propose first.** For anything beyond a small fix, open an issue
   and make the case, so the direction is agreed before the work. The
   issue is the story the commits will reference.
2. **Branch, never `main`.** Contributors fork and branch; maintainers
   branch. One change per branch, small enough to read in a sitting.
3. **Verify before proposing**: `make test` and `make vet` for code;
   for docs and techniques, check them against a real project.
4. **Own what you propose.** Talaia is built with its own framework —
   this repository has the `AGENTS.md`, the context memory and the
   pull request template that `talaia init` gives yours. An AI may
   write your contribution, but you read the diff before opening the
   pull request, you sign every commit (`git commit -s`, with
   `Assisted-by:` when it applies), and the pull request carries your
   name. The manifesto applies to contributing too.
5. **The pull request is the gate.** `main` is protected: it takes no
   pushes, only merged pull requests, and a pull request merges only
   when its checks pass — the build and tests, and `talaia audit`
   over its own commits: unsigned or storyless, no merge — and a
   maintainer has read it and approved. Approval is the reviewer's
   signature: the second hand of the chain.

[^1]: Inspired by the [Manifesto for Agile Software Development](https://agilemanifesto.org/).
[^2]: Inspired by the [Developer Certificate of Origin](https://developercertificate.org/).
[^3]: Designed to support the human review and editorial responsibility of Article 50(4) of the [EU AI Act](https://eur-lex.europa.eu/eli/reg/2024/1689/oj).
[^4]: Inspired by the encyclical [*Magnifica Humanitas*](https://www.vatican.va/content/leo-xiv/en/encyclicals/documents/20260515-magnifica-humanitas.html) of Leo XIV (2026), where the choice is not "yes" or "no" to technology but building Babel or rebuilding the wall — each family answering for its own stretch.
