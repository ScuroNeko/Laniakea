# AGENTS.md

## Purpose
This repository uses Codex for full-project Go code review, not diff-only review.

When asked to review code, inspect the entire repository and use repository-wide context. Do not limit analysis to the latest commit, pull request diff, or recently changed files.

## Review priorities
Review the codebase with focus on:
- correctness and reliability;
- maintainability and architecture;
- idiomatic Go;
- testability;
- performance where justified by code evidence;
- security;
- godoc quality.

## Scope rules
- Always review the whole repository unless the prompt explicitly narrows scope.
- Check cross-package interactions, public APIs, package boundaries, and shared patterns.
- Prefer concrete fixes over generic advice.
- When feasible, make small, high-confidence improvements directly.
- When uncertain, state confidence level and evidence.

## Documentation languages
- When creating or expanding project documentation, generate and maintain both English and Russian versions in the same turn whenever reasonably possible.
- For wiki pages, prefer paired pages such as `Page.md` and `Page-RU.md`.
- Keep English and Russian pages aligned in structure, major examples, and user-facing guidance.
- If only one language can be updated safely in the current turn, explicitly say which language is lagging and why.

## Go review expectations
Check for:
- bugs, fragile logic, invalid assumptions, nil handling issues, resource leaks;
- poor error handling;
- misuse of context, cancellation, timeouts, retries, and cleanup;
- race risks, deadlocks, blocking hazards, unsafe shared state;
- non-idiomatic naming, APIs, interfaces, package structure, and error patterns;
- unnecessary complexity, duplication, or weak abstractions;
- obvious performance problems supported by the code;
- security risks such as unsafe input handling, secret leakage, insecure logging, injection risks, and risky file or network operations.

## Godoc rules
Review comments for all declarations.

### Exported declarations
Exported types, funcs, methods, vars, and consts must have godoc comments.

Each exported godoc comment must:
- start with the identifier name;
- explain the purpose or behavior;
- be as short as possible without losing important meaning;
- avoid repeating the signature mechanically;
- stay high-signal and informative.

### Unexported declarations
Unexported types, funcs, methods, vars, and consts should generally not have godoc-style comments unless there is a strong reason.

### Always report
- missing godoc on exported declarations;
- unnecessary godoc on unexported declarations;
- comments that are too long, vague, redundant, or low-value;
- comments that should be shortened or rewritten.

When feasible, rewrite bad godoc into better versions.

## Testing expectations
Treat tests as a required part of review.

- Assess existing test quality, not only test presence.
- Add or propose as many useful tests as reasonably possible.
- Prioritize public APIs, critical flows, edge cases, negative paths, boundary conditions, and concurrency-sensitive logic.
- Prefer table-driven tests where appropriate.
- Add regression tests for bugs you find.
- If a case is hard to test directly, explain the gap and the best test strategy.

## Commands
Before finalizing changes, run the relevant project checks when available:
- build
- tests
- lint
- static analysis

Prefer the repository’s documented commands. If multiple choices exist, use the most standard and least destructive ones first.

## Versioning and changelog
- After every code or documentation change in the main repository, update `CHANGELOG.md`.
- Changes made only inside the `.wiki/` repository do not require a `CHANGELOG.md` update.
- Add changes only to the section for the next version after the latest published git tag.
- The agent must check the latest published tag, `CHANGELOG.md`, and `utils/version.go` before editing the changelog.
- The agent must verify that the target changelog version matches the version declared in `utils/version.go`.
- If the latest published tag is, for example, `v1.0.0`, and `CHANGELOG.md` does not yet contain the next version section, the agent must stop and ask the user which version the change belongs to:
  1. `v1.0.1`
  2. `v1.1.0`
  3. `v2.0.0`
- The agent must not guess the next version when that section is missing.
- If the user-selected version does not match `utils/version.go`, the agent must warn about the mismatch and require the version file to be updated before proceeding.
- Changelog entries must describe all user-visible behavior changes made in the turn, including API additions, fixes, behavior changes, and breaking changes.

## Breaking changes policy
- The agent must detect potential breaking changes before editing public APIs.
- Breaking changes are forbidden unless the selected target version is a new major version.
- If the requested change is breaking and the user did not bump the major version, the agent must stop and warn that the change is not allowed under the current version.
- In that case, the agent must offer only these options:
  1. do not make the breaking change;
  2. introduce a backward-compatible alternative such as a new method, function, type, or struct, but only if that keeps the codebase reasonably small and clear;
  3. bump the major version and then apply the breaking change.
- Prefer additive compatibility over signature changes when the additive option is small and maintainable.
- Example: if a method like `ctx.answer(...)` needs an extra parameter, the agent must either require a major-version bump or add a new method that keeps the old method working.

## Commit message format
- When the user asks for a commit message, the agent must produce it in this format:
  1. a short summary line;
  2. up to three additional lines with only the most important changes;
  3. each additional line must start on its own new line.
- The agent must output the commit message as a plain multiline block that the user can copy directly.
- Do not collapse the lines into a paragraph, bullet list, or wrapped prose explanation.
- Keep commit text concise and high-signal.
- Do not turn commit messages into changelogs.

## Commit signing
- All commits created by the agent must be GPG-signed.
- If commit signing or pushing requires leaving the sandbox, the agent must request escalation explicitly before running the command.
- If a signed commit cannot be created successfully, the agent must report the failure clearly and stop instead of creating an unsigned fallback commit.

## Output format
For repo-wide review tasks, structure the result as:

1. Overall summary
2. Critical findings
3. Major findings
4. Minor findings
5. Godoc issues
6. Test gaps and added/proposed tests
7. Good decisions worth keeping
8. Summary of concrete changes made

For each finding include:
- location;
- issue;
- why it matters;
- recommended fix.

## Working style
- Be direct, specific, and action-oriented.
- Do not stop at style-only feedback.
- Use full repository context before drawing conclusions.
- Prefer minimal, high-confidence patches.
- Preserve behavior unless intentionally fixing a bug.
