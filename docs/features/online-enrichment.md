# Feature: Online Enrichment

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| ENRICH-FR-01 | This feature | owns |
| ENRICH-FR-02 | This feature | owns |
| ENRICH-FR-03 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Online Enrichment
**ID Prefix:** ENRICH
**Summary:** Optional, additive online metadata — latest upstream tool versions and
documentation-link validation — that never blocks, invalidates, or changes the
local deterministic report.
**Dependencies:** Project Foundation, Tool Catalog, Probing and Version Matching
**Priority:** Should

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| ENRICH-US-01 | newcomer developer | to see the latest upstream version next to a requirement | I know whether a pin is stale | Should |
| ENRICH-US-02 | offline user | the tool to work fully without network | I am not blocked on connectivity | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"ENRICH-FR-01","kind":"requirement","text":"Online enrichment is disabled by default and runs only when the user passes --enrich, so that no outbound network request occurs during a default invocation."}
```

```forge-requirement
{"id":"ENRICH-FR-02","kind":"requirement","text":"When --enrich is enabled, look up the latest upstream version per tool and validate documentation links, and on any failure note the missing enrichment inline and continue without blocking or invalidating the local report."}
```

```forge-requirement
{"id":"ENRICH-FR-03","kind":"requirement","text":"Enrichment requests use a bounded per-request timeout and an overall time budget, are additive metadata only, and never change finding statuses or process exit codes."}
```

---

## 4. UI / Interaction Design

Enrichment appears as inline notes in the report (for example "latest: 20.11.1")
and as an explicit online/offline indicator; absence is stated, not hidden.

---

## 5. Implementation Tasks

### Phase 1: Enrichment Client and Link Validation

```forge-task
{
  "id": "ENRICH-1",
  "title": "Implement the enrichment client and latest-version lookup",
  "description": "Implement internal/enrich with an HTTP client that performs latest-version lookups for catalog tools with a bounded per-request timeout and an overall budget, plus a function that returns an Enrichment for each finding or a structured not-available marker. Every network error, timeout, or non-200 response degrades to not-available without returning a fatal error. The package must make no requests when not explicitly invoked.",
  "ownerAgent": "enrichment-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/enrich/client.go", "internal/enrich/latest.go", "internal/enrich/client_test.go"],
  "validationCommands": ["go test ./internal/enrich/... -run TestLatest", "go vet ./internal/enrich/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/online-enrichment.md#ENRICH-FR-01", "docs/features/online-enrichment.md#ENRICH-FR-02", "docs/features/online-enrichment.md#ENRICH-FR-03"],
    "acceptanceCriteria": ["No request is made unless the lookup is invoked", "Timeouts and error responses degrade to not-available", "A successful lookup returns the latest version for the tool", "Tests use an httptest server and run offline"],
    "constraints": ["Only tool identifiers and documentation URLs leave the machine; no repository contents are sent"],
    "constraintRefs": ["docs/PRD.md#CONST-10", "docs/PRD.md#CONST-08"],
    "references": []
  }
}
```

```forge-task
{
  "id": "ENRICH-2",
  "title": "Validate documentation links and attach inline enrichment notes",
  "description": "Implement documentation-link validation for catalog entries and a function that attaches enrichment results to findings as additive notes, including an explicit not-available note when a lookup could not complete. Ensure the attachment step never mutates finding status, never removes local data, and is safe to skip entirely when enrichment is disabled.",
  "ownerAgent": "enrichment-engineer",
  "dependencies": ["ENRICH-1"],
  "expectedOutputs": ["internal/enrich/links.go", "internal/enrich/report.go", "internal/enrich/links_test.go"],
  "validationCommands": ["go test ./internal/enrich/... -run TestLinks", "go vet ./internal/enrich/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/online-enrichment.md#ENRICH-FR-02"],
    "acceptanceCriteria": ["A broken link is reported as a note and does not fail the run", "Attaching enrichment leaves statuses and exit codes unchanged", "A disabled enricher produces no notes and no requests", "Tests run offline with an httptest server"],
    "constraints": ["Enrichment is additive and degrades gracefully"],
    "constraintRefs": ["docs/PRD.md#CONST-08", "docs/PRD.md#CONST-10"],
    "references": []
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Unit Tests | latest lookup, link validation, note attachment | `httptest` servers with success, timeout, and error cases |
| Degradation Tests | offline and failure paths | Assert not-available notes and unchanged statuses |

Key test scenarios:
1. A default invocation performs zero network requests.
2. A timed-out lookup produces a not-available note and a successful scan.
3. A broken docs link is noted without changing any status.

---

## 7. Acceptance Criteria

1. Enrichment is opt-in and performs no network I/O by default.
2. All enrichment failures degrade to inline notes and never change exit codes.
3. No repository contents are transmitted.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Which upstream source supplies latest versions | A single configurable HTTPS endpoint with graceful fallback |
| 2 | Whether link validation is on for every --enrich run | On by default when --enrich is set, bounded by budget |
