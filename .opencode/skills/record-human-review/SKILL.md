---
name: record-human-review
description: "Prepare and record a repo-ready human-review gate: assemble per-category catalog sample or TUI transcript evidence, capture observed results with the JSON review-file schema under docs/reviews/, list requested corrections with dispositions, and record an explicit reviewer decision without self-approving."
---

# Skill: Record a Human-Review Gate

Prepare and record a repo-ready human-review gate as durable JSON evidence. Use
this for the catalog accuracy gate at `docs/reviews/catalog-accuracy.json`
(CATALOG-3) or the TUI legibility gate at `docs/reviews/tui-legibility.json`
(REPORT-6). These gates judge accuracy and legibility that an automated check
cannot decide, so an agent assembles evidence and only a human approves.

## Process

### Step 1: Identify the gate and its file

Match the request to a known gate: catalog accuracy or TUI legibility. If the
request is a different gate, then treat it as a new schema, record a schema
version, and do not key off an unrelated review file.

### Step 2: Assemble the evidence

For catalog accuracy, sample across every category and capture the resolved docs
URL plus the Linux, macOS, and Windows commands for each sampled tool. For TUI
legibility, capture transcripts or screenshots for the list, detail, filter, and
plain modes. Load the review-file schema from `references/review-schema.md` when
you need the exact field names for a gate.

### Step 3: Record observed results with the schema

Write a single JSON object under `docs/reviews/` with stable fields: gate id,
schema version, reviewer, timestamp, sampled ids or modes, observed results,
requested corrections, dispositions, and the decision. Keep field names identical
to the schema so downstream tooling can read the file.

### Step 4: List corrections with dispositions

For each requested correction, record a severity and a disposition: applied,
accepted with a rationale, or deferred. If a correction is applied, then re-run the
relevant catalog or TUI test and update the observed result before recording the
disposition.

### Step 5: Record the reviewer decision without self-approving

Only a human reviewer sets the decision to approved. An agent records evidence and
leaves the decision pending by default. If no human has reviewed the sample, then the gate
stays pending; never infer approval from a passing test.

### Step 6: Validate and hand off

Validate the file parses, confirm the catalog or TUI source was not edited to
match a fabricated sample, and report the gate status. Do not modify
`docs/SKILL-CANDIDATES.json`.

## Gotchas

- **An agent cannot approve a human gate.** Marking the decision approved without a human review is the exact failure this gate exists to prevent (CATALOG-3, REPORT-6).
- **Sample across every category.** A catalog review that only covers runtimes leaves package managers and cloud CLIs unverified.
- **Capture all four TUI modes.** List, detail, filter, and plain; a missing plain-mode transcript leaves CONST-12 unproven.
- **Record what was observed, not what was expected.** Copy the actual command output or URL resolution result into the observed field.
- **Dispositions are mandatory.** A correction with no disposition reads as ignored and blocks the gate.
- **Keep the JSON field names stable.** The review schema is consumed by later stages; renaming a field breaks reconciliation.
- **Do not fix the sample.** Editing `tools.json` or `internal/tui` to match the review after the fact invalidates the evidence.

## Validation

Run this self-check before recording the gate:

- [ ] The review file parses cleanly, for example `python3 -m json.tool docs/reviews/<file>.json`
- [ ] The sample spans every category (catalog) or all four modes (TUI)
- [ ] Every requested correction has a severity and a disposition
- [ ] The decision is left pending unless a human reviewer set it
- [ ] The catalog or TUI source was not edited to match the sample

The reference `references/review-schema.md` holds both gate schemas when you need
to pin the fields.
