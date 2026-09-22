# Feature: Project Foundation

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| FOUND-FR-01 | This feature | owns |
| FOUND-FR-02 | This feature | owns |
| FOUND-FR-03 | This feature | owns |
| CONST-03 | [Vision](../PRD.md#15. Glossary) | participates |
| CONST-04 | [Vision](../PRD.md#15. Glossary) | participates |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Project Foundation
**ID Prefix:** FOUND
**Summary:** Establish the Go module, package skeleton, build tooling, and shared
domain model that every other feature compiles against.
**Dependencies:** None
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| FOUND-US-01 | developer | a single `make build` to produce the binary | I can build and run the tool locally without extra setup | Must |
| FOUND-US-02 | maintainer | one shared domain model | features integrate without redefining core types | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"FOUND-FR-01","kind":"requirement","text":"The repository is a Go module named github.com/mcfuzzysquirrel/repo-ready with a Go 1.27 toolchain directive, the cmd/repo-ready and internal/<domain> package layout, and a Makefile providing build, test, vet, and fmt targets."}
```

```forge-requirement
{"id":"FOUND-FR-02","kind":"requirement","text":"A shared internal/model package defines the domain types Requirement, Component, InstalledTool, Finding, Status with values ok|missing|outdated|unknown, and Report, with stable JSON tags that match the public --json schema."}
```

```forge-requirement
{"id":"FOUND-FR-03","kind":"requirement","text":"Build version metadata is injectable at build time via -ldflags into internal/version and defaults to the string dev when no version is injected."}
```

---

## 4. UI / Interaction Design

None for this feature; it provides build and domain foundations only.

---

## 5. Implementation Tasks

### Phase 1: Module and Domain Foundations

```forge-task
{
  "id": "FOUND-1",
  "title": "Initialize the Go module, layout, and build tooling",
  "description": "Create the Go module github.com/mcfuzzysquirrel/repo-ready with a go 1.27 toolchain directive and the cmd and internal package layout. Add a Makefile with build, test, vet, and fmt targets. Add a minimal cmd/repo-ready/main.go that prints the version, and internal/version with a Version variable defaulting to dev plus a test. Keep the module dependency-free at this stage.",
  "ownerAgent": "release-engineer",
  "dependencies": [],
  "expectedOutputs": ["go.mod", "Makefile", "cmd/repo-ready/main.go", "internal/version/version.go", "internal/version/version_test.go"],
  "validationCommands": ["go build ./...", "go test ./internal/version/...", "make build"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/project-foundation.md#FOUND-FR-01", "docs/features/project-foundation.md#FOUND-FR-03"],
    "acceptanceCriteria": ["go build ./... succeeds on a clean checkout", "The version defaults to dev and is overridable with -ldflags", "make build produces a runnable binary"],
    "constraints": [],
    "constraintRefs": ["docs/PRD.md#CONST-03", "docs/PRD.md#CONST-04"],
    "references": []
  }
}
```

```forge-task
{
  "id": "FOUND-2",
  "title": "Define the shared domain model",
  "description": "Implement internal/model with Requirement (tool, constraint, source file, component, category, confidence), Component, InstalledTool (tool, version, path), Finding (requirement, installed, Status), and Report. Status is an enumerated type with exactly the values ok, missing, outdated, and unknown, and the types carry JSON tags that match the public --json schema. Do not implement detection, probing, or rendering logic here.",
  "ownerAgent": "cli-engineer",
  "dependencies": ["FOUND-1"],
  "expectedOutputs": ["internal/model/model.go", "internal/model/model_test.go"],
  "validationCommands": ["go test ./internal/model/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/project-foundation.md#FOUND-FR-02"],
    "acceptanceCriteria": ["Model round-trips through JSON with stable field names", "Status parsing rejects values outside ok|missing|outdated|unknown", "A table-driven test covers every status value"],
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
| Unit Tests | version injection, domain model JSON | `go test ./internal/version/... ./internal/model/...` |
| Build | whole module | `go build ./...` and `make build` |

Key test scenarios:
1. Version defaults to `dev` and is replaced by an injected value.
2. Every status value serializes and parses, and an invalid status is rejected.

---

## 7. Acceptance Criteria

1. `go build ./...`, `make build`, and `make test` succeed from the repository root.
2. The shared model compiles with no import cycles and is the single source of core types.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Whether to add a golangci-lint config in this feature | Deferred to Release and Distribution CI |
