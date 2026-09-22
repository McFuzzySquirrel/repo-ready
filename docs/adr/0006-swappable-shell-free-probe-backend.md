# ADR-0006: Swappable, shell-free, bounded probe backend

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** probe-engineer

## Context

Discovering installed versions could be done by invoking a shell, by running
`--version`, or by enumerating version managers (nvm, asdf, mise, pyenv). Shell
interpolation of tool names is an injection risk, and unbounded parallel
probing can hang or flood the machine. v1 intentionally excludes version-manager
enumeration but should not have to be rewritten to add it.

## Decision

Expose installed-tool discovery behind a **backend interface** so enumeration
can be added later without changing callers (PROBE-FR-01). The v1 backend
locates a tool with `exec.LookPath`, runs its known safe version-query
invocation directly **without a shell** under a 5-second context timeout, and
parses the first version-like token from the output (PROBE-FR-02). A missing
binary yields a `missing` result rather than an error; an unrecognized binary
yields an `unknown` version. Probe scheduling de-duplicates tool identifiers and
runs at most 8 concurrent workers, and cancellation stops outstanding probes
(PROBE-FR-03, CONST-11). Tool names are validated against the known catalog and
never interpolated into a shell string (CONST-09).

## Alternatives Considered

- **Shell out to `bash -c "<tool> --version"`.** Rejected: shell injection and
  platform inconsistency.
- **Enumerate version managers in v1.** Rejected: additional scope; the backend
  interface defers it without a rewrite.
- **Unbounded parallelism.** Rejected: a single misbehaving binary could stall
  the run; bounded workers plus timeouts keep the budget.

## Consequences

- Only known safe version-query commands ever execute, with no shell.
- A tool-to-command map must be explicit and validated.
- v1 does not report versions installed only via a version manager that is not
  on `PATH`; this limitation is documented.

## Implementation References

- Planned: `internal/probe/backend.go`, `internal/probe/path.go`,
  `internal/probe/parse.go`, `internal/probe/scheduler.go`.
- Requirements: `docs/features/probing-and-matching.md#PROBE-FR-01` through
  `docs/features/probing-and-matching.md#PROBE-FR-03`, `docs/PRD.md#CONST-09`,
  `docs/PRD.md#CONST-11`.
