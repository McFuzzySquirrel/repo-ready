# Architecture Decision Records

This directory records the durable architectural decisions for `repo-ready`.
Each record captures the context, the decision, the alternatives that were
rejected, the consequences, and the planned implementation references.

> [!NOTE]
> Every decision below is **Proposed**, not yet implemented. The repository has
> no Go source today; the decisions are ratified as design in
> [`docs/PRD.md`](../PRD.md) and the feature specifications under
> [`docs/features/`](../features/). Each ADR's Implementation References point
> to the planned files and the requirement IDs that will satisfy them.

| ADR | Title | Status | Requirement IDs |
|-----|-------|--------|-----------------|
| [0001](0001-single-static-go-binary.md) | Ship as a single static Go binary with an embedded catalog | Proposed | CONST-03, CONST-04 |
| [0002](0002-report-only-never-install.md) | Report-only: never install or execute repository code | Proposed | CONST-01, CONST-02 |
| [0003](0003-offline-by-default-opt-in-enrichment.md) | Offline by default with opt-in additive enrichment | Proposed | CONST-08, CONST-10, ENRICH-FR-01..03 |
| [0004](0004-declarative-detection-authoritative.md) | Declarative files are the authoritative source | Proposed | SCAN-FR-07..09, CONST-05 |
| [0005](0005-loose-transparent-version-matching.md) | Loose, transparent version matching that warns | Proposed | PROBE-FR-04, PROBE-FR-05 |
| [0006](0006-swappable-shell-free-probe-backend.md) | Swappable, shell-free, bounded probe backend | Proposed | PROBE-FR-01..03, CONST-09, CONST-11 |
| [0007](0007-versioned-json-and-exit-codes.md) | Stable versioned JSON and a fixed exit-code contract | Proposed | REPORT-FR-01..03 |
| [0008](0008-accessible-tui-with-plain-mode.md) | Keyboard-driven TUI with an accessible plain mode | Proposed | CONST-12, REPORT-FR-04..06 |
| [0009](0009-bounded-scan-and-safe-clone.md) | Bounded scanning and a safe remote-clone lifecycle | Proposed | SCAN-FR-02..06, SCAN-FR-10, CONST-06, CONST-13 |

## Conventions

- **Status** is one of `Proposed`, `Accepted`, `Superseded`, or `Deprecated`.
  A decision moves to `Accepted` when implementation and tests corroborate it.
- **Date** is the date the decision was recorded, not the date it shipped.
- IDs are stable and never reused; superseding a decision adds a new ADR and
  marks the old one `Superseded`.
