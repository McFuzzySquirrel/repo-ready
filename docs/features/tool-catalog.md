# Feature: Tool Catalog

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| CATALOG-FR-01 | This feature | owns |
| CATALOG-FR-02 | This feature | owns |
| CATALOG-FR-03 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Tool Catalog
**ID Prefix:** CATALOG
**Summary:** Bundle a curated, versioned catalog of tool metadata — display name,
category, documentation URL, and per-platform install commands — embedded in the
binary and looked up by tool identifier.
**Dependencies:** Project Foundation
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| CATALOG-US-01 | newcomer developer | an install command and docs link for each required tool | I can set up without searching the web | Must |
| CATALOG-US-02 | maintainer | catalog data validated at load time | bad entries fail loudly instead of printing wrong commands | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"CATALOG-FR-01","kind":"requirement","text":"Bundle in the binary a curated catalog of 25 to 40 tools, each with a stable identifier, display name, category, https documentation URL, and explicit Linux, macOS, and Windows install commands, versioned alongside the binary."}
```

```forge-requirement
{"id":"CATALOG-FR-02","kind":"requirement","text":"Load and validate the catalog at startup, requiring unique identifiers, non-empty required fields, and valid https documentation URLs, and fail with a clear error if the bundled catalog is invalid."}
```

```forge-requirement
{"id":"CATALOG-FR-03","kind":"requirement","text":"A tool present in a repository but absent from the catalog is still reported with its declared constraint and source, using no catalog metadata, and never causes a hard failure."}
```

---

## 4. UI / Interaction Design

Catalog metadata surfaces in the Reporting detail pane. No catalog-specific UI.

---

## 5. Implementation Tasks

### Phase 1: Catalog Schema and Content

```forge-task
{
  "id": "CATALOG-1",
  "title": "Implement the catalog schema, loader, and validation",
  "description": "Implement internal/catalog with a ToolEntry type (id, name, category, docs URL, and linux/macos/windows install commands), an embedded JSON catalog via go:embed, a loader, and validation requiring unique ids, non-empty name/category/docs URL, an https docs URL, and at least one non-empty install command per platform. Provide a lookup that returns an entry or a not-found result so callers can report unknown tools without failing. Include a small seed catalog so the package compiles before the full catalog is authored.",
  "ownerAgent": "catalog-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/catalog/catalog.go", "internal/catalog/catalog_test.go", "internal/catalog/data/tools.json", "internal/catalog/testdata/invalid_catalog.json"],
  "validationCommands": ["go test ./internal/catalog/...", "go vet ./internal/catalog/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/tool-catalog.md#CATALOG-FR-01", "docs/features/tool-catalog.md#CATALOG-FR-02", "docs/features/tool-catalog.md#CATALOG-FR-03"],
    "acceptanceCriteria": ["Duplicate ids, empty fields, and non-https docs URLs fail validation", "Lookup of an unknown tool returns not-found without an error", "The embedded catalog loads without network access"],
    "constraints": ["Catalog data is embedded and versioned with the binary"],
    "constraintRefs": ["docs/PRD.md#CONST-03"],
    "references": []
  }
}
```

```forge-task
{
  "id": "CATALOG-2",
  "title": "Author the curated 25-40 tool catalog",
  "description": "Populate internal/catalog/data/tools.json with 25 to 40 curated tools spanning language runtimes, package managers and dev tooling, containers and orchestration, cloud CLIs, database clients, and system-level dependencies. Each entry carries a stable id, display name, category, a canonical https docs URL, and install commands for Linux, macOS, and Windows. Add a coverage test that asserts every entry validates, the category set matches the intended spread, and each platform command is non-empty for every tool.",
  "ownerAgent": "catalog-engineer",
  "dependencies": ["CATALOG-1"],
  "expectedOutputs": ["internal/catalog/data/tools.json", "internal/catalog/tools_coverage_test.go"],
  "validationCommands": ["go test ./internal/catalog/... -run TestCatalogCoverage"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/tool-catalog.md#CATALOG-FR-01"],
    "acceptanceCriteria": ["The catalog contains between 25 and 40 validated entries", "Every entry has Linux, macOS, and Windows install commands", "Categories span runtimes, package managers, containers, cloud CLIs, database clients, and system dependencies"],
    "constraints": ["Install commands must be official, copy-pasteable, and non-destructive"],
    "constraintRefs": ["docs/PRD.md#CONST-01"],
    "references": []
  }
}
```

```forge-task
{
  "id": "CATALOG-3",
  "title": "Human review of catalog install-command and link accuracy",
  "description": "A human reviewer spot-checks the curated catalog across platforms. For a sample spanning every category, verify that each docs URL resolves to the official documentation and that the Linux, macOS, and Windows install commands are correct, current, and non-destructive. Record the sampled tool ids, the observed results, any corrections required, and the reviewer decision in the review file. This gate cannot be satisfied by an agent.",
  "dependencies": ["CATALOG-2"],
  "expectedOutputs": [],
  "validationCommands": [],
  "contract": {
    "version": 2,
    "kind": "human-review",
    "requirements": [],
    "requirementRefs": ["docs/features/tool-catalog.md#CATALOG-FR-01"],
    "acceptanceCriteria": ["A per-category sample is recorded with observed URL and command results", "Every flagged correction is either applied or explicitly accepted with a rationale", "The reviewer records an explicit approval decision"],
    "constraints": [],
    "constraintRefs": [],
    "references": [],
    "reviewFile": "docs/reviews/catalog-accuracy.json"
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Unit Tests | schema validation and lookup | Table-driven tests with valid and invalid fixtures |
| Data Tests | full catalog contents | Coverage test over the embedded catalog |
| Human Review | link and command accuracy | Manual spot-check recorded as review evidence |

Key test scenarios:
1. An invalid catalog entry fails loading with a descriptive error.
2. An unknown tool id returns not-found and does not abort a report.
3. Every catalog entry has all three platform commands.

---

## 7. Acceptance Criteria

1. The binary embeds a validated catalog of 25-40 tools.
2. Unknown tools are reported without catalog metadata and without failure.
3. The catalog accuracy human-review gate is recorded and approved.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Exact 25-40 tool selection | Choose high-frequency tools per category during authoring |
| 2 | Whether to include a JSON-LD or schema version field | Include a catalog version string |
