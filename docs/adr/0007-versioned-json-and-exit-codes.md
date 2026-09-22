# ADR-0007: Stable versioned JSON and a fixed exit-code contract

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** cli-engineer

## Context

CI pipelines and scripts need machine-readable output and predictable exit
codes. Interactive users need a TUI. If the JSON shape changes between releases,
consumers break; if the exit code conflates "scan succeeded but tools are
missing" with "the scan failed", callers cannot distinguish incomplete
environments from operational failures.

## Decision

`--json` emits a **single versioned object** containing `schemaVersion`, input
metadata, scan metadata (including skipped paths and depth), components,
findings with declared and installed values, conflict declarations, and an
optional enrichment section that is omitted when absent; field names are stable
across releases (REPORT-FR-03). Exit codes are fixed: `0` on a successful scan
(even with missing tools), `1` only when `--strict` is set and at least one
finding is missing or outdated, and `2` on a fatal input, clone, or usage error
(REPORT-FR-02). Invalid usage exits `2` (REPORT-FR-01).

## Alternatives Considered

- **Exit non-zero whenever a tool is missing.** Rejected: it would make the
  default report unusable in scripts and conflate data with failure.
- **Unversioned JSON.** Rejected: it offers no compatibility contract.
- **Emit JSON logs/streams.** Rejected: a single document is easier for callers
  to consume and golden-test.

## Consequences

- `--json` without `--strict` exits `0` whenever the scan succeeds; the caller
  decides what the data means.
- The schema version is a constant, and a golden test pins field names.
- The TUI and JSON share one `Report` model so they cannot drift.

## Implementation References

- Planned: `internal/report/json.go`, `internal/app/cli.go`,
  `cmd/repo-ready/main.go`, and a future `docs/JSON-SCHEMA.md` describing the
  emitted fields.
- Requirements: `docs/features/reporting.md#REPORT-FR-01` through
  `docs/features/reporting.md#REPORT-FR-03`.
