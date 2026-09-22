---
name: add-detector
description: "Add one declarative-file detector to repo-ready internal/detect: implement the Detector interface for a file kind, register it in the deterministic registry, extract declared constraints verbatim into model.Requirement, assign category and confidence, return structured errors for malformed input, and ship a table-driven fixture test under testdata/detect/."
---

# Skill: Add a Declarative-File Detector

Add one detector to `internal/detect` that reads a single declarative file kind and
emits `model.Requirement` values for repo-ready. Use this when the registry must
cover another file kind, an existing detector family gains a format, or a fixture
proves an extraction gap. Detectors are pure functions over file bytes; they never
resolve constraints, probe tools, or execute repository code.

## Process

### Step 1: Confirm the file kind is in scope

Compare the requested kind against the eleven v1 detectors: `.tool-versions`,
`.nvmrc`, `mise.toml`, `package.json` engines, `rust-toolchain.toml`, `go.mod`,
`.python-version`, `devcontainer.json`, `Dockerfile`, `Makefile`, and README
prerequisites (SCAN-FR-08). If the kind is not one of these, then stop and record
it as future work; v1 ships exactly this set and extras break the fixed detector
count.

### Step 2: Pick the detector family and declare file names

Choose the family file the detector belongs to by file type: dotfiles and
version-manager files in `internal/detect/dotfiles.go`, ecosystem manifests in
`internal/detect/manifests.go`, and text-derived files in
`internal/detect/textfiles.go`. Implement the `Detector` interface and return the
exact file names it handles. If two detectors claim the same file name, then merge
them or make the claim more specific; a claimed file must invoke exactly one
detector.

### Step 3: Extract declared constraints verbatim

Start from the raw declared string. Populate `tool`, `constraint`, `sourceFile`,
`component`, `category`, and `confidence` on each `model.Requirement`. Normalize
only enough to identify the tool. If the file names a tool without a version, then
emit the requirement with no constraint rather than inventing one. Never resolve a
constraint and never read a version from anywhere but the file (SCAN-FR-07).

### Step 4: Return structured errors instead of guessing

For malformed or unreadable content, return an explicit error and zero
requirements for that file. Do not treat malformed input as an empty success. Keep
README detection conservative: only fenced install blocks or explicit
`tool`-plus-version list lines count, with reduced confidence and a README source
(SCAN-FR-09).

### Step 5: Register deterministically and add the fixture

Register the detector so the registry runs every detector once per component.
Sort the requirements a detector returns so identical inputs produce byte-identical
output (SCAN-FR-10). Add a fixture directory under `testdata/detect/<kind>/` and a
table-driven test matching the family test name (`TestDotfileDetectors`,
`TestManifestDetectors`, or `TestTextDetectors`). Load `references/detector-contract.md` when you need the exact interface, model fields, category taxonomy, or error semantics.

### Step 6: Run the quality gate

Run the gate in `run-quality-gates` by default for every detector change, and fix
any formatting, vet, or test failure before committing. Cite `SCAN-FR-07` through
`SCAN-FR-10` in the commit message.

## Gotchas

- **Never invent a version.** A tool named without a version yields a requirement with an empty constraint, never a guessed one (SCAN-FR-07).
- **Multi-tool `.tool-versions` yields one requirement per line.** Do not collapse a file into a single requirement.
- **Malformed is not empty.** A syntactically invalid manifest is an error with no requirement, not a zero-requirement success.
- **Registry order must be stable.** Sort within the detector; downstream byte-identical reports depend on it (SCAN-FR-10).
- **Category and confidence are required fields.** Leave them unset and the report and catalog stages cannot group or label the finding.
- **Detectors stay pure.** No subprocess, no shell, no repository scripts; CONST-01 makes any execution a defect.
- **README is the false-positive risk.** General prose that mentions a tool is not a prerequisite declaration.

## Validation

Run this self-check before committing:

- [ ] `gofmt -l .` prints nothing and `go vet ./internal/detect/...` is clean
- [ ] `go test ./internal/detect/...` passes, including the new fixture test
- [ ] The new detector is registered and its claimed file invokes it exactly once
- [ ] A malformed fixture returns an error and no requirement
- [ ] Repeated runs of the registry produce identical ordering
- [ ] The fixture lives under `testdata/detect/<kind>/` and is committed, not generated

The reference `references/fixture-shape.md` holds the fixture layout and table
shape when you need to pin one.
