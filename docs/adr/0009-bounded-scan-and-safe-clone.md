# ADR-0009: Bounded scanning and a safe remote-clone lifecycle

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** scan-engineer

## Context

Scanning an arbitrary repository can descend into huge dependency trees, follow
into vendored code, traverse submodules, or modify the input. Remote inputs
introduce a temporary checkout that must never be left behind and must reuse the
user's existing credentials without handling them directly. A typical answer is
expected in under two seconds.

## Decision

Recursively scan to a default maximum depth of **6**, skipping `.git`,
submodule directories, a built-in denylist (`node_modules`, `vendor`, `.venv`,
`dist`, `build`, `target`, and similar), and paths matched by `.gitignore`;
v1 does not traverse submodules (SCAN-FR-05, CONST-06). Group requirements by
the component directory that declares them and record every skipped path with a
reason so the report can state that the scan was not exhaustive (SCAN-FR-06,
SCAN-FR-10). Local directories are scanned in place and fully offline without
modifying any file (SCAN-FR-02). Remote inputs are shallow-cloned with
`--depth 1` into a uniquely named `0700` temporary directory that is removed on
every exit path — success, error, and signal — and no scan ever modifies the
input repository (SCAN-FR-03, CONST-13). Clone operations reuse the ambient git
configuration and credential helpers and never prompt for or store credentials
(SCAN-FR-04, CONST-02).

## Alternatives Considered

- **Full-depth recursive scan with no limits.** Rejected: slow and noisy on real
  monorepos.
- **Ignore `.gitignore` and rely only on a denylist.** Rejected: it would scan
  generated and already-ignored trees.
- **Clone into the current directory or a fixed temp path.** Rejected: risks
  collisions and leaves artifacts behind; a unique `0700` dir with guaranteed
  cleanup is safer.

## Consequences

- An incomplete scan is explicit: skipped paths are part of both the TUI and the
  JSON report.
- Deterministic ordering is required so identical inputs produce byte-identical
  reports (SCAN-FR-10).
- The `.gitignore` matcher must be covered by fixture tests because it may
  diverge from git.

## Implementation References

- Planned: `internal/scan/scanner.go`, `internal/scan/gitignore.go`,
  `internal/scan/denylist.go`, `internal/input/resolve.go`,
  `internal/input/clone.go`, fixtures under `testdata/scan/` and `e2e/testdata/`.
- Requirements: `docs/features/scan-and-detect.md#SCAN-FR-02` through
  `docs/features/scan-and-detect.md#SCAN-FR-06`, `#SCAN-FR-10`,
  `docs/PRD.md#CONST-06`, `docs/PRD.md#CONST-13`.
