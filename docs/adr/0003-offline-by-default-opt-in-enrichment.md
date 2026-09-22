# ADR-0003: Offline by default with opt-in additive enrichment

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** enrichment-engineer, probe-engineer

## Context

Developers often need an answer on a locked-down machine, in CI, or on a flaky
network. The value of the report is the local comparison between declared and
installed versions. Latest-upstream versions and live link checks are useful
but are not required to answer "what do I install?".

## Decision

`repo-ready` performs no network I/O by default. Online enrichment runs only
when the user passes `--enrich`, and it is strictly additive: it looks up the
latest upstream version per tool and validates documentation links, notes any
failure inline, and never blocks, invalidates, or changes the local report,
finding statuses, or exit codes (CONST-08, ENRICH-FR-01..03). When enrichment
is enabled, the only outbound data is tool identifiers and documentation URLs;
no repository contents and no telemetry leave the machine (CONST-10).

## Alternatives Considered

- **Enrich by default.** Rejected: it would break the offline guarantee and
  make the default path nondeterministic.
- **Cache enrichment results on disk.** Deferred: it would add state and a
  config surface that v1 explicitly excludes (CONST-05).
- **Treat a failed lookup as an error.** Rejected: enrichment must degrade
  gracefully; a missing feed must never block the deterministic report.

## Consequences

- Default runs are deterministic and reproducible, matching the byte-identical
  report requirement (SCAN-FR-10).
- The report must represent enrichment absence explicitly rather than hiding it.
- Enrichment needs bounded per-request timeouts and an overall budget.

## Implementation References

- Planned: `internal/enrich/client.go`, `internal/enrich/latest.go`,
  `internal/enrich/links.go`, `internal/enrich/report.go`.
- Requirements: `docs/PRD.md#CONST-08`, `docs/PRD.md#CONST-10`,
  `docs/features/online-enrichment.md#ENRICH-FR-01` through
  `docs/features/online-enrichment.md#ENRICH-FR-03`.
