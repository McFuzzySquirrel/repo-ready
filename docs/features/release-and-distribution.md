# Feature: Release and Distribution

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| RELEASE-FR-01 | This feature | owns |
| RELEASE-FR-02 | This feature | owns |
| RELEASE-FR-03 | This feature | owns |
| RELEASE-FR-04 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Release and Distribution
**ID Prefix:** RELEASE
**Summary:** Continuous integration across all target platforms, a release
workflow that cross-compiles versioned archives with checksums, and
documentation that lets users install via `go install`, a script, or a prebuilt
binary.
**Dependencies:** Project Foundation, Reporting
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| RELEASE-US-01 | user | to install from a prebuilt binary, a script, or go install | I can adopt the tool however I prefer | Must |
| RELEASE-US-02 | maintainer | tagged releases to build all platforms automatically | releases are reproducible and consistent | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"RELEASE-FR-01","kind":"requirement","text":"A GitHub Actions CI workflow runs formatting checks, go vet, and the Go test suite on Linux, macOS, and Windows for every push and pull request."}
```

```forge-requirement
{"id":"RELEASE-FR-02","kind":"requirement","text":"A release workflow triggered by a version tag cross-compiles Linux, macOS, and Windows binaries for amd64 and arm64, injects the version via -ldflags, produces per-platform archives plus checksums, and attaches them to a GitHub Release."}
```

```forge-requirement
{"id":"RELEASE-FR-03","kind":"requirement","text":"The project supports go install github.com/mcfuzzysquirrel/repo-ready/cmd/repo-ready@latest and provides a documented install script for prebuilt binaries."}
```

```forge-requirement
{"id":"RELEASE-FR-04","kind":"requirement","text":"The README documents installation, usage, every flag, the exit-code contract, and links the JSON schema documentation."}
```

---

## 4. UI / Interaction Design

None; this feature covers CI, packaging, and documentation.

---

## 5. Implementation Tasks

### Phase 1: CI, Release, and Documentation

```forge-task
{
  "id": "RELEASE-1",
  "title": "Add the cross-platform CI workflow",
  "description": "Add a GitHub Actions workflow that runs on push and pull request with a matrix over ubuntu, macos, and windows, setting up the Go toolchain, verifying gofmt with no diffs, running go vet, and running go test ./... . Fail the job on any formatting diff, vet finding, or failing test. Do not add release publishing here.",
  "ownerAgent": "release-engineer",
  "dependencies": [],
  "expectedOutputs": [".github/workflows/ci.yml"],
  "validationCommands": ["go vet ./...", "go test ./..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/release-and-distribution.md#RELEASE-FR-01"],
    "acceptanceCriteria": ["The workflow matrix includes Linux, macOS, and Windows", "Formatting, vet, and test steps fail the job on failure", "The workflow file is valid YAML"],
    "constraints": ["CI must not require network beyond module downloads"],
    "constraintRefs": ["docs/PRD.md#CONST-08"],
    "references": []
  }
}
```

```forge-task
{
  "id": "RELEASE-2",
  "title": "Add the tagged release build and checksum workflow",
  "description": "Add GoReleaser configuration and a GitHub Actions release workflow triggered by a version tag. Cross-compile linux, darwin, and windows for amd64 and arm64, inject the version through -ldflags into internal/version, produce archives and a checksums file, and publish them to the GitHub Release for the tag. Do not publish Homebrew or Scoop packages in v1.",
  "ownerAgent": "release-engineer",
  "dependencies": ["RELEASE-1", "REPORT-5"],
  "expectedOutputs": [".goreleaser.yaml", ".github/workflows/release.yml"],
  "validationCommands": ["go build -ldflags \"-X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=test\" ./..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/release-and-distribution.md#RELEASE-FR-02"],
    "acceptanceCriteria": ["The build matrix covers all three operating systems and both architectures", "The version is injected via -ldflags", "Archives and a checksums file are produced and attached to the release"],
    "constraints": ["Produce static binaries with no runtime dependencies and no Homebrew/Scoop packaging in v1"],
    "constraintRefs": ["docs/PRD.md#CONST-03"],
    "references": []
  }
}
```

```forge-task
{
  "id": "RELEASE-3",
  "title": "Write install script and user documentation",
  "description": "Add a POSIX install script that detects OS and architecture, downloads the matching release archive, verifies its checksum, and installs the binary to a user-writable location. Write the README covering installation via prebuilt binary, the install script, and go install, plus usage, every flag, the 0/1/2 exit-code contract, and a link to the JSON schema document. Add docs/JSON-SCHEMA.md describing the versioned --json fields.",
  "ownerAgent": "release-engineer",
  "dependencies": ["RELEASE-2"],
  "expectedOutputs": ["scripts/install.sh", "README.md", "docs/JSON-SCHEMA.md"],
  "validationCommands": ["bash -n scripts/install.sh"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/release-and-distribution.md#RELEASE-FR-03", "docs/features/release-and-distribution.md#RELEASE-FR-04"],
    "acceptanceCriteria": ["The install script passes shell syntax checking and verifies checksums", "The README documents all install paths, flags, and exit codes", "The JSON schema document lists every emitted field"],
    "constraints": ["Installing is performed by the user; the tool itself never installs anything"],
    "constraintRefs": ["docs/PRD.md#CONST-01"],
    "references": []
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| CI | fmt, vet, tests on three platforms | GitHub Actions matrix |
| Build | cross-compilation and version injection | Local `go build` with ldflags and release dry run |
| Docs | script syntax | `bash -n` and manual review of documented flags |

Key test scenarios:
1. A tag produces archives for all six platform/architecture combinations with a checksums file.
2. `go build` with the version ldflag embeds the requested version.

---

## 7. Acceptance Criteria

1. CI is green across Linux, macOS, and Windows.
2. Tagged releases publish versioned archives and checksums for all target platforms.
3. Installation and flag documentation is complete.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Whether GoReleaser or a plain Actions matrix is used | GoReleaser with a GitHub Actions workflow |
| 2 | Signing or notarization | Out of scope for v1 |
