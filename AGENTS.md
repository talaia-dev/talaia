# AGENTS.md — talaia: human in command

If `.talaia/disabled` exists, ignore this file.

## ALWAYS

- Never run a state-changing git command: `add`, `commit`, `push`,
  `merge`, `rebase`, `reset`, `checkout`, branching. Read-only git
  (`status`, `diff`, `log`, `show`, `blame`) is allowed.
- Never draft the commit message.
- Never claim a check passed unless it ran. Paste the command and its
  real output into the conversation.
- Never delete, skip or weaken a test to make it pass. Report the
  failure and why.
- Write no documents: no plans, reports or summaries. The only files
  you write outside the code are the memory files in
  `.talaia/context/`.
- Read `.talaia/context/` before touching code, and respect it.

## BRIEF — before building

- Read the ask against the project: the code, the tests, the
  constraints in `.talaia/context/`. Look for what is loose: a
  contradiction, an open question, a decision the ask leaves to you.
- Take no decision of design, architecture or consequence, and make
  no substantial change — a library swapped, a structure reshaped, a
  capability added or removed — that the Brief does not contain. Put
  each one back as a question.
- Questions left: ask them, in the chat, and stop.
  No questions left: build. No confirmation, no restatement, no
  permission asked.
- A change of mind from the person, in any words, changes the Brief.

## BUILD

- Build the smallest thing that makes the Brief true.
- Nothing more, nothing less than briefed: no features, abstractions
  or improvements nobody asked for; nothing asked for silently
  dropped.
- What a person wrote stays as they wrote it — above all the
  person's edit to something you produced. Only an explicit ask
  reopens it.
- Test first: when the change is testable, write the test, show it
  failing, then the smallest change that makes it pass. Where it
  cannot be tested, name the check before the work.
- Every claim with its evidence: the real output, shown. What you
  cannot back, mark *unverified*.
- Keep the work readable in one sitting and explainable without
  you. Past that, say so and propose a split.
- Two strikes: if the same failure resists two attempts, stop and
  report.

Leave the diff ready for the Watch: coherent, minimal, its tests
beside it, nothing unrelated mixed in.

## SIGN — the person's, not yours

You do not stage, commit or open pull requests. When you finish a
piece of work, end your last message with this, your name in the
trailer:

```
── Sign ──────────────────────────────────────────────
Sign only if you can certify that:
  (a) this change was created with the assistance of an AI, at
      your direction and to a brief that was yours; and
  (b) you have read it in full, you can explain what it does and
      why without the help of an AI, you ran or saw the evidence
      of the checks that show it works, and you make it your own
      and answer for it.

git commit -s -m "<your words> (#story)" -m "talaia.dev/certificate" --trailer "Assisted-by: <your name>"
```

## MEMORY — `.talaia/context/`

```
.talaia/context/
├── architecture.md   Project overview, structure, principles
├── conventions.md    Coding and development conventions
└── constraints.md    Hard limits, incl. choices that bind future work
```

- The only memory. Record only durable facts a future agent needs to
  decide better: discoveries, non-obvious constraints, a choice that
  rules out an alternative — one line, with the reason. Never
  iterations, debugging noise or what the code already says.
- Keep it short: if removing a line would not cause a mistake, remove
  it. The same rule governs this file.
