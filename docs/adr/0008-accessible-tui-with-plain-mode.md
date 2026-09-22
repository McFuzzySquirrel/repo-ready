# ADR-0008: Keyboard-driven TUI with an accessible plain mode

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** tui-engineer

## Context

The primary interface is an interactive terminal report. Color and box drawing
are useful visually but are inaccessible to screen readers and to terminals
with limited support. Users must be able to operate the report entirely without
a mouse, and the meaning of a status must survive when color is unavailable.

## Decision

Build the TUI with Bubble Tea, Lip Gloss, and Bubbles. Present a scrollable
findings list and a per-tool detail pane showing why the tool is required, the
declaring file, the documentation link, and per-platform install commands, with
filters for all, missing-or-outdated, and category (REPORT-FR-04). Operation is
fully keyboard-driven (arrows or `j`/`k`, `Enter`, `f`, `c`, `q`), and the
selected install command is copied via OSC52 with a clear "copy unsupported"
notice when the terminal cannot comply (REPORT-FR-05). Status is **never**
conveyed by color alone; the report honors `--no-color` and `NO_COLOR`, and
`--plain` renders screen-reader-friendly text without box-drawing characters
(REPORT-FR-06, CONST-12).

## Alternatives Considered

- **Color-only status.** Rejected on accessibility grounds.
- **Mouse-first interaction.** Rejected: keyboard completeness is required; mouse
  is optional.
- **A separate plain renderer with its own model.** Rejected: one model with
  presentation modes avoids drift.

## Consequences

- A human "TUI legibility" review gate evaluates status readability without
  color, declared/installed clarity, and key discoverability.
- Plain mode must contain no ANSI color codes and no box-drawing characters.

## Implementation References

- Planned: `internal/tui/model.go`, `internal/tui/view.go`,
  `internal/tui/update.go`, `internal/tui/filter.go`, `internal/tui/copy.go`,
  `internal/tui/style.go`.
- Requirements: `docs/PRD.md#CONST-12`,
  `docs/features/reporting.md#REPORT-FR-04` through
  `docs/features/reporting.md#REPORT-FR-06`.
