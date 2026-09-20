# Task: populate this project's context memory (.talaia/context/)

Analyze this repository **file by file** before writing anything:
directory layout, build and dependency files, source code, tests,
docs, and read-only git history if available. Then fill in the context
memory. Follow the rules in AGENTS.md.

## What to fill

1. `.talaia/context/architecture.md` — replace every TODO: project
   overview (what it is, what it does, for whom), system structure
   (the real components/directories and how they relate), and the
   architectural principles actually visible in the code.
2. `.talaia/context/conventions.md` — the real build/test/lint
   commands, the code style the codebase actually follows, and its
   testing conventions.
3. `.talaia/context/constraints.md` — constraints evidenced by the
   project: dependency policy, compatibility targets, performance,
   security or compliance requirements, and any decision visible in
   the code or history that rules out an alternative for future work
   (one line each, with the reason). If a section has no evidence,
   state that briefly.

## Rules

- Ground every statement in files you actually read; cite paths.
- Never guess: if something cannot be determined from the repository,
  say so — do not invent facts, commands or requirements.
- Keep every file concise and high-signal: this is working memory for
  future agents, not documentation for its own sake. For every line
  ask whether removing it would cause a future mistake; if not, leave
  it out.
- Do not create any other file: no feature files, decision logs or
  plans. The story lives on the team's board, the change in git, the
  explanation in the pull request.
- Do not modify anything outside `.talaia/context/`.
- Do not run any state-changing git command.

## When done

Report which files you filled and what you could not determine from
the code.
