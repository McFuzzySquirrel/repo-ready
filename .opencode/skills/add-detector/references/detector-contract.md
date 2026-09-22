# Detector contract reference

Authoritative sources: `docs/PRD.md` §6.4, `docs/features/scan-and-detect.md`
SCAN-FR-07..10, SCAN-3..6, and `internal/model` (owned by cli-engineer).

## Interface

```go
type Detector interface {
    // Handles reports the exact file base names this detector reads.
    Handles() []string
    // Detect reads the candidate files in componentDir and returns requirements.
    Detect(componentDir string, files []string) ([]model.Requirement, error)
}
```

The caller passes only files present in the component directory. A detector must
not read sibling directories, the network, or repository scripts.

## Registry behavior

- One registration per detector in `internal/detect/registry.go`.
- The registry runs every registered detector once per component.
- A file name claimed by a detector is offered only to that detector.
- The registry returns a deterministically ordered slice. Sort requirements by
  source file, then tool, then constraint before returning them.
- Identical inputs must yield byte-identical registry output (SCAN-FR-10).

## model.Requirement fields

| Field | Meaning | Rule |
|-------|---------|------|
| `tool` | Tool identifier | Lowercase, matches the catalog identifier where one exists |
| `constraint` | Declared version string | Copied verbatim from the file; empty when none is declared |
| `sourceFile` | Declaring file | Repository-relative path |
| `component` | Component directory | Directory that declares the requirement |
| `category` | Tool category | One of the catalog taxonomy values |
| `confidence` | Extraction confidence | High for structured files, reduced for README |

## Category taxonomy

Use the same vocabulary the catalog uses so grouping does not fragment: language
runtime, package manager, container, orchestration, cloud CLI, database client,
system dependency. Pick the single best fit per requirement.

## Confidence rules

- Structured manifests and dotfiles: high confidence.
- `Dockerfile`, `Makefile`: high only for explicit version declarations.
- README prerequisites: reduced confidence, always, even for fenced blocks
  (SCAN-FR-09).

## Error semantics

| Input | Result |
|-------|--------|
| Well-formed file | Requirements, no error |
| Tool named without a version | Requirement with empty constraint, no error |
| Malformed or unreadable file | Explicit error, zero requirements for that file |
| File kind not handled | Detector is not invoked |

Never return a guessed version and never downgrade a parse error to an empty
success. Wrap the underlying parse error with the source file name so the caller
can surface it.
