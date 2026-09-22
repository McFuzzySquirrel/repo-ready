# repo-ready User Guide

> [!WARNING]
> **This guide describes the specified (planned) v1 behavior. `repo-ready` has
> not been implemented or released.** There is no binary to download and no
> command that runs today. Everything below is the design captured in
> [`docs/PRD.md`](PRD.md) and [`docs/features/`](features/); commands and flags
> will only work once the implementation lands. See
> [`docs/releases/unreleased.md`](releases/unreleased.md) for validation status.

## Overview

`repo-ready` is a single-binary terminal app that inspects a local or remote git
repository, works out which tools and versions the repository expects from the
declarative files it already carries, compares them against what is installed
on your machine, and shows copy-pasteable install commands and documentation
links. It **only reports** — it never installs anything and never executes
repository code. For the design rationale, see
[`docs/adr/`](adr/README.md); for exact internals, see the feature
specifications in [`docs/features/`](features/).

It answers one question: *"I cloned this — now what do I install?"*

## Who this guide is for

- **Newcomer developer** — just cloned or was handed an unfamiliar repository
  and needs a fast, accurate setup list.
- **Monorepo contributor** — needs findings grouped by component, with version
  conflicts surfaced.
- **CI / automation engineer** — needs stable JSON and a strict exit code to
  gate a pipeline.

Installing or hosting the tool itself is covered by the
[administrator guide](admin-guide.md).

## Install and first use (planned)

Prerequisites for a future install:

- `git` on `PATH` — required only for **remote** repository inputs.
- Network — required only for remote inputs and for `--enrich`; a local scan is
  fully offline.

The specified installation paths (not yet available) are:

| Method | Specified command |
|--------|-------------------|
| Prebuilt binary | Download the archive for your OS/architecture from the GitHub release, verify the checksum, and put the binary on `PATH`. |
| Install script | `curl -fsSL https://raw.githubusercontent.com/McFuzzySquirrel/repo-ready/main/scripts/install.sh \| sh` |
| Go install | `go install github.com/mcfuzzysquirrel/repo-ready/cmd/repo-ready@latest` |
| From source | `git clone https://github.com/McFuzzySquirrel/repo-ready.git && cd repo-ready && make build` |

The first successful invocation is planned to be:

```bash
repo-ready ./path/to/repo
```

## Core workflow (planned)

1. **Point at a repository.** Pass one positional argument: an existing local
   directory path or a remote git URL. The kind is classified automatically;
   no flag is needed.
2. **Let it scan.** Local directories are read in place and offline. Remote URLs
   are shallow-cloned (`--depth 1`) into a temporary directory that is always
   cleaned up.
3. **Read the report.** The TUI lists one row per finding with a text status,
   the tool, the declared constraint, and the installed version.
4. **Open a tool's detail.** Press `Enter` to see why the tool is required, the
   declaring file and component, the documentation link, and per-platform
   install commands.
5. **Copy and install.** Press `c` to copy the selected install command (via
   OSC52), then run it yourself. `repo-ready` never runs it for you.

## Recipes

### Check a freshly cloned local repository

```bash
repo-ready ./path/to/repo
```

Expected: an interactive TUI grouping findings by component and pairing each
declared constraint with the installed version.

### Check a remote repository without cloning it yourself

```bash
repo-ready https://github.com/example/project.git
```

Expected: same report; the temporary clone is removed before the command exits.

### Produce machine-readable output for a script or CI

```bash
repo-ready ./path/to/repo --json
```

Expected: one versioned JSON document on stdout. Without `--strict` this exits
`0` whenever the scan succeeds, even if tools are missing.

### Gate a pipeline on a complete environment

```bash
repo-ready ./path/to/repo --strict
```

Expected: exit `1` if any finding is `missing` or `outdated`, otherwise `0`.

### Read the report with a screen reader

```bash
repo-ready ./path/to/repo --plain
```

Expected: linear text, no color, no box-drawing characters.

## Command and feature reference (planned)

### Flags

| Flag | Specified behavior |
|------|--------------------|
| `--json` | Emit a single versioned JSON document instead of the TUI. |
| `--strict` | Exit `1` when a required tool is missing or outdated (CI gating). |
| `--depth N` | Maximum directory scan depth. Default `6`. |
| `--enrich` | Enable additive online enrichment (latest versions, link checks). |
| `--no-color` | Disable color output. |
| `--plain` | Screen-reader-friendly output without color or box drawing. |
| `--version` | Print the version and exit. |
| `--help` | Show usage. |

### Finding statuses

| Status | Meaning |
|--------|---------|
| `ok` | The installed version satisfies the declared constraint. |
| `outdated` | Installed, but below the required minimum or mismatched major/minor. |
| `missing` | The tool was not found on `PATH`. |
| `unknown` | The constraint or installed version could not be parsed. |

Status is never conveyed by color alone.

### Version matching

| Declared | Treated as |
|----------|------------|
| `18.2.0`, `>=18.2.0` | minimum version |
| `^18.2.0` | same major as `18.2.0` |
| `~18.2.0`, `~>18.2.0` | same minor as `18.2.0` |
| `18` | same major |

A mismatch is `outdated`, which warns — it never silently fails. See
[ADR-0005](adr/0005-loose-transparent-version-matching.md).

### Declarations detected (planned)

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

README detection is deliberately conservative and marks findings with lower
confidence. Tools that are not in the embedded catalog are still reported with
their declared constraint and source, but have no install metadata.

### TUI keys (planned)

| Key | Action |
|-----|--------|
| `↑` / `↓`, `j` / `k` | Navigate the findings list |
| `Enter` | Open/close the detail pane |
| `f` | Cycle filters (all, missing-or-outdated, category) |
| `c` | Copy the selected install command via OSC52 |
| `q` | Quit |

## Configuration

There is **no user configuration and no per-repo manifest in v1** by design
(CONST-05). The only runtime inputs are the positional argument and the flags
above. `--enrich` is the sole control that enables network access.

## Safety and data handling

- **Report only.** `repo-ready` never installs tools and never executes
  repository code or scripts.
- **Credential-neutral.** Remote clones reuse your existing git configuration
  and credential helpers. The tool never reads, stores, or manages credentials;
  on clone failure it surfaces git's own error and exits non-zero.
- **Shell-free probing.** Recognized tools are run directly (no shell) with a
  5-second timeout, using only known safe version queries.
- **No telemetry.** Nothing is stored remotely. With `--enrich`, the only
  outbound data is tool identifiers and documentation URLs.
- **Non-destructive.** A local scan never modifies the repository; a remote scan
  cleans up its temporary clone on every exit path.

See [ADR-0002](adr/0002-report-only-never-install.md) and
[ADR-0009](adr/0009-bounded-scan-and-safe-clone.md).

## Troubleshooting (planned behavior)

| Symptom | Likely cause | Next action |
|---------|--------------|-------------|
| `repo-ready: command not found` | Not installed yet (no release exists) | Track the [release notes](releases/unreleased.md); build from source once implemented. |
| Scan reports fewer tools than expected | Depth limit, denylist, or `.gitignore` excluded the declaring file | Check the skipped-paths section of the report; raise `--depth` if appropriate. |
| A tool shows `unknown` | Constraint or installed version could not be parsed | Inspect the declared and installed strings shown on the finding. |
| Remote scan fails | `git` missing or clone rejected | Run the same `git clone` manually; the tool surfaces git's stderr. |
| `--enrich` notes "not available" | Feed unreachable or offline | None needed; enrichment is additive and the local report remains valid. |

## Further help

- [Product requirements](PRD.md) and [feature specifications](features/)
- [Architecture decisions](adr/README.md)
- [Changelog](../CHANGELOG.md) and [release notes](releases/README.md)
- [Administrator guide](admin-guide.md) — building, releasing, and operating
- Issue tracker: <https://github.com/McFuzzySquirrel/repo-ready/issues>
