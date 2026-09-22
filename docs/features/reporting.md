# Feature: Reporting

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| REPORT-FR-01 | This feature | owns |
| REPORT-FR-02 | This feature | owns |
| REPORT-FR-03 | This feature | owns |
| REPORT-FR-04 | This feature | owns |
| REPORT-FR-05 | This feature | owns |
| REPORT-FR-06 | This feature | owns |
| REPORT-FR-07 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Reporting
**ID Prefix:** REPORT
**Summary:** Wire the pipeline into a single run and present results as an
interactive, keyboard-driven TUI or a versioned JSON document, with a
screen-reader-friendly plain mode and strict exit-code semantics.
**Dependencies:** Project Foundation, Scan and Detection, Tool Catalog, Probing and Version Matching, Online Enrichment
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| REPORT-US-01 | newcomer developer | a scrollable report with per-tool detail and a copy action | I can act on the output immediately | Must |
| REPORT-US-02 | CI engineer | stable JSON and a strict exit code | I can gate a pipeline | Must |
| REPORT-US-03 | screen-reader user | a plain mode without color or box drawing | I can read the report | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"REPORT-FR-01","kind":"requirement","text":"The CLI accepts one positional path-or-URL argument plus the flags --json, --strict, --depth, --no-color, --plain, --enrich, --version, and --help, and rejects invalid usage with exit code 2."}
```

```forge-requirement
{"id":"REPORT-FR-02","kind":"requirement","text":"Exit code is 0 on a successful scan, 1 only when --strict is set and at least one finding is missing or outdated, and 2 on a fatal input, clone, or usage error."}
```

```forge-requirement
{"id":"REPORT-FR-03","kind":"requirement","text":"--json emits a single versioned JSON object containing schemaVersion, input metadata, scan metadata including skipped paths, components, findings with declared and installed values, and optional enrichment, all with stable field names."}
```

```forge-requirement
{"id":"REPORT-FR-04","kind":"requirement","text":"The TUI presents a scrollable report, a per-tool detail pane showing why the tool is required, the declaring file, the documentation link, and the install commands, and filters for all tools, missing-or-outdated tools, and category."}
```

```forge-requirement
{"id":"REPORT-FR-05","kind":"requirement","text":"The TUI is fully keyboard-driven, copies the selected install command via OSC52, and shows a clear copy-unsupported notice when the terminal does not support OSC52."}
```

```forge-requirement
{"id":"REPORT-FR-06","kind":"requirement","text":"Status is never conveyed by color alone, the report honors --no-color and the NO_COLOR environment variable, and --plain renders screen-reader-friendly text without box-drawing characters."}
```

```forge-requirement
{"id":"REPORT-FR-07","kind":"requirement","text":"A single orchestrator runs input, scan, detect, probe, and match in order, attaches enrichment only when requested, and surfaces skipped paths and the online/offline enrichment state in the report."}
```

---

## 4. UI / Interaction Design

- **Report list:** one row per finding with status label, tool, declared
  constraint, and installed version; color is redundant with a text status.
- **Detail pane:** why required, declaring file and component, docs link, and
  per-platform install commands.
- **Filters:** all, missing-or-outdated, and category; a persistent count is shown.
- **Keys:** arrow keys and j/k navigate, Enter opens detail, f cycles filters,
  c copies the selected install command, q quits. Mouse is optional.
- **Plain mode:** linear text, no box drawing, no color, stable reading order.

---

## 5. Implementation Tasks

### Phase 1: Encoders, Orchestration, and TUI

```forge-task
{
  "id": "REPORT-1",
  "title": "Encode the versioned JSON report",
  "description": "Implement internal/report with a JSON encoder that renders the Report model as a single object with schemaVersion, input metadata, scan metadata including skipped paths and depth, components, findings with declared and installed values and conflict declarations, and an optional enrichment section that is omitted when absent. Enforce stable field ordering and names with a schema version constant, and return an error only for unencodable input.",
  "ownerAgent": "cli-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/report/json.go", "internal/report/json_test.go"],
  "validationCommands": ["go test ./internal/report/...", "go vet ./internal/report/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-03"],
    "acceptanceCriteria": ["Output is a single object with a schemaVersion field", "Skipped paths and scan depth are present", "Enrichment is omitted when absent and present when supplied", "A golden test pins the field names"],
    "constraints": [],
    "constraintRefs": [],
    "references": []
  }
}
```

```forge-task
{
  "id": "REPORT-2",
  "title": "Orchestrate the scan, detect, probe, and match pipeline",
  "description": "Implement the orchestrator in internal/app that runs input resolution, scanning, detection, probing, and reconciliation in order, assembles the Report model, and records skipped paths plus whether enrichment is enabled and online. Accept an optional enricher hook that is nil by default so no enrichment code runs unless supplied. Fail with a typed fatal error for input or clone failures. Do not parse CLI flags or render output here.",
  "ownerAgent": "cli-engineer",
  "dependencies": ["SCAN-6", "CATALOG-1", "PROBE-4", "REPORT-1"],
  "expectedOutputs": ["internal/app/orchestrate.go", "internal/app/orchestrate_test.go"],
  "validationCommands": ["go test ./internal/app/... -run TestOrchestrate", "go vet ./internal/app/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-07"],
    "acceptanceCriteria": ["A fixture repo produces a complete Report with components and findings", "Skipped paths and enrichment state are recorded", "A nil enricher performs no enrichment work", "An input failure returns a typed fatal error"],
    "constraints": ["The orchestrator never installs or executes repository code"],
    "constraintRefs": ["docs/PRD.md#CONST-01"],
    "references": []
  }
}
```

```forge-task
{
  "id": "REPORT-3",
  "title": "Build the Bubble Tea report and detail view",
  "description": "Implement internal/tui with a Bubble Tea model, view, and update that render a scrollable findings list and a per-tool detail pane. Each row shows a text status label, tool, declared constraint, and installed version. The detail pane shows why the tool is required, the declaring file and component, the docs link, and per-platform install commands. Navigation is via arrow keys and j/k, Enter opens detail, and q quits. Styling uses Lip Gloss.",
  "ownerAgent": "tui-engineer",
  "dependencies": ["REPORT-2"],
  "expectedOutputs": ["internal/tui/model.go", "internal/tui/view.go", "internal/tui/update.go", "internal/tui/tui_test.go"],
  "validationCommands": ["go test ./internal/tui/... -run TestModel", "go vet ./internal/tui/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-04"],
    "acceptanceCriteria": ["The view renders a row per finding with status, declared, and installed values", "The detail pane shows requirement reason, declaring file, docs link, and install commands", "Navigation and detail toggling update the model as asserted by update tests"],
    "constraints": ["Status is never rendered by color alone"],
    "constraintRefs": ["docs/PRD.md#CONST-12"],
    "references": []
  }
}
```

```forge-task
{
  "id": "REPORT-4",
  "title": "Add filters, OSC52 copy, and color/plain modes",
  "description": "Extend the TUI with filters for all, missing-or-outdated, and category, a persistent visible count, OSC52 clipboard copy of the selected install command with a copy-unsupported notice, and color handling that honors --no-color and NO_COLOR plus a --plain mode that removes box drawing and color. Keep keyboard shortcuts documented in the view footer.",
  "ownerAgent": "tui-engineer",
  "dependencies": ["REPORT-3"],
  "expectedOutputs": ["internal/tui/filter.go", "internal/tui/copy.go", "internal/tui/style.go", "internal/tui/filter_test.go"],
  "validationCommands": ["go test ./internal/tui/... -run TestFilters", "go vet ./internal/tui/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-04", "docs/features/reporting.md#REPORT-FR-05", "docs/features/reporting.md#REPORT-FR-06"],
    "acceptanceCriteria": ["Each filter yields the expected visible subset", "Copy emits the correct OSC52 sequence and shows an unsupported notice when unavailable", "NO_COLOR and --no-color suppress color and --plain removes box drawing", "Status remains readable without color"],
    "constraints": ["Keyboard-only operation must be complete; no action requires a mouse"],
    "constraintRefs": ["docs/PRD.md#CONST-12"],
    "references": []
  }
}
```

```forge-task
{
  "id": "REPORT-5",
  "title": "Implement CLI flags, wiring, and exit codes",
  "description": "Implement the CLI entrypoint and flag parsing in cmd/repo-ready and internal/app/cli.go. Accept one positional path-or-URL and the flags --json, --strict, --depth, --no-color, --plain, --enrich, --version, and --help. Wire the orchestrator to either the JSON encoder or the TUI, construct the enricher only when --enrich is set, print usage on invalid usage, and return exit codes 0 on success, 1 only under --strict with a missing or outdated finding, and 2 on fatal input, clone, or usage error. Keep main.go minimal.",
  "ownerAgent": "cli-engineer",
  "dependencies": ["REPORT-2", "REPORT-4"],
  "expectedOutputs": ["cmd/repo-ready/main.go", "internal/app/cli.go", "internal/app/cli_test.go"],
  "validationCommands": ["go test ./internal/app/... -run TestCLI", "go build ./..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-01", "docs/features/reporting.md#REPORT-FR-02"],
    "acceptanceCriteria": ["Each documented flag is parsed and invalid usage exits 2", "--json selects JSON and the default selects the TUI", "--strict returns 1 only with a missing or outdated finding", "A good scan without --strict returns 0"],
    "constraints": ["--enrich is the only path that enables network enrichment"],
    "constraintRefs": ["docs/PRD.md#CONST-08"],
    "references": []
  }
}
```

### Phase 2: Human Review

```forge-task
{
  "id": "REPORT-6",
  "title": "Human review of TUI legibility and pass/warn/fail readability",
  "description": "A human reviewer runs the TUI against a fixture repository and evaluates legibility: whether status reads clearly without color, whether declared and installed versions are unambiguous, whether the detail pane is sufficient to act, and whether filters and keys are discoverable. Capture screenshots or transcripts, list any requested changes with severity, and record an explicit approval decision in the review file. This judgment cannot be made by an agent.",
  "dependencies": ["REPORT-4", "REPORT-5"],
  "expectedOutputs": [],
  "validationCommands": [],
  "contract": {
    "version": 2,
    "kind": "human-review",
    "requirements": [],
    "requirementRefs": ["docs/features/reporting.md#REPORT-FR-04"],
    "acceptanceCriteria": ["A transcript or screenshot set is captured for list, detail, filter, and plain modes", "Requested changes are listed with severity and disposition", "The reviewer records an explicit approval decision"],
    "constraints": [],
    "constraintRefs": [],
    "references": [],
    "reviewFile": "docs/reviews/tui-legibility.json"
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Unit Tests | JSON encoding, orchestrator, TUI update logic, filters, copy | Table-driven and golden tests; Bubble Tea update assertions |
| Integration Tests | CLI end-to-end on a fixture repo | Build and run the binary; assert output and exit codes |
| Human Review | TUI legibility | Recorded transcript and approval decision |

Key test scenarios:
1. `--json` on a fixture repo emits a single valid object with the expected fields.
2. `--strict` exits 1 when a tool is missing and 0 otherwise.
3. `--plain` output contains no ANSI color codes or box-drawing characters.

---

## 7. Acceptance Criteria

1. The CLI implements every documented flag and the 0/1/2 exit-code contract.
2. JSON output is versioned and stable; TUI is keyboard-complete and accessible.
3. The TUI legibility human-review gate is recorded and approved.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Exact key bindings and footer wording | The set documented above |
| 2 | Whether enrichment state shows in JSON as well as the TUI | Include an enrichment section when enrichment runs |
