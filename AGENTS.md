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
