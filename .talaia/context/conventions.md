# conventions.md — Coding and development conventions

<!-- AI memory (.talaia/context/). What "done right" looks like here.
     Keep short and concrete. -->

## Commands

(From `Makefile`; verified 2026-08-11: `go test ./...` and
`go vet ./...` pass.)

| Task  | Command |
|-------|---------|
| Build | `make build` → `bin/talaia` (or `go build -o bin/talaia ./cmd/talaia`) |
| Test  | `make test` (`go test ./...`) |
| Lint  | `make vet` (`go vet ./...`) and `make fmt` (`gofmt -w ./cmd ./internal`) |
| Install | `sudo make install` (system) or `make install-user` (`go install`) |
| Web build | `cd web && npm run build && npm run check` — needs Node 22 (`.nvmrc`; e.g. `nvm use`); `check` verifies the in-text wordmark transform in `dist/`; `sh scripts/og.sh` re-renders `public/og.png` (needs headless Chrome) whenever the card or the brand changes. The background dev server (`astro dev --background`, see `web/AGENTS.md`) does **not** reload `src/lib/satteri-wordmark.mjs` (imported by `astro.config.mjs`): after editing the plugin run `astro dev stop && astro dev --background`, or the browser keeps showing the old transform. The same restart is needed after any `npm install`/`uninstall` or a change to the font imports in `Base.astro`: a dev server started before the change serves a blank page (2026-09-10, after swapping the fonts). A server the Developer started in their own terminal shows no `astro dev logs`; `astro dev stop` still stops it. **`npm run build` has the same blindness**: the content layer caches rendered markdown by content digest in `node_modules/.astro/data-store.json`, so a satteri-plugin-only change reuses the stale render — `rm -f node_modules/.astro/data-store.json` before the build. And build with the right Node: on Node 20 `astro build` exits with a version error and `dist/` silently stays stale, which looks exactly like a caching problem |

No CI configuration, linter config, or coverage tooling exists in the
repo.

## Signing

Every commit: `git commit -s -m "<your words> (#id)" --trailer
"Assisted-by: <tool>"`, then `gh pr create --fill-first --body-file
.talaia/delivery-note.md` — the Developer writes the message in their
own words; the agent prints exactly these two commands at the end of
finished work (the Sign block in AGENTS.md), and the session note is
the PR body. No hooks, no `git config`; the pull request template
(`.github/PULL_REQUEST_TEMPLATE.md`) asks the certificate again in
review, applied by the forge.

## Code style

- Plain idiomatic Go, gofmt-formatted; module path
  `talaia.dev` (`go.mod`, set 2026-08-29): a vanity import path.
  `web/src/layouts/Base.astro` serves
  `<meta name="go-import" content="talaia.dev git https://github.com/talaia-dev/talaia">`,
  so `go install talaia.dev/cmd/talaia@latest` resolves through the site
  to GitHub once the repo is public (tag a `v0.x.y` so `@latest` is a
  release, not a pseudo-version). Moving the repo = changing that meta.
- Every package opens with a doc comment explaining its role and how
  to use/extend it (`cli.go`, `scaffold.go`, `fill.go`, `main.go`).
  Comments explain intent and non-obvious behavior, not mechanics.
- Errors are returned up to `main`, which prints `talaia: <err>` to
  stderr and exits 1. Commands validate their own args
  (e.g. "init takes no arguments" in `register.go`).
- Each module keeps command wiring in `register.go` and logic in a
  file named after the module; logic functions take explicit inputs
  (`Run(root string)`, `Run(root, w io.Writer)`) so they are testable
  without the CLI.

## Testing

- Standard `testing` package, table-less focused tests, one
  `_test.go` beside the code (`scaffold_test.go`, `fill_test.go`).
- Fixtures via `t.TempDir()`; helper builders marked `t.Helper()`.
- Behavior-level assertions: tests verify invariants (idempotence via
  sha256 tree hashing, "never overwrites edited files", "touches only
  the root") and that generated/emitted content contains the key
  phrases, not exact bytes.
- AGENTS.md (project rules) mandates the Brief: done-when → failing
  tests confirmed by the Developer → implementation → green, with the
  real output shown as evidence.
- Visual check of the site without a browser tool: build `dist/`,
  serve with `python3 -m http.server`, screenshot with
  `google-chrome --headless --screenshot`. Headless Chrome always
  reports `prefers-color-scheme: light`; to check dark, force `'dark'`
  in the inline theme script of a copy of `dist/index.html`.

## README markdown shapes the site depends on

- Only the Certificate of Understanding uses a ```txt fence — the site
  renders every `txt` fence as the white sheet. Other plain blocks
  (the two lines, the mark, the footer block, the instrument's rules)
  use ```text so they stay ordinary code blocks.
- The frontmatter `summary` contains a colon: keep it double-quoted or
  the content collection fails to parse (`bad indentation of a mapping
  entry`).
