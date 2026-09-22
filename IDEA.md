# Project Idea

**repo-ready** — a TUI application that inspects a local or remote git repository
and tells you which applications and tool versions you need to work on it,
including links to where to get them and copy-pasteable instructions.

The core promise: point it at any repo, and it answers "I cloned this, now what
do I install?" in seconds, without reading the README.

---

## Problem and Users

- **Primary user:** a developer who just encountered an unfamiliar repo (cloned
  it, been handed it, or is onboarding) and needs a fast, accurate list of what
  to install.
- **Trigger:** "I cloned this — now what do I need to set up to work on it?"
- The tool optimizes for correctness on real repos first, and fast
  time-to-answer as the experience goal.

## Scope

### What it does

- Points at a **local path or a remote git URL**.
- For remote repos, **shallow-clones into a temp directory**, scans, then cleans
  up. Host-agnostic; local paths are fully offline.
- Reads **declarative files the repo already carries** to determine
  requirements — `.tool-versions`, `.nvmrc`, `mise.toml`, `package.json`
  `engines`, `rust-toolchain.toml`, `go.mod`, `.python-version`,
  `devcontainer.json`, `Dockerfile`, `Makefile`, README prerequisites, and
  similar.
- **Checks what is installed on the user's machine** and diffs it against those
  requirements: per-tool pass/warn/fail showing "have X, need Y".
- Covers a broad range of requirements: language runtimes, dev tooling and
  package managers, containers/orchestration, cloud CLIs, database clients, and
  system-level dependencies.
- **Reports only.** It produces copy-pasteable install commands and links; it
  never installs anything or executes install steps itself.
- Presents an interactive TUI: a scrollable report, a per-tool detail pane (why
  the tool is required, which file declared it, docs link, install commands),
  filters (e.g. show only missing, filter by category), and copy-to-clipboard.
  Also offers a non-interactive **JSON output** for scripting and CI.

### Version matching

- **Loose matching:** a tool passes when the installed version is at or above
  the declared minimum / the same major; otherwise it warns (never silently
  fails).
- The declared constraint and the installed version are **always shown
  together**, so the user can see the raw comparison rather than trusting an
  opaque status.

### Repos that declare nothing

- Deterministic, declarative detection is the authoritative source. Where a repo
  or subdirectory declares nothing recognizable, the tool may show a curated
  suggestion list for the apparent language(s), **clearly framed as
  unverified** — never presented as fact.

### Multiple manifests / monorepos

- Recursively scans and **groups findings by directory/component**, so it is
  clear which package needs which version.
- Conflicts (multiple declarations of the same tool) show all declarations and
  are flagged.
- Scan respects `.gitignore` and a built-in denylist (`node_modules`, `vendor`,
  `.venv`, `dist`, `build`, …) with a sane max depth; skipped paths are reported
  so the user knows the scan was not exhaustive. No submodule traversal.

### Installed-version probing

- v1 probes `PATH` (`exec.LookPath`) and runs `--version`, parsing the output.
- Implemented behind a **swappable probing backend** so version-manager
  enumeration (nvm/asdf/mise/pyenv) can be added later without a rewrite.

### Network behavior

- Online enrichment is normal: latest upstream versions and live-validated
  links.
- The local, deterministic report is the product; enrichment is additive and
  **degrades gracefully** when a feed is unreachable or the machine is offline.
  Missing enrichment is noted inline and never blocks or invalidates the report.

## Technology and Platform

- **Go**, with a TUI library (e.g. Bubble Tea), shipped as a **single binary**.
- Targets **Linux, macOS, and Windows** in v1.
- Install metadata (per-tool name, docs URL, official install commands per
  platform) is **curated and bundled** with the binary, versioned alongside it.
- **Distribution:** prebuilt binaries on GitHub Releases plus an install script
  and `go install`; Homebrew/Scoop installers are a later follow.

## Constraints and Non-Goals (v1)

- Does **not** execute installs or run any repo code.
- Does **not** manage credentials — remote clones reuse the user's existing git
  configuration (SSH keys, credential helpers). If `git clone` fails, surface
  git's own error; the tool never touches credentials.
- **No user config and no per-repo manifest** — built-in detectors only. The
  architecture should allow config/manifests to be layered on later.
- No submodule traversal.
- No version-manager enumeration (a later milestone).
- No Windows-out: all three platforms are included in v1.

## Success Signals

- **Primary:** correct identification of the required toolset across real repos
  (accuracy).
- **Experience:** time-to-answer for a new developer drops to seconds, with no
  README reading.
- Evidence in v1: per-detector **unit tests against manifest fixtures**, plus a
  small **end-to-end corpus of real repos** to catch integration, dedupe, and
  conflict bugs.

## Automation Semantics

- Non-interactive (`--json`) mode **exits 0 when the scan succeeds**, regardless
  of missing tools — the data is the output and the caller decides.
- An opt-in `--strict` / `--check` flag returns non-zero when a required tool is
  missing or the wrong version, enabling CI/script gating.

## Open Questions

- **TUI look and feel** — colors, density, and how pass/warn/fail reads visually
  cannot be settled by discussion. Needs a static mockup or a one-screen spike.
- **Exact detector list** and per-platform install-command coverage — expected to
  grow over time.
- **Curation/quality bar** for the "unverified suggestions" list on
  under-declared repos.
- **Presentation of version conflicts** in the UI (detail pane vs. inline).
- **Which enrichment sources** to use for latest versions and link validation.
- **Whether and when** submodule traversal is added.

---

> Generated by forge-launcher on 2026-09-22T20:27:14Z
> Sharpened via forge-grill-idea on 2026-09-22.
> Use this file as input for: `@workspace /forge-auto-build-prd Use docs/IDEA.md as the project idea`
