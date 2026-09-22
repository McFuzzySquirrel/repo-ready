# Feature: Probing and Version Matching

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| PROBE-FR-01 | This feature | owns |
| PROBE-FR-02 | This feature | owns |
| PROBE-FR-03 | This feature | owns |
| PROBE-FR-04 | This feature | owns |
| PROBE-FR-05 | This feature | owns |
| PROBE-FR-06 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Probing and Version Matching
**ID Prefix:** PROBE
**Summary:** Query the local machine for installed tools behind a swappable
backend, then compare installed versions to declared constraints loosely and
transparently, de-duplicating and flagging conflicts.
**Dependencies:** Project Foundation
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| PROBE-US-01 | newcomer developer | to see "have X, need Y" per tool | I can judge the comparison myself | Must |
| PROBE-US-02 | monorepo contributor | conflicting declarations surfaced | I know which component disagrees | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"PROBE-FR-01","kind":"requirement","text":"Expose installed-tool probing behind a backend interface so version-manager enumeration can be added later without changing callers."}
```

```forge-requirement
{"id":"PROBE-FR-02","kind":"requirement","text":"The v1 backend locates a tool on PATH, runs a known safe version-query invocation without a shell with a 5 second timeout, parses the first version-like token from its output, and reports a missing binary as missing rather than an error."}
```

```forge-requirement
{"id":"PROBE-FR-03","kind":"requirement","text":"Probe requests run with at most 8 concurrent workers and identical tool probes are de-duplicated so each distinct tool is queried once."}
```

```forge-requirement
{"id":"PROBE-FR-04","kind":"requirement","text":"Loose matching parses a numeric pin as a minimum, >=X as a minimum, ^X as the same major as X, ~X or ~>X as the same minor as X, and a bare major as the same major, and treats an unparseable constraint as unknown."}
```

```forge-requirement
{"id":"PROBE-FR-05","kind":"requirement","text":"Each finding has a status of ok, missing, outdated, or unknown; a mismatch warns rather than silently failing, and the declared constraint and installed version are always retained and shown together."}
```

```forge-requirement
{"id":"PROBE-FR-06","kind":"requirement","text":"Findings are de-duplicated by tool across components, and when a tool is declared with differing constraints in multiple components all declarations are listed and the finding is flagged as a conflict."}
```

---

## 4. UI / Interaction Design

Status and raw declared/installed strings produced here drive the Reporting list
and detail pane.

---

## 5. Implementation Tasks

### Phase 1: Probing Backend and Matching

```forge-task
{
  "id": "PROBE-1",
  "title": "Implement the swappable probe backend with PATH lookup",
  "description": "Define the probe backend interface in internal/probe, then implement the v1 backend that locates a tool with exec.LookPath, runs its known safe version-query invocation directly without a shell under a 5 second context timeout, and parses the first version-like token from combined output. A binary that is absent yields a missing result rather than an error, and an unrecognized binary yields an unknown version. Keep tool-to-command mapping explicit and validated against known tool identifiers.",
  "ownerAgent": "probe-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/probe/backend.go", "internal/probe/path.go", "internal/probe/parse.go", "internal/probe/probe_test.go"],
  "validationCommands": ["go test ./internal/probe/... -run TestPathBackend", "go vet ./internal/probe/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/probing-and-matching.md#PROBE-FR-01", "docs/features/probing-and-matching.md#PROBE-FR-02"],
    "acceptanceCriteria": ["A missing binary yields a missing result, not an error", "The invocation runs without a shell and honors the 5 second timeout", "Version output variants are parsed by table-driven tests", "The backend is used only through its interface"],
    "constraints": ["Only known safe version-query commands run, with no shell and no credential access"],
    "constraintRefs": ["docs/PRD.md#CONST-09", "docs/PRD.md#CONST-01"],
    "references": []
  }
}
```

```forge-task
{
  "id": "PROBE-2",
  "title": "Add bounded-concurrency probe scheduling with de-duplication",
  "description": "Implement a probe scheduler that accepts a set of tool identifiers, de-duplicates them, runs probes with at most 8 concurrent workers, respects each probe's timeout, and returns results keyed by tool. Cancellation stops outstanding work promptly. Do not implement matching or reporting here.",
  "ownerAgent": "probe-engineer",
  "dependencies": ["PROBE-1"],
  "expectedOutputs": ["internal/probe/scheduler.go", "internal/probe/scheduler_test.go"],
  "validationCommands": ["go test ./internal/probe/... -run TestScheduler", "go vet ./internal/probe/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/probing-and-matching.md#PROBE-FR-03"],
    "acceptanceCriteria": ["A tool present twice in the input is probed once", "No more than 8 probes run concurrently", "Cancellation terminates outstanding probes"],
    "constraints": ["Bound concurrency to the documented worker limit"],
    "constraintRefs": ["docs/PRD.md#CONST-11"],
    "references": []
  }
}
```

```forge-task
{
  "id": "PROBE-3",
  "title": "Parse constraints and compute loose match statuses",
  "description": "Implement internal/match constraint parsing and status computation. Support numeric pins as minimums, >=X minimums, ^X same-major, ~X and ~>X same-minor, and bare major same-major, using the pinned semver library for parsing. Compute ok, missing, outdated, or unknown, where an unparseable constraint or version is unknown and never fails. Always retain the declared constraint string and the installed version string on the finding.",
  "ownerAgent": "probe-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/match/constraint.go", "internal/match/match.go", "internal/match/match_test.go"],
  "validationCommands": ["go test ./internal/match/... -run TestMatch", "go vet ./internal/match/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/probing-and-matching.md#PROBE-FR-04", "docs/features/probing-and-matching.md#PROBE-FR-05"],
    "acceptanceCriteria": ["Each constraint form is covered by a table-driven test including boundary versions", "An unparseable constraint yields unknown, never a failure", "Declared and installed strings are preserved verbatim on the finding"],
    "constraints": ["Mismatches warn; they never silently fail"],
    "constraintRefs": [],
    "references": []
  }
}
```

### Phase 2: Reconciliation

```forge-task
{
  "id": "PROBE-4",
  "title": "Reconcile findings across components and flag conflicts",
  "description": "Implement reconciliation that merges probe results with detected requirements, de-duplicates by tool across components, and, when one tool is declared with differing constraints in multiple components, produces a single finding that lists every declaration and is flagged as a conflict. Preserve component attribution for each declaration and keep deterministic ordering. Do not render output here.",
  "ownerAgent": "probe-engineer",
  "dependencies": ["PROBE-3"],
  "expectedOutputs": ["internal/match/reconcile.go", "internal/match/reconcile_test.go"],
  "validationCommands": ["go test ./internal/match/... -run TestReconcile", "go vet ./internal/match/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/probing-and-matching.md#PROBE-FR-06"],
    "acceptanceCriteria": ["A tool declared in three components yields one finding listing all declarations", "Differing constraints set the conflict flag and preserve every declaration", "Identical constraints do not set the conflict flag", "Reconciliation order is deterministic"],
    "constraints": [],
    "constraintRefs": [],
    "references": []
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Unit Tests | backend parsing, scheduler, constraint grammar, reconciliation | Table-driven tests; scheduler uses fake backends |
| Boundary Tests | version edges and invalid inputs | Explicit min/major/minor boundary cases |

Key test scenarios:
1. `^18.2.0` with installed `19.0.0` is outdated and with `18.9.1` is ok.
2. An unknown constraint or unparsable version is unknown and non-fatal.
3. The same tool declared with two constraints yields one conflict-flagged finding.

---

## 7. Acceptance Criteria

1. Probing is bounded, timed, shell-free, and swappable behind an interface.
2. Every constraint form has a documented, tested status outcome.
3. Conflicts surface all declarations with the conflict flag.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Which version-query command each tool uses | A curated map in the catalog; default `--version` |
| 2 | Behavior when the installed version is newer than a caret/major constraint | Same major passes; different major is outdated |
