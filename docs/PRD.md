# PRD: repo-ready

## 1. Overview

**Product Name:** repo-ready
**Summary:** A single-binary terminal application that inspects a local or remote
git repository, derives the applications and tool versions the repo expects, diffs
them against what is installed on the user's machine, and reports copy-pasteable
install commands and documentation links. It never installs anything.
**Target Platform:** Linux, macOS, and Windows developer workstations, distributed
as a Go single binary.
**Key Constraints:** Report-only, no credential handling, offline by default,
declarative detection only, no user config or per-repo manifest in v1.
**Historical Sources:** [docs/IDEA.md](IDEA.md) (source material only; not an execution source).

---

## 2. Version History

| Version | Date | Author | Changes |
|---------|------|--------|---------|
| 1.0 | 2026-09-22 | - | Initial PRD authored from docs/IDEA.md via forge-auto-build-prd |

---

## 3. Goals and Non-Goals

### 3.1 Goals

- Answer "I cloned this, now what do I install?" in seconds from declarative
  files the repository already carries.
- Show correctness-first results: the declared constraint and the installed
  version are always displayed together, never hidden behind a status.
- Keep the local, deterministic report authoritative; online enrichment is
  additive.
- Ship one static binary with no runtime dependencies and an embedded install
  metadata catalog.

### 3.2 Non-Goals

- Executing installs, repository scripts, or repository code.
- Managing credentials, tokens, or SSH keys.
- User configuration, per-repo manifests, or plugin systems in v1.
- Submodule traversal, version-manager enumeration, or Windows-out.
- Curated "unverified suggestions" for repos that declare nothing (deferred past v1).
- Homebrew and Scoop installers (deferred past v1).

---

## 4. Personas

| Persona | Description | Key Needs |
|---------|-------------|-----------|
| Newcomer developer | Just cloned or was handed an unfamiliar repo and needs to set up | A fast, accurate list of required tools with install commands |
| Monorepo contributor | Works in a repo with many components declaring different versions | Findings grouped by component with conflicts surfaced |
| CI / automation engineer | Wants to gate pipelines on environment readiness | Stable JSON output and a strict exit code |

---

## 5. Research Findings

- The repository has no prior product code or requirements documents; the only
  product source is [docs/IDEA.md](IDEA.md), sharpened via forge-grill-idea.
- Declarative files such as `.tool-versions`, `.nvmrc`, `mise.toml`,
  `package.json` engines, `rust-toolchain.toml`, `go.mod`, `.python-version`,
  `devcontainer.json`, `Dockerfile`, `Makefile`, and README prerequisites
  already encode the intended environment for real repositories.
- Version matching must be loose and transparent: a pin is treated as a minimum
  and mismatches warn rather than fail, so the tool never blocks a developer on
  a technically-satisfiable environment.

---

## 6. Technical Architecture

### 6.1 Technology Stack

- **Language:** Go 1.27.1 (module `github.com/mcfuzzysquirrel/repo-ready`).
- **TUI:** [Bubble Tea v1.3.10](https://github.com/charmbracelet/bubbletea) with
  [Lip Gloss v1.1.0](https://github.com/charmbracelet/lipgloss) and
  [Bubbles v1.0.0](https://github.com/charmbracelet/bubbles).
- **Version parsing:** [Masterminds/semver v3.5.0](https://github.com/Masterminds/semver).
- **Ignore matching:** `github.com/sabhiram/go-gitignore` (pinned pseudo-version),
  used only to honor `.gitignore`; the built-in denylist is the primary guard.
- **CLI parsing:** Go standard library `flag`.
- Versions above were checked against their module proxies on 2026-09-22.

### 6.2 Project Structure

```
cmd/repo-ready/main.go        # CLI entrypoint
internal/app/                 # orchestration wiring
internal/input/               # local path vs remote URL resolution, shallow clone
internal/scan/                # filesystem walk, gitignore/denylist, components
internal/detect/              # per-manifest detectors + registry
internal/model/               # shared domain types
internal/probe/               # swappable probe backend, PATH + --version
internal/match/               # constraint parsing, loose matching, reconciliation
internal/catalog/             # embedded curated tool/install metadata
internal/report/              # JSON encoder
internal/tui/                 # Bubble Tea model/view/update, filters, copy
internal/enrich/              # optional online enrichment
internal/version/             # build-injected version metadata
testdata/                     # detector fixtures
e2e/                          # offline end-to-end and opt-in network tests
```

### 6.3 Platform and Packaging Constraints

```forge-requirement
{"id":"CONST-03","kind":"constraint","text":"repo-ready ships as one static Go binary with no runtime dependencies, embedding its curated install-metadata catalog and versioning that catalog alongside the binary."}
```

```forge-requirement
{"id":"CONST-04","kind":"constraint","text":"Linux, macOS, and Windows are first-class v1 target platforms for both runtime behavior and build artifacts."}
```

```forge-requirement
{"id":"CONST-07","kind":"constraint","text":"v1 does not enumerate version managers (nvm, asdf, mise, pyenv); installed-tool discovery is exposed through a swappable probing backend so that enumeration can be added later without changing callers."}
```

### 6.4 Key APIs / Interfaces

- **Detector:** `Detect(componentDir string, files []string) ([]model.Requirement, error)`;
  one implementation per declarative file, registered in a registry.
- **Probe backend:** `Probe(ctx context.Context, tool string) (model.InstalledTool, error)`
  behind an interface so version-manager enumeration can be added later.
- **Enricher:** `Enrich(ctx context.Context, findings []model.Finding) ([]model.Enrichment, error)`,
  optional and additive.

---

## 7. Non-Functional Requirements

| ID | Requirement | Priority |
|----|-------------|----------|
| CONST-11 | Performance budget and hard limits | Must |
| CONST-13 | Temporary clone safety | Must |
| CONST-05 | No user configuration or per-repo manifest in v1 | Must |
| CONST-06 | No submodule traversal | Must |
| CONST-08 | Offline by default | Must |

```forge-requirement
{"id":"CONST-11","kind":"constraint","text":"A local scan and report for a typical repository completes in under 2 seconds; probe commands each have a 5 second timeout; probing uses at most 8 concurrent workers; remote clones use depth 1; directory scanning defaults to a maximum depth of 6."}
```

```forge-requirement
{"id":"CONST-13","kind":"constraint","text":"Remote repositories are shallow-cloned into a uniquely named temporary directory created with owner-only permissions (0700) and removed on every exit path including success, error, and signal; no scan ever modifies the input repository."}
```

```forge-requirement
{"id":"CONST-05","kind":"constraint","text":"v1 has no user configuration, config file, or per-repository manifest; detection uses built-in detectors only, and the architecture must allow configuration and manifests to be layered on later."}
```

```forge-requirement
{"id":"CONST-06","kind":"constraint","text":"v1 does not traverse git submodules; submodule directories are skipped and reported as skipped paths."}
```

```forge-requirement
{"id":"CONST-08","kind":"constraint","text":"repo-ready is offline by default and performs no network I/O unless the user passes --enrich; when enrichment is requested it is additive and degrades gracefully when a feed is unreachable, so the local deterministic report always remains valid."}
```

---

## 8. Security and Privacy

| ID | Requirement | Priority |
|----|-------------|----------|
| CONST-01 | Report-only behavior | Must |
| CONST-02 | Credential neutrality | Must |
| CONST-09 | No-shell command probing | Must |
| CONST-10 | No telemetry, minimal network payload | Must |

```forge-requirement
{"id":"CONST-01","kind":"constraint","text":"repo-ready never installs tools and never executes repository code or repository scripts; it only reports requirements, install commands, and links for the user to run themselves."}
```

```forge-requirement
{"id":"CONST-02","kind":"constraint","text":"repo-ready never reads, stores, or manages credentials; remote clone operations reuse the user's existing git configuration and credential helpers, and on failure the tool surfaces git's own error and exits non-zero without touching credentials."}
```

```forge-requirement
{"id":"CONST-09","kind":"constraint","text":"Command probing runs binaries directly without a shell, only using known safe version-query invocations for recognized tools, with tool names validated against the bundled catalog and never interpolated into a shell string."}
```

```forge-requirement
{"id":"CONST-10","kind":"constraint","text":"repo-ready emits no telemetry and stores nothing remotely; when online enrichment is explicitly enabled the only outbound data is the tool identifiers and documentation URLs required for that lookup."}
```

---

## 9. Accessibility

| ID | Requirement | Priority |
|----|-------------|----------|
| CONST-12 | Accessible TUI presentation | Must |

```forge-requirement
{"id":"CONST-12","kind":"constraint","text":"Status is never encoded by color alone; the TUI is fully keyboard-navigable, honors the NO_COLOR environment variable and a --no-color flag, and a --plain mode renders screen-reader-friendly text without box-drawing characters."}
```

---

## 10. System States / Lifecycle

1. **Input resolution** — classify the positional argument as a local directory or
   a remote git URL; acquire a local scan root.
2. **Scan** — walk the tree within depth, denylist, and `.gitignore` limits; group
   by component; record skipped paths.
3. **Detect** — run declarative detectors per component to produce requirements.
4. **Probe** — query the local machine for installed tools and versions.
5. **Match** — compute statuses, de-duplicate by tool, and flag conflicts.
6. **Report** — render the interactive TUI or the versioned JSON document; apply
   enrichment only when `--enrich` is set.
7. **Exit** — 0 on a successful scan, 1 with `--strict` when a finding is missing
   or outdated, 2 on a fatal input/clone/usage error.

---

## 11. Analytics / Success Metrics

| Metric | Target | Measurement Method |
|--------|--------|--------------------|
| Detector accuracy | Correct toolset on fixture manifests | Unit tests with golden fixtures |
| Time to answer | < 2s for a typical local repo | Performance test in e2e suite |
| Cross-platform coverage | Linux, macOS, Windows | CI matrix green |

---

## 12. Dependencies and Risks

### 12.1 Dependencies

- Git available on `PATH` for remote input.
- Network only for optional enrichment.
- Upstream module availability for the pinned Go dependencies.

### 12.2 Risks

| Risk | Impact | Mitigation |
|------|--------|------------|
| Fuzzy README parsing produces false positives | Wrong requirements | Conservative, explicit-only parsing; README findings marked lower confidence |
| `.gitignore` matcher diverges from git | Scanning wrong files | Built-in denylist is primary; matcher covered by fixture tests |
| Probing hangs on a misbehaving binary | Slow scans | Per-command 5s timeout and bounded concurrency |
| Enrichment outages | Missing latest-version hints | Additive only; graceful degradation, never blocks or changes exit codes |

---

## 13. Future Considerations

Version-manager enumeration (nvm/asdf/mise/pyenv), user config and per-repo
manifests, curated unverified suggestions for under-declared repos, submodule
traversal, Homebrew/Scoop packaging, and additional detectors.

---

## 14. Features

| # | Feature | File | Dependencies | Priority |
|---|---------|------|-------------|----------|
| 1 | Project Foundation | [docs/features/project-foundation.md](features/project-foundation.md) | None | Must |
| 2 | Scan and Detection | [docs/features/scan-and-detect.md](features/scan-and-detect.md) | Project Foundation | Must |
| 3 | Tool Catalog | [docs/features/tool-catalog.md](features/tool-catalog.md) | Project Foundation | Must |
| 4 | Probing and Version Matching | [docs/features/probing-and-matching.md](features/probing-and-matching.md) | Project Foundation | Must |
| 5 | Online Enrichment | [docs/features/online-enrichment.md](features/online-enrichment.md) | Project Foundation, Tool Catalog, Probing and Version Matching | Should |
| 6 | Reporting | [docs/features/reporting.md](features/reporting.md) | Project Foundation, Scan and Detection, Tool Catalog, Probing and Version Matching, Online Enrichment | Must |
| 7 | Release and Distribution | [docs/features/release-and-distribution.md](features/release-and-distribution.md) | Project Foundation, Reporting | Must |
| 8 | End-to-End Validation | [docs/features/e2e-validation.md](features/e2e-validation.md) | Reporting, Scan and Detection, Probing and Version Matching | Must |

### Feature Dependency Graph

```
Project Foundation
├── Scan and Detection
├── Tool Catalog
├── Probing and Version Matching
├── Online Enrichment (needs Tool Catalog + Probing and Version Matching)
├── Reporting (needs Scan and Detection + Tool Catalog + Probing and Version Matching + Online Enrichment)
│   └── Release and Distribution
└── End-to-End Validation (needs Reporting + Scan and Detection + Probing and Version Matching)
```

---

## 15. Glossary

| Term | Definition |
|------|------------|
| Component | A directory within a scanned repo that declares one or more requirements |
| Constraint | The version expression a repository declares for a tool |
| Detector | A unit that reads one declarative file kind and emits requirements |
| Enrichment | Optional additive metadata fetched online (latest versions, link checks) |
| Finding | A reconciled requirement plus the local probe result and status |

---

## 16. Open Questions

| # | Question | Default Assumption |
|---|----------|--------------------|
| 1 | Exact TUI colors, density, and pass/warn/fail legibility | Settled by the TUI legibility human-review gate after a spike |
| 2 | Whether `go-gitignore` should be replaced by a maintained matcher | Keep it pinned for v1; revisit if it breaks |
| 3 | Exact detector list growth beyond the v1 set | v1 ships exactly the listed detectors |
| 4 | Which enrichment sources to use for latest versions and link checks | A single configurable HTTPS source with graceful degradation |
| 5 | Whether and when submodule traversal is added | Out of scope for v1 |
| 6 | When curated unverified suggestions ship | Post-v1 |
| 7 | Catalog tool count beyond 25 | 25-40 entries in v1; grow later |
