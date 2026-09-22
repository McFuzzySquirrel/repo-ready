# Feature: End-to-End Validation

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| E2E-FR-01 | This feature | owns |
| E2E-FR-02 | This feature | owns |
| E2E-FR-03 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** End-to-End Validation
**ID Prefix:** E2E
**Summary:** Prove the whole pipeline against real fixture repositories offline,
provide an opt-in corpus test against remote repositories, and enforce a time
budget for a typical scan.
**Dependencies:** Reporting, Scan and Detection, Probing and Version Matching
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| E2E-US-01 | maintainer | deterministic end-to-end tests without network | CI is reliable | Must |
| E2E-US-02 | maintainer | an opt-in corpus against real repos | integration regressions are caught | Should |

---

## 3. Functional Requirements

```forge-requirement
{"id":"E2E-FR-01","kind":"requirement","text":"Offline end-to-end tests build or invoke the binary against committed fixture repositories and assert deterministic scan, detection, matching, and reporting output with no network access."}
```

```forge-requirement
{"id":"E2E-FR-02","kind":"requirement","text":"An opt-in network corpus test runs against a small set of real remote repositories and is skipped by default so that standard CI remains deterministic and offline."}
```

```forge-requirement
{"id":"E2E-FR-03","kind":"requirement","text":"A performance test asserts that a typical fixture repository completes a full scan and report within the documented time budget."}
```

---

## 4. UI / Interaction Design

None; this feature validates behavior established elsewhere.

---

## 5. Implementation Tasks

### Phase 1: Offline End-to-End and Opt-In Corpus

```forge-task
{
  "id": "E2E-1",
  "title": "Build the offline end-to-end fixture harness",
  "description": "Implement offline end-to-end tests that run the CLI against committed fixture repositories under e2e/testdata, covering a single-ecosystem repo, a monorepo with multiple components and a version conflict, and a repo declaring nothing. Assert exit codes, JSON structure, component grouping, conflict flagging, and deterministic repeatability with no network access.",
  "ownerAgent": "qa-engineer",
  "dependencies": ["REPORT-5", "SCAN-6", "PROBE-4"],
  "expectedOutputs": ["e2e/e2e_test.go", "e2e/fixtures_test.go", "e2e/testdata/node-rust/.nvmrc", "e2e/testdata/monorepo/web/package.json", "e2e/testdata/monorepo/api/go.mod"],
  "validationCommands": ["go test ./e2e/... -run TestOffline", "go vet ./e2e/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/e2e-validation.md#E2E-FR-01"],
    "acceptanceCriteria": ["A single-ecosystem fixture produces the expected tools and statuses", "A monorepo fixture groups findings by component and flags the conflict", "A repo declaring nothing produces an empty finding set without error", "Repeated runs are byte-identical", "Tests make no network calls"],
    "constraints": ["Offline tests must not depend on network or on tools installed on the host"],
    "constraintRefs": ["docs/PRD.md#CONST-08"],
    "references": []
  }
}
```

```forge-task
{
  "id": "E2E-2",
  "title": "Add the opt-in network corpus and performance budget tests",
  "description": "Add an opt-in test that, only when an explicit environment variable is set, runs the tool against a small configured corpus of real remote repositories and checks that scanning succeeds and findings are produced. Add a performance test that times a full scan and report of a typical fixture repository and fails when it exceeds the documented budget. Both tests must skip cleanly by default so standard CI stays offline and deterministic.",
  "ownerAgent": "qa-engineer",
  "dependencies": ["E2E-1"],
  "expectedOutputs": ["e2e/network_test.go", "e2e/perf_test.go"],
  "validationCommands": ["go test ./e2e/... -run 'TestNetwork|TestPerf'"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/e2e-validation.md#E2E-FR-02", "docs/features/e2e-validation.md#E2E-FR-03"],
    "acceptanceCriteria": ["The network corpus test skips unless its environment variable is set", "A real remote repo scans successfully when enabled", "The performance test fails when a fixture exceeds the time budget", "Both tests are safe to run repeatedly"],
    "constraints": ["Standard CI must remain offline; network tests are opt-in only"],
    "constraintRefs": ["docs/PRD.md#CONST-08", "docs/PRD.md#CONST-11"],
    "references": []
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Offline E2E | full pipeline against fixture repos | Build/invoke the binary; assert output and codes |
| Opt-in Corpus | real remote repositories | Env-gated; skipped by default |
| Performance | time budget for a typical repo | Timed run with an asserted threshold |

Key test scenarios:
1. The monorepo fixture reports two components and a conflict-flagged tool.
2. The empty-declaration fixture exits 0 with an empty finding list.
3. A fixture exceeding the time budget fails the performance test.

---

## 7. Acceptance Criteria

1. Offline end-to-end tests cover success, monorepo grouping, conflict, and empty cases.
2. Network corpus tests are opt-in and skipped by default.
3. A performance threshold test guards the time budget.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Which real repositories form the network corpus | Two or three small public repos configured in the test |
| 2 | Exact performance threshold | 2 seconds for the typical fixture on CI |
