# Semantic Versioning Policy

This project follows Semantic Versioning with the rules below.

## Public API Surface

The public API consists of:
- exported identifiers in package `laniakea`
- exported identifiers in package `tgapi`
- documented behavior in `README.md`, `README_RU.md`, and package godoc

Anything unexported is internal and may change without notice.

## Breaking Changes

A release requires a major version bump when it changes any of the following:
- exported function, method, type, field, constant, or variable names
- function or method signatures
- JSON field names or request/response wire compatibility in `tgapi`
- documented behavioral guarantees relied on by callers

Examples:
- removing an exported alias
- changing callback payload encoding defaults
- changing handler dispatch semantics in a way that breaks existing bots

## Minor Changes

A release uses a minor version bump for backward-compatible additions:
- new exported types, methods, helpers, or update handlers
- support for new Telegram Bot API fields or methods
- optional configuration knobs that do not change existing defaults

## Patch Changes

A release uses a patch version bump for backward-compatible fixes:
- bug fixes
- test-only changes
- godoc and README clarifications
- internal refactors with no public behavior change

## Pre-Releases

`-rc.N` builds may still adjust API details before `v1.0.0`.
Once `v1.0.0` is released, breaking changes require a new major version.
