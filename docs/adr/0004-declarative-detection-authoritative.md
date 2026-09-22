# ADR-0004: Declarative files are the authoritative source

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** detector-engineer, scan-engineer

## Context

Repositories encode their intended environment in machine-readable files such
as `.tool-versions`, `.nvmrc`, `mise.toml`, `package.json` `engines`,
`rust-toolchain.toml`, `go.mod`, `.python-version`, `devcontainer.json`,
`Dockerfile`, and `Makefile`. Inferring requirements from general prose is
error-prone and produces false positives that erode trust.

## Decision

Detection is deterministic and declarative: each detector reads exactly one
file kind and emits `tool`, `constraint`, `source file`, `component`,
`category`, and `confidence`, never inventing a version absent from the file
(SCAN-FR-07). v1 ships exactly the eleven listed detectors (SCAN-FR-08).
README prerequisite detection is deliberately conservative — it recognizes only
explicit prerequisite declarations, marks them with reduced confidence and a
README source, and never derives requirements from general prose (SCAN-FR-09).
v1 has no user configuration and no per-repo manifest; the architecture leaves
room to layer them on later (CONST-05).

## Alternatives Considered

- **LLM or heuristic inference over README prose.** Rejected: nondeterministic
  and prone to false requirements.
- **Require a `repo-ready.yaml` manifest.** Rejected for v1: it adds an authoring
  burden and a config surface the product explicitly avoids.
- **Merge all declaration sources and pick one winner silently.** Rejected:
  conflicts must be surfaced, not resolved (see ADR-0005).

## Consequences

- Detectors are pure functions over files and must not execute repository code.
- Malformed manifests produce a structured error and no requirement rather than
  a guessed one.
- Repos that declare nothing yield an empty finding set without error.

## Implementation References

- Planned: `internal/detect/detector.go`, `internal/detect/registry.go`,
  `internal/detect/dotfiles.go`, `internal/detect/manifests.go`,
  `internal/detect/textfiles.go`, fixtures under `testdata/detect/`.
- Requirements: `docs/features/scan-and-detect.md#SCAN-FR-07` through
  `docs/features/scan-and-detect.md#SCAN-FR-09`, `docs/PRD.md#CONST-05`.
