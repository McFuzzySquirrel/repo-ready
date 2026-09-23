# repo-ready

**Point it at any repository and find out exactly what you need to install.**

![Platforms](https://img.shields.io/badge/platforms-linux%20%7C%20macOS%20%7C%20windows-blue)
![Go](https://img.shields.io/badge/go-1.27-00ADD8?logo=go&logoColor=white)
![Report only](https://img.shields.io/badge/behavior-report--only-success)

> [!IMPORTANT]
> **Specification stage — CLI not implemented or released.** The Go module,
> package skeleton, and `Makefile` now exist (see [Development](#development)),
> but the scanning behavior described below is still under construction: running
> the binary only reports its version, and no release has shipped. Everything
> else in this README describes the **specified v1 design**. See
> [CHANGELOG.md](CHANGELOG.md) and [docs/releases/](docs/releases/README.md) for
> status.

`repo-ready` is a single-binary terminal app that inspects a local or remote git
repository, works out which tools and versions the repo expects, compares them
against what is already installed on your machine, and hands you copy-pasteable
install commands and documentation links.

It answers one question in seconds: **"I cloned this — now what do I install?"**
No README reading, no guesswork.

> [!NOTE]
> `repo-ready` only **reports**. It never installs anything and never executes
> repository code or scripts.

The v1 specification — the authoritative design for the planned implementation —
lives in [docs/PRD.md](docs/PRD.md) and [docs/features/](docs/features/). Design
rationale is recorded in [docs/adr/](docs/adr/README.md).

## Why

A repository already encodes its intended environment in the declarative files
it carries — `.nvmrc`, `.tool-versions`, `go.mod`, `package.json` engines, and
friends. `repo-ready` uses those as the authoritative source instead of guessing
from prose, then diffs them against the tools actually installed on your machine.

- **Local or remote** — scan a directory in place, or shallow-clone a git URL
  into a temporary directory that is always cleaned up.
- **Correctness first** — the declared constraint and the installed version are
  always shown together, never hidden behind a status.
- **Monorepo aware** — findings are grouped by component, and conflicting
  declarations are surfaced instead of silently resolved.
- **Transparent matching** — version pins are treated as minimums and mismatches
  **warn** rather than fail.
- **Report only** — it prints install commands and links for you to run. It
  never touches credentials and never modifies the repository.
- **One static binary** — the curated install-metadata catalog is embedded and
  versioned alongside the binary. No runtime dependencies.

## How it works

```
input resolution → scan → detect → probe → match → report
```

1. **Input** — classify the argument as a local directory or a remote git URL.
2. **Scan** — walk the tree within depth, denylist, and `.gitignore` limits, and
   group candidate files by component. Skipped paths are recorded.
3. **Detect** — run declarative detectors to extract declared requirements.
4. **Probe** — look up each tool on `PATH` and query its version.
5. **Match** — compute pass/warn statuses, de-duplicate by tool, and flag
   conflicts.
6. **Report** — render the interactive TUI, or emit versioned JSON with
   `--json`.

## Installation

### Prebuilt binary

Download the archive for your OS and architecture from the
[latest release](https://github.com/McFuzzySquirrel/repo-ready/releases), verify
the checksum, and put the binary on your `PATH`.

### Install script

```bash
curl -fsSL https://raw.githubusercontent.com/McFuzzySquirrel/repo-ready/main/scripts/install.sh | sh
```

The script detects your platform, downloads the matching release, and verifies
its checksum before installing.

### Go install

```bash
go install github.com/mcfuzzysquirrel/repo-ready/cmd/repo-ready@latest
```

### From source

```bash
git clone https://github.com/McFuzzySquirrel/repo-ready.git
cd repo-ready
make build
```

## Usage

```bash
# Scan a local repository
repo-ready ./path/to/repo

# Scan a remote repository (shallow-cloned, then cleaned up)
repo-ready https://github.com/example/project.git

# Machine-readable output for scripting and CI
repo-ready ./path/to/repo --json

# Gate a pipeline on a complete environment
repo-ready ./path/to/repo --strict
```

### Flags

| Flag | Description |
|------|-------------|
| `--json` | Emit a single versioned JSON document instead of the TUI. |
| `--strict` | Exit `1` when a required tool is missing or outdated (CI gating). |
| `--depth N` | Maximum directory scan depth. Default `6`. |
| `--enrich` | Enable additive online enrichment (latest versions, link checks). |
| `--no-color` | Disable color output. |
| `--plain` | Screen-reader-friendly output without color or box drawing. |
| `--version` | Print the version and exit. |
| `--help` | Show usage. |

> [!IMPORTANT]
> No network I/O happens by default. Remote input needs `git`, and online
> enrichment only runs when you pass `--enrich`.

## Reading the report

Each finding pairs a text status with the tool, the declared constraint, and the
installed version.

| Status | Meaning |
|--------|---------|
| `ok` | The installed version satisfies the declared constraint. |
| `outdated` | Installed, but below the required minimum or major/minor. |
| `missing` | The tool was not found on `PATH`. |
| `unknown` | The constraint or installed version could not be parsed. |

Status is never conveyed by color alone. The TUI is fully keyboard-driven
(navigate with arrows or `j`/`k`, `Enter` for details, `f` to cycle filters, `c`
to copy an install command).

In the TUI, the detail pane explains **why** a tool is required, which file
declared it, its category, the documentation link, and per-platform install
commands.

### Version matching

| Declared | Treated as |
|----------|------------|
| `18.2.0`, `>=18.2.0` | minimum version |
| `^18.2.0` | same major as `18.2.0` |
| `~18.2.0`, `~>18.2.0` | same minor as `18.2.0` |
| `18` | same major |

A tool is `ok` when the installed version is at or above a minimum or matches
the allowed major/minor. Anything else is `outdated`, which warns — it never
silently fails.

## Supported declarations

v1 detects requirements from the declarative files a repository already carries:

| File | What is read |
|------|--------------|
| `.tool-versions` | asdf/mise tool pins |
| `.nvmrc` | Node.js version |
| `mise.toml` | mise tool versions |
| `package.json` | `engines` |
| `rust-toolchain.toml` | Rust toolchain |
| `go.mod` | Go version directive |
| `.python-version` | Python version |
| `devcontainer.json` | Dev container features |
| `Dockerfile` | `FROM` image tags |
| `Makefile` | Explicit tool/version declarations |
| `README` | Explicit prerequisite blocks only |

> [!NOTE]
> README detection is deliberately conservative: it only recognizes explicit
> prerequisite declarations, marks them with lower confidence, and never derives
> requirements from general prose.

Tools that appear in a repository but are not in the embedded catalog are still
reported with their declared constraint and source — they simply have no install
metadata.

## JSON output

`--json` emits a single versioned object covering input metadata, scan metadata
(including skipped paths), components, findings with declared and installed
values, conflict declarations, and optional enrichment. Field names are planned
to be stable across releases. A `docs/JSON-SCHEMA.md` reference document will be
added with the implementation; the schema and compatibility contract are
specified in [ADR-0007](docs/adr/0007-versioned-json-and-exit-codes.md).

## Exit codes

| Code | Meaning |
|------|---------|
| `0` | The scan completed successfully (even with missing tools). |
| `1` | `--strict` was set and at least one finding is missing or outdated. |
| `2` | A fatal input, clone, or usage error. |

`--json` without `--strict` exits `0` whenever the scan succeeds — the data is
the output, and the caller decides what to do with it.

## Development

> [!NOTE]
> The Go module skeleton and `Makefile` now exist, so these four targets work
> today. The scanning behavior they build is still being implemented; see the
> status note at the top of this file.

Requires Go 1.27 or later. `git` is needed only for remote inputs.

```bash
make build   # build the binary
make test    # run the test suite
make vet     # run go vet
make fmt     # format the code
```

The test suite is offline and deterministic by default. An opt-in corpus test
runs against real remote repositories only when its environment variable is set,
and a performance test guards the 2-second budget for a typical local scan.

## Documentation

- [User guide](docs/user-guide.md)
- [Administrator guide](docs/admin-guide.md)
- [Architecture decisions](docs/adr/README.md)
- [Changelog](CHANGELOG.md) and [release notes](docs/releases/README.md)
- [Product requirements](docs/PRD.md) and [feature specifications](docs/features/)

## Security and privacy

- Never installs tools and never runs repository code.
- Never reads, stores, or manages credentials; remote clones reuse your existing
  git configuration and credential helpers.
- Probes recognized tools directly without a shell, with a 5-second timeout.
- No telemetry. When enrichment is explicitly enabled, the only outbound data is
  tool identifiers and documentation URLs.
