# ADR-0005: Loose, transparent version matching that warns

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** probe-engineer

## Context

Repositories declare constraints in inconsistent forms: bare pins (`18.2.0`),
minimums (`>=18.2.0`), caret ranges (`^18.2.0`), tilde ranges (`~18.2.0`,
`~>18.2.0`), and bare majors (`18`). A strict resolver would reject
technically-satisfiable environments and block a developer on a version the
repo would actually accept.

## Decision

Treat a bare numeric pin as a minimum, `>=X` as a minimum, `^X` as the same
major as X, `~X`/`~>X` as the same minor as X, and a bare major as the same
major (PROBE-FR-04). Each finding carries one of `ok`, `missing`, `outdated`,
or `unknown`; a mismatch **warns** rather than silently failing, and the
declared constraint and installed version strings are always retained and shown
together (PROBE-FR-05). An unparseable constraint or version is `unknown`, never
a failure.

## Alternatives Considered

- **Strict semver satisfaction.** Rejected: too many false "outdated" results on
  real repos that use ranges loosely.
- **Only compare exact equality.** Rejected: it would mark almost every
  environment as outdated.
- **Hide the raw strings and show a status only.** Rejected: the product
  promises correctness-first, transparency over opaque status.

## Consequences

- Reconciliation must de-duplicate by tool across components and, when the same
  tool is declared with differing constraints, list every declaration and flag
  the finding as a conflict (PROBE-FR-06).
- Only `--strict` converts a missing/outdated finding into a non-zero exit.
- Boundary versions (same major/minor edges) require table-driven tests.

## Implementation References

- Planned: `internal/match/constraint.go`, `internal/match/match.go`,
  `internal/match/reconcile.go`, `internal/match/match_test.go`.
- Requirements: `docs/features/probing-and-matching.md#PROBE-FR-04` through
  `docs/features/probing-and-matching.md#PROBE-FR-06`.
