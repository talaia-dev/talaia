# architecture.md — Project overview, structure and principles

<!-- AI memory (.talaia/context/). Practical guidance, not exhaustive
     docs. Update when the structure or the principles change. -->

## Project overview

**The Developer's vision (2026-09-11, stated after three brand passes):**
talaia's centre is the **Certificate of Understanding**, not the
framework. Everything reads as a chain — *Manifesto* (why the person
stands above the tool) → *Certificate* (since 2026-09-16 two clauses a human
certifies with `Signed-off-by:`: made with the help of an AI; read in
full, understood, answered for, made their own) → *Framework* (Brief/Build/Watch/Sign,
which includes the certificate and makes its clauses possible) → *CLI
and `git commit -s`* (which install and perform it). README and site
should present that order and that emphasis; the framework and the
tool are in service of the certificate. The site's look is that of a
published standard (references: NIST SSDF, SLSA, Green Software
Foundation), never editorial, never playful/terminal. Applied to the README on 2026-09-14 and widened the same day: the
certificate serves any work done with AI — published content (a one-line
mark under the post, the brief as optional evidence), documents and
reports (a footer block, one `Signed-off-by:` per hand up the chain, the
receiver may ask the three questions) and software (the only medium
with a machine check). Parts: Manifesto / The Certificate of
Understanding / Framework / In practice / Using talaia / Development.
The hero, the `<meta description>` in `index.astro` and the lede in
`llms.txt.ts` still say "software development" — not yet updated.

Talaia is a human-in-command accountability framework for AI-assisted
development — a **standard plus its check**, like the DCO. Method:
Brief → Build → Watch → Sign, mapped onto what developers already do
(Brief = the prompting loop; Build = the agent under AGENTS.md rules;
Watch = reading the diff in the IDE, staging = first acceptance; Sign
= `git commit -s` with the Certificate of Understanding and the
`Assisted-by:` trailer; the PR = the reviewer's sign). The Go CLI is
deliberately tiny — `init` (installs the harness: AGENTS.md +
`CLAUDE.md`/`GEMINI.md` redirects, PR template, `.talaia/context/`),
`fill`, `audit`, `help`. The story lives on the team's board, the
change in git, the explanation in the pull request — the framework
generates no document a human must read (`docs/README.md`,
`cmd/talaia/main.go`).

Talaia executes no AI itself. `talaia init` scaffolds files;
`talaia fill` prints a prompt for the user to pipe into their agent
(`talaia fill | claude -p`). This repository is itself managed with
its own framework (root `AGENTS.md`, `.talaia/context/`). The
`.talaia/decisions/` and `.talaia/features/` folders in this repo are
history from before 2026-08-28 and are no longer maintained; anything
still binding from them lives in `constraints.md`.

The public site (`web/`, Astro, static) renders `docs/README.md` —
`web/src/content/docs` is a **directory symlink** to `../../../docs`
(verified 2026-09-01; earlier noted here as a file hardlink), so both
paths are literally the same file: edit either in place; replacing the
symlink with a copy would silently fork the content.

## System structure

| Area | Purpose |
|------|---------|
| `cmd/talaia/` | Entry point: creates the `cli.App` and registers each module (`scaffold.Register`, `fill.Register`, `audit.Register`). |
| `internal/cli/` | Minimal command dispatcher: `Command` structs registered on an `App`, name-based routing, built-in `help`. |
| `internal/modules/scaffold/` | The `init` command: creates `AGENTS.md`, redirects, and `.talaia/` — never overwrites; returns a `Report` (Created/Kept). |
| `internal/modules/scaffold/templates/` | Editable source of every generated file, embedded with `go:embed`. Edit here, rebuild. Since 2026-08-29 it also holds `pull-request-template.md` (→ `.github/PULL_REQUEST_TEMPLATE.md`: the delivery note's six headings, an "Assisted by:" line and the four certificate checkboxes) Since 2026-09-02, `talaia-workflow.yml` (→ `.github/workflows/talaia.yml`: init installs the audit check itself — branch protection stays a one-time manual/`gh api` act, since settings are not files) and `talaia-gitignore` (→ `.talaia/.gitignore`: hides `delivery-note.md`, the agent's **session note** — removed with the worktable on 09-01, restored 09-02 as flight-recorder + PR body; shape: What changed / Why / How it was tested / Three questions). Template names avoid leading dots because `go:embed` skips dotfiles. Git hooks and a `.gitmessage` were built and then removed the same day: they depended on per-clone `git config` that may never run; the certificate is instead reminded by the agent (AGENTS.md's Sign block) and asked by the PR template, both zero-config. |
| `internal/modules/fill/` | The `fill` command: emits `templates/fill.prompt.md` to stdout; refuses to run without an existing `.talaia/`. |
| *(removed 2026-09-01)* | An entire MWI system — `internal/mwi` (frontmatter model, status machine), `internal/term` (TTY styling), `internal/modules/work` (gates, then no-gates sign/close), `internal/modules/board` (terminal observation portal), a worktable in `.talaia/work/` and a generated delivery note — was built and **deleted the same day**, never committed. Verdict recorded in constraints.md: generated artifacts nobody reads; talaia is a standard + check, not tooling. Do not rebuild it. |
| `internal/modules/audit/` | The `audit` command (2026-08-29): one `git log --no-merges --shortstat --format=<RS/US-separated>` parsed by `Parse`, summarized by `Summarize` (signed / story `#id` / Assisted-by / Reviewed-by among assisted / within budget / reverts, gaps list, per-author), printed by `WriteText`/`WriteJSON`. Pure core testable without git; `TestRunOnARealRepository` builds a temp repo and skips if git is absent. `--strict` fails only on unsigned or story-less commits; the budget (default 300) is advisory. |
| `action.yml` (repo root) | The enforcement check as a composite GitHub Action (`uses: talaia-dev/talaia@<ref>`): setup-go with `go-version-file` from the action path, builds the binary there, computes the range from the PR event (`base.sha..head.sha`, overridable via `inputs.range`), runs `audit --budget --strict`, tees the report into `$GITHUB_STEP_SUMMARY` and emits `::error::` on failure. Consumers need `actions/checkout` with `fetch-depth: 0`. |
| `.github/workflows/ci.yml` | The repo's own CI: Go build/test/vet/gofmt on push+PR, and the audit check (`uses: ./`) on PRs — talaia enforced on talaia. |
| `Makefile` | build / install / install-user / test / fmt / vet / clean. |
| `docs/README.md` | The single source of the README and the website content (hardlinked into `web/src/content/docs/`). |
| `web/` | talaia.dev site: `src/pages/index.astro` builds the sidebar from the rendered headings: each README `#` opens a group (label and anchor = the h1), each `##` under it is an item — since 2026-08-29 nothing is hardcoded (the old `PART_STARTS` slug map broke whenever a section was renamed). The top bar and the footer link to the `#` ids (one per README `#`; `/certificate` is its own page, `src/pages/certificate.astro`: the certificate text only, extracted from the README's ```txt fence at build time); the hero's *Get started* points at `#the-instruments-rules` (In practice › The instrument's rules, where talaia.dev/agentsmd is; since 2026-09-20, so a non-developer does not land on `go install`) — rename that section and update the button. Since 2026-09-10 the README opens with a **summary before its first `#`** (lede + the core as a list + one closing line): the site styles it by position — `.readme > p:first-child` / `ul:not(h1 ~ *)` select everything not preceded by an h1, the first `p` gets an `Overview` overline — so any heading placed there would break both the styling and the index (a `##` before the first `#` is dropped from the toc). Manifesto and Framework each open with a `## Introduction` that exists for the index only: `.readme h2[id^="introduction"]` hides it visually while keeping it in the flow as a scroll target; the second one gets the `-1` slug; `src/pages/llms.txt.ts` emits the README as plain text; brand: since 2026-09-03 the mark is **the Developer's own definitive design**, adopted verbatim from their sheet at `web/brand/logo_talaia_negativo_variantes.svg` (+ .png) — a machicolated pixel watchtower: fire = one free 12×12 amber pixel at (24,0); three 12×12 merlons at y18; 60×12 parapet; 48×36 body inset 6 (the machicolation); 8×13 door at (14,49), rendered as an even-odd **hole** so any background shows through; native units 60×78. Brand palette from the sheet: navy `#1A2238`, cream `#FAF6EF`, amber `#E8A33D` — on the site (2026-09-11) the amber is **only** the fire and the focus ring; the **navy is the institutional colour**: the hero is a full-width navy band (`slot="hero"` in `Base.astro`, rendered above `.shell`; cream tower, white name, an uppercase status line "Open standard · Working draft", the tagline as the `h1`), links and the current index item are navy (`--accent`), pull-quotes carry a navy bar. Light is the **default** regardless of system preference (`localStorage 'theme' ?? 'light'`); dark is the same page in the navy (`#161c2e`). No `body::before` glow or grain, radii 4px. Page colours: white paper `#fff`/`#f4f5f7` with navy ink; a cream page + amber everywhere read as "educational" (09-10), a near-black dev-tool dark mode as "childish" (09-11). `Logo.astro` keeps the native 60×78 viewBox with no padding so the tower's foot is the box's bottom edge: lockups (`.brand`, hero `.lockup`) are `inline-flex; align-items: baseline` and the tower **stands on the word's baseline**, `.98em` tall (merlons = ascender height, fire alone above), gap `.42em`, the name in Inter 600 — **no cursor** (the sheet's amber bar was dropped 2026-09-11: with the pixel tower it read as a game), no animation, no ground line, and the mark stays small (30px in the hero, 18px in the bar). Set 2026-09-10/11 after the image lead's reviews: a 56px tower vertically centred beside a bold soft-serif word read as "two separate things" and as a children's/educational brand; a big mono/terminal lockup read as "still childish" — do not go back to centring, to a tower larger than the type, or to a large mark. The in-text `.wm` transform is untouched (now plain bold). Favicon = the sheet's favicon variant: cream 62-tile rx12, navy tower, amber fire (`favicon.svg`, mirrored in `icons.py`: translate(15,9) scale(0.52), same at every size). A week of agent-drawn candidates (engraved tower, chevron towers, pixel castles, layer stack "cake", `t` monogram that read as a cross — never retry that one, the README is deliberately secular, commit-chain "thermometer", prompt marks) all failed or were superseded; the Developer drew the final one. `public/favicon.svg` duplicates the path with a prefers-color-scheme swap; `favicon.ico` (16/32/48) and `apple-touch-icon.png` (180, navy badge, white mark) are generated by `web/scripts/icons.py` — pure Python, no deps: mirror any path change there and rerun. The sprite reads at every icon size (2×2-cell door survives 16px), so `icons.py` no longer has a small-size variant. Layout (2026-08-28, Vercel/PostHog-leaning iteration): sticky blurred top bar (`.topbar` in `Base.astro`: brand, one link per `#` part, GitHub, theme toggle; links hidden ≤800px), hero = mark (56px) + wordmark, tagline ("You direct. The AI codes. The commit is your signature."), subline and two `.btn` CTAs. Markdown-driven styling via selectors, so the README stays plain: the manifesto pairs paragraph (`p:has(> strong:first-child):has(> br)`) renders as a two-column grid (each `strong` and each text run is a grid cell, `br` hidden); the four items of "The project" (`#the-project + p + ul`, items written `**Title.** text`) render as 2×2 cards; tables use horizontal hairlines and uppercase small headers; `pre` has 10px radius; a `blockquote` renders as a serif pull-quote with the beacon bar (used for the framework's one-sentence summary in the Framework introduction), (the Certificate of Understanding is no longer a blockquote: it is the ```txt fence, rendered as the white sheet). Every plain-text "talaia"/"Talaia" in the README body becomes the in-text wordmark at build time (lowercase, serif, a plain i — the fire-dotted i sat too high in running text and was dropped there on 2026-08-28; it stays in the hero and top bar) via `web/src/lib/satteri-wordmark.mjs`, a Sätteri hast plugin registered as `markdown.processor: satteri({ hastPlugins })` in `astro.config.mjs` — Astro 7 uses Sätteri, not unified: `rehypePlugins` would need `@astrojs/markdown-remark`, deliberately not added. A second Sätteri plugin, `src/lib/satteri-towers.mjs`, finds the manifesto paragraph that opens with "Two kinds of tower" and inserts before it the one illustration of the site — a `figure.towers` diptych: Babel (stepped, no fire, 38% opacity) beside the talaia with its fire and the next tower on the sea line, mono captions — and gives the paragraph `p.towers-text`. Since 2026-09-03 figure+paragraph render as **one framed block** (border+bg-alt, `--measure` width), drawn by CSS in two halves across the adjacent siblings (`.towers` = top half, `p.towers-text` = bottom); the serif lede voice is gone. Sätteri's `wrapNode` was tried for a real wrapper element and **silently drops the wrapped node** — don't use it; the two-halves CSS is the working pattern. Rename that opening sentence and the figure silently disappears. Code, headings, domains (`talaia.dev`, `talaia-dev`) and path segments stay literal; `npm run check` (`web/scripts/check-wordmark.mjs`) verifies that against `dist/`. Identity pass 2026-08-29 (serene / open-source / talaies), re-set 2026-09-10 and again 2026-09-11 after the image lead's reviews ("editorial / school", then "still childish"). The references the Developer gave for the target are **published standards**: NIST SSDF, SLSA, the Green Software Foundation — light page, one deep institutional colour as a field, a humanist sans at ordinary weights and tracking, the vocabulary of a spec (status, versions, overview, numbered parts), no gadgets. So: **no serif anywhere** (Fraunces soft, then Source Serif 4 sober, both read as a publisher) and **no mono in the UI** either (a JetBrains-Mono name/nav/labels/cursor "terminal voice" read as childish): Inter only — `@fontsource-variable/inter/opsz.css`, `"opsz" 32` and weight 600 for the tagline and titles, tracking never tighter than −.02em; uppercase 11.5px 600 sans overlines for labels; parts numbered by `counter(part)` (six since 2026-09-14); JetBrains Mono 400 only for code. Do not bring back a serif for "manifesto" feel nor the mono for "developer" feel — the voice is the standard's. Part titles once kept a short amber bar (a roman "Part I…IV" overline with a fire-and-horizon signal was tried on 2026-08-29 and rejected as too book-like; the mono `01` number is the code-like replacement, 2026-09-10); hero ends with `.hero-seal` (mono line: open source · any agent · nothing to configure); footer = mark + wordmark, one line, four links; smooth scroll, amber `:focus-visible`, thin scrollbars, print stylesheet; motion slowed (logo fires 12s, wordmark fire 9s, gentler keyframes). Social card: `src/pages/og-card.astro` (noindex, excluded from the sitemap via the sitemap `filter`) rendered to `public/og.png` by `scripts/og.sh` with headless Chrome after a build (captured at 1200×900 and cropped by `scripts/pngcrop.py`, pure Python: this Chrome's `--window-size` includes ~87px of window chrome, so a 630px window yields a 543px viewport — the same clamp that makes mobile captures 500px wide); `Base.astro` carries og:/twitter meta pointing at https://talaia.dev/og.png. Rejected earlier the same day: bottom horizon strip, hero seascape — keep the page calm. The index is a scroll-spy (inline script in `Base.astro`): the `##` whose heading last passed under the top bar (96px line) gets `aria-current="true"` on its link, styled in `.toc li a[aria-current]`; the hidden Introduction headings count because `position:absolute` keeps their static place. On phones the index (`.toc`) goes after the content (`order: 2`) and the grid column is `minmax(0,1fr)` — plain `1fr` let wide code blocks overflow the viewport. In the mark the single fire breathes — opacity 1→.5, 12s loop (`.beacon-main` in `global.css`, off under `prefers-reduced-motion`; the two-fires-take-turns animation left with the distant tower). The `body::before` light layer (amber glow / paper grain) was removed on 2026-09-10 with the rest of the ambient amber. A hero "seascape" (big tower + sea line + far towers) was tried and rejected on 2026-08-28, as was a fixed bottom strip of tiny towers: keep the hero to logo + name. |

## Architectural principles

- One self-contained module per capability under
  `internal/modules/<name>/`, exposing `Register(app *cli.App)`;
  wired in `cmd/talaia/main.go`. New capabilities follow the same
  pattern.
- Generated content lives as editable template files embedded at
  build time (`go:embed` in `scaffold.go` and `fill.go`), never as Go
  string literals. A missing embedded template panics — treated as a
  build defect, not a runtime condition.
- Non-destructive by design: `scaffold.Run` skips any existing file
  and reports it as Kept; re-running `init` is idempotent
  (`scaffold_test.go`).
- Standard library only; no external dependencies (`go.mod` has no
  requires). Fully offline.
- The tool orchestrates no AI: it writes files and emits prompts; the
  human pipes prompts into whatever agent they use.
