# Feature: Scan and Detection

## Traceability

| Canonical ID | Owner / Source Link | Relationship |
|--------------|---------------------|--------------|
| SCAN-FR-01 | This feature | owns |
| SCAN-FR-02 | This feature | owns |
| SCAN-FR-03 | This feature | owns |
| SCAN-FR-04 | This feature | owns |
| SCAN-FR-05 | This feature | owns |
| SCAN-FR-06 | This feature | owns |
| SCAN-FR-07 | This feature | owns |
| SCAN-FR-08 | This feature | owns |
| SCAN-FR-09 | This feature | owns |
| SCAN-FR-10 | This feature | owns |

**PRD:** [docs/PRD.md](../PRD.md)

---

## 1. Feature Overview

**Feature Name:** Scan and Detection
**ID Prefix:** SCAN
**Summary:** Acquire a local scan root from a path or a remote URL, walk it under
depth/denylist/gitignore limits grouping by component, and derive requirements
from the declarative files the repository carries.
**Dependencies:** Project Foundation
**Priority:** Must

---

## 2. User Stories

| ID | As a... | I want to... | So that... | Priority |
|----|---------|-------------|-----------|----------|
| SCAN-US-01 | newcomer developer | point the tool at a repo by path or URL | I get requirements without manual setup | Must |
| SCAN-US-02 | monorepo contributor | findings grouped by component | I know which package needs which version | Must |

---

## 3. Functional Requirements

```forge-requirement
{"id":"SCAN-FR-01","kind":"requirement","text":"Accept one positional input that is either an existing local directory path or a remote git URL, and classify its kind automatically without requiring an explicit flag."}
```

```forge-requirement
{"id":"SCAN-FR-02","kind":"requirement","text":"For a local directory input, scan it in place and fully offline without modifying any file in the directory."}
```

```forge-requirement
{"id":"SCAN-FR-03","kind":"requirement","text":"For a remote input, shallow-clone with depth 1 into a uniquely named temporary directory and remove that directory on every exit path including success, error, and signal."}
```

```forge-requirement
{"id":"SCAN-FR-04","kind":"requirement","text":"Remote clone operations reuse the user's existing git configuration and credential helpers, never prompt for or store credentials, and on failure surface git's own stderr and cause a non-zero exit."}
```

```forge-requirement
{"id":"SCAN-FR-05","kind":"requirement","text":"Recursively scan the input tree to a default maximum depth of 6, skipping .git, submodules, a built-in denylist (node_modules, vendor, .venv, dist, build, and similar), and paths matched by .gitignore."}
```

```forge-requirement
{"id":"SCAN-FR-06","kind":"requirement","text":"Group detected requirements by the component directory that declares them, and record every skipped path so the report can state that the scan was not exhaustive."}
```

```forge-requirement
{"id":"SCAN-FR-07","kind":"requirement","text":"Provide a detector registry and a detector interface in which each detector reads one declarative file kind and returns tool, constraint, source file, component, category, and confidence for each requirement, never inventing versions absent from the file."}
```

```forge-requirement
{"id":"SCAN-FR-08","kind":"requirement","text":"Ship exactly these v1 detectors: .tool-versions, .nvmrc, mise.toml, package.json engines, rust-toolchain.toml, go.mod, .python-version, devcontainer.json, Dockerfile, Makefile, and README prerequisites."}
```

```forge-requirement
{"id":"SCAN-FR-09","kind":"requirement","text":"README prerequisite detection is conservative: it recognizes only explicit prerequisite declarations such as fenced install blocks or lines naming a tool with a version, marks those findings with lower confidence and a README source, and never derives requirements from general prose."}
```

```forge-requirement
{"id":"SCAN-FR-10","kind":"requirement","text":"Detector and component results are returned in a deterministic, stable order so identical inputs produce byte-identical reports."}
```

---

## 4. UI / Interaction Design

No dedicated UI; this feature produces the requirement model consumed by Reporting.

---

## 5. Implementation Tasks

### Phase 1: Input Acquisition and Scanning

```forge-task
{
  "id": "SCAN-1",
  "title": "Resolve local paths and shallow-clone remote inputs",
  "description": "Implement input classification and acquisition in internal/input. Classify a single argument as an existing local directory or a remote git URL. For local directories return the absolute path as the scan root without modifying it. For remotes create a 0700 temporary directory, run git clone --depth 1 into it without a shell, and return the clone path plus a cleanup function that removes the directory on every exit path including error and signal contexts. On clone failure return the git stderr and a fatal error. Do not detect requirements or probe tools here.",
  "ownerAgent": "scan-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/input/resolve.go", "internal/input/clone.go", "internal/input/input_test.go"],
  "validationCommands": ["go test ./internal/input/...", "go vet ./internal/input/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-01", "docs/features/scan-and-detect.md#SCAN-FR-02", "docs/features/scan-and-detect.md#SCAN-FR-03", "docs/features/scan-and-detect.md#SCAN-FR-04"],
    "acceptanceCriteria": ["A local directory and a remote URL are each classified correctly", "The clone runs depth 1 with no shell and a 0700 temp dir", "Cleanup removes the temp dir on success and on error", "A clone failure returns git stderr and a fatal error"],
    "constraints": ["Never read or write credentials; rely on the ambient git configuration"],
    "constraintRefs": ["docs/PRD.md#CONST-02", "docs/PRD.md#CONST-09", "docs/PRD.md#CONST-13"],
    "references": []
  }
}
```

```forge-task
{
  "id": "SCAN-2",
  "title": "Walk the tree and group requirements by component",
  "description": "Implement internal/scan. Walk a scan root to a configurable maximum depth (default 6), skipping .git, submodule directories, a built-in denylist, and paths matched by the repository's .gitignore files. Collect candidate declarative files and group them by their containing component directory. Return the sorted component list and an explicit list of skipped paths with reasons. Do not parse file contents here.",
  "ownerAgent": "scan-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/scan/scanner.go", "internal/scan/gitignore.go", "internal/scan/denylist.go", "internal/scan/scanner_test.go", "testdata/scan/.gitignore", "testdata/scan/nested/app/.tool-versions"],
  "validationCommands": ["go test ./internal/scan/...", "go vet ./internal/scan/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-05", "docs/features/scan-and-detect.md#SCAN-FR-06", "docs/features/scan-and-detect.md#SCAN-FR-10"],
    "acceptanceCriteria": ["Depth, denylist, and .gitignore exclusions are each covered by a fixture test", "Skipped paths are reported with a reason", "Results are deterministically ordered across repeated runs"],
    "constraints": ["Never traverse submodules or mutate the scanned tree"],
    "constraintRefs": ["docs/PRD.md#CONST-06", "docs/PRD.md#CONST-11"],
    "references": []
  }
}
```

### Phase 2: Declarative Detectors

```forge-task
{
  "id": "SCAN-3",
  "title": "Build the detector registry and interface",
  "description": "Implement internal/detect with a Detector interface that maps a component directory and its candidate files to model.Requirement values, and a registry that runs every registered detector and returns the combined, deterministically ordered requirements. A detector declares which file names it handles. Do not implement individual manifest parsers here.",
  "ownerAgent": "detector-engineer",
  "dependencies": [],
  "expectedOutputs": ["internal/detect/detector.go", "internal/detect/registry.go", "internal/detect/registry_test.go"],
  "validationCommands": ["go test ./internal/detect/...", "go vet ./internal/detect/..."],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-07"],
    "acceptanceCriteria": ["The registry runs every registered detector exactly once per component", "A detector claiming a file name is invoked for that file", "Registry output order is stable", "A detector never emits a version not present in its input"],
    "constraints": ["Detectors are pure functions over files and must not execute repository code"],
    "constraintRefs": ["docs/PRD.md#CONST-01"],
    "references": []
  }
}
```

```forge-task
{
  "id": "SCAN-4",
  "title": "Implement dotfile and version-manager detectors",
  "description": "Implement detectors for .tool-versions, .nvmrc, mise.toml, .python-version, and rust-toolchain.toml. Each detector reads its file, extracts tool names and declared constraints exactly as written, and emits a model.Requirement with the declaring file, component, category, and confidence. Cover plain pins, major-only pins, and multi-tool files, and return an explicit error for unreadable content rather than guessing.",
  "ownerAgent": "detector-engineer",
  "dependencies": ["SCAN-3"],
  "expectedOutputs": ["internal/detect/dotfiles.go", "internal/detect/dotfiles_test.go", "testdata/detect/toolversions/.tool-versions", "testdata/detect/nvmrc/.nvmrc"],
  "validationCommands": ["go test ./internal/detect/... -run TestDotfileDetectors"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-08"],
    "acceptanceCriteria": ["Each of the five dotfile formats is covered by a table-driven fixture test", "Multi-tool .tool-versions yields one requirement per tool", "Malformed input returns an error and no requirement"],
    "constraints": [],
    "constraintRefs": [],
    "references": []
  }
}
```

```forge-task
{
  "id": "SCAN-5",
  "title": "Implement ecosystem manifest detectors",
  "description": "Implement detectors for package.json engines, go.mod, and devcontainer.json. Parse the language's declared tool and version fields using the standard library or the pinned dependencies, emit model.Requirement values with the declaring file and component, and ignore fields that declare no tool version. Return a structured error for syntactically invalid manifests.",
  "ownerAgent": "detector-engineer",
  "dependencies": ["SCAN-3"],
  "expectedOutputs": ["internal/detect/manifests.go", "internal/detect/manifests_test.go", "testdata/detect/packagejson/package.json", "testdata/detect/gomod/go.mod"],
  "validationCommands": ["go test ./internal/detect/... -run TestManifestDetectors"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-08"],
    "acceptanceCriteria": ["package.json engines, go.mod directives, and devcontainer.json features are each parsed by a fixture test", "Invalid manifests produce an error and no requirement", "Declared constraints are preserved verbatim"],
    "constraints": [],
    "constraintRefs": [],
    "references": []
  }
}
```

```forge-task
{
  "id": "SCAN-6",
  "title": "Implement conservative text-file detectors",
  "description": "Implement detectors for Dockerfile, Makefile, and README prerequisites. For Dockerfile and Makefile, extract explicit tool/version declarations such as FROM image tags and documented tool versions. For README, recognize only explicit prerequisite blocks: fenced install command blocks referencing a known tool, or bullet/list lines of the form tool plus a version or minimum. Mark every README-derived requirement with reduced confidence and a README source, and never derive a requirement from general prose.",
  "ownerAgent": "detector-engineer",
  "dependencies": ["SCAN-3"],
  "expectedOutputs": ["internal/detect/textfiles.go", "internal/detect/textfiles_test.go", "testdata/detect/readme/README.md"],
  "validationCommands": ["go test ./internal/detect/... -run TestTextDetectors"],
  "contract": {
    "version": 2,
    "kind": "implementation",
    "requirements": [],
    "requirementRefs": ["docs/features/scan-and-detect.md#SCAN-FR-08", "docs/features/scan-and-detect.md#SCAN-FR-09"],
    "acceptanceCriteria": ["Explicit README prerequisite blocks yield requirements with reduced confidence", "Prose that merely mentions a tool yields no requirement", "Dockerfile and Makefile fixtures each yield the expected requirements"],
    "constraints": ["README detection must never guess from unstructured prose"],
    "constraintRefs": [],
    "references": []
  }
}
```

---

## 6. Testing Strategy

| Level | Scope | Approach |
|-------|-------|----------|
| Unit Tests | input resolution, cloning, scanning, each detector | Table-driven fixture tests under `testdata/` |
| Integration Tests | scan then detect on a fixture repo | Component grouping and ordering assertions |

Key test scenarios:
1. A local fixture repo with nested components is grouped correctly and reproducibly.
2. A malformed manifest produces an error instead of a guessed requirement.
3. A README that only mentions a tool in prose yields nothing.

---

## 7. Acceptance Criteria

1. All eleven v1 detectors are implemented and covered by fixture tests.
2. Remote inputs are shallow-cloned, scanned, and cleaned up on every path.
3. Scan output order is deterministic for identical inputs.

---

## 8. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Exact denylist contents beyond the listed names | Start with node_modules, vendor, .venv, dist, build, target, .git |
| 2 | Whether Makefile version extraction should be enabled broadly | Only explicit version declarations; otherwise skip |
