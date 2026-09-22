# repo-ready — Unreleased (planned v1)

**Release date:** Not released. Specification stage as of 2026-09-22.

> [!WARNING]
> This is a **draft of planned scope**, not a release. Nothing described here
> is implemented, and no artifact is available. It exists so the scope is
> reviewable before implementation begins.

## Summary

`repo-ready` is planned as a single-binary terminal app that inspects a local or
remote git repository, derives the tools and versions the repository expects
from its declarative files, diffs them against what is installed on the machine,
and reports copy-pasteable install commands and documentation links. It is
report-only: it never installs anything and never executes repository code.

## Highlights (planned)

- **Local or remote input** — scan a directory in place, or shallow-clone a git
  URL into a temporary directory that is always cleaned up.
- **Declarative detection** — eleven v1 detectors covering `.tool-versions`,
  `.nvmrc`, `mise.toml`, `package.json` engines, `rust-toolchain.toml`,
  `go.mod`, `.python-version`, `devcontainer.json`, `Dockerfile`, `Makefile`,
  and conservative README prerequisites.
- **A curated embedded catalog** — 25 to 40 tools with display name, category,
  https docs URL, and Linux/macOS/Windows install commands, versioned with the
  binary.
- **Transparent matching** — declared constraint and installed version shown
  together; mismatches warn, never fail; conflicts surfaced.
- **Interactive TUI and JSON** — a keyboard-driven report plus a single
  versioned JSON document for scripting.
- **Accessible output** — status never by color alone; `--no-color`, `NO_COLOR`,
  and a `--plain` screen-reader mode.
- **Offline by default** — no network I/O unless `--enrich` is passed.

## Compatibility and prerequisites (planned)

- **Platforms:** Linux, macOS, and Windows; `amd64` and `arm64`.
- **Install methods:** prebuilt archive, install script, `go install`, or build
  from source.
- **Runtime dependencies:** none for the binary; `git` is required only for
  remote inputs.
- **Go toolchain (build):** 1.27.1.

## Installation or upgrade

No installation is possible yet. The specified paths are documented in the
[user guide](../user-guide.md#install-and-first-use-planned) and the
[administrator guide](../admin-guide.md#installation). Upgrade and rollback
procedures will be defined once the first release exists.

## Known limitations (planned for v1)

- No version-manager enumeration (nvm, asdf, mise, pyenv); discovery is `PATH`
  based behind a swappable backend.
- No submodule traversal; skipped paths are reported.
- No user configuration and no per-repo manifest.
- No curated "unverified suggestions" for repositories that declare nothing.
- No Homebrew or Scoop packaging, and no code signing or notarization.
- Enrichment requires `--enrich` and degrades to inline "not available" notes.

## Validation

**Not validated.** No tests exist. The planned validation is an offline,
deterministic end-to-end suite over committed fixture repositories, an opt-in
env-gated network corpus that skips by default, and a performance test guarding
the 2-second budget for a typical local scan. See
[`docs/features/e2e-validation.md`](../features/e2e-validation.md). Do not cite
any test result for this project until those suites actually run.

## Security and privacy

- Report-only; never installs tools and never executes repository code
  ([ADR-0002](../adr/0002-report-only-never-install.md)).
- Credential-neutral; remote clones reuse the ambient git configuration.
- Shell-free probing with known safe version queries and a 5-second timeout.
- No telemetry; with `--enrich`, only tool identifiers and documentation URLs
  leave the machine.
- Remote clones are shallow, isolated in a `0700` temp dir, and removed on every
  exit path.

## Documentation

- [User guide](../user-guide.md)
- [Administrator guide](../admin-guide.md)
- [Architecture decisions](../adr/README.md)
- [Changelog](../../CHANGELOG.md)
- [Product requirements](../PRD.md) and [feature specifications](../features/)
