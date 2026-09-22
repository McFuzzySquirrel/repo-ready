# Changelog

All notable changes to `repo-ready` are documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project intends to follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

> [!IMPORTANT]
> **No version of `repo-ready` has been released.** The repository currently
> contains specifications, an agent team, and project skills only: there is no
> Go source, `go.mod`, `Makefile`, CI workflow, release configuration, or built
> binary. Every capability described elsewhere in the docs is **planned**, not
> implemented. The planned v1 scope is captured in
> [docs/releases/unreleased.md](docs/releases/unreleased.md); durable design
> decisions are recorded in [docs/adr/](docs/adr/README.md).

## [Unreleased]

### Added

- Product requirements, [`docs/PRD.md`](docs/PRD.md), and eight feature
  specifications under [`docs/features/`](docs/features/), decomposed from
  [`docs/IDEA.md`](IDEA.md).
- A nine-agent engineering team under [`.opencode/agents/`](.opencode/agents/)
  and six project skills under [`.opencode/skills/`](.opencode/skills/).
- Execution planning artifacts: [`docs/EXECUTION-MANIFEST.json`](docs/EXECUTION-MANIFEST.json)
  and [`docs/agent-responsibility-matrix.md`](docs/agent-responsibility-matrix.md).
- Project documentation set: this changelog, the [ADRs](docs/adr/README.md),
  the [user guide](docs/user-guide.md), the [administrator guide](docs/admin-guide.md),
  and [release notes](docs/releases/README.md).

### Notes

- No functional changes are recorded yet because no implementation exists.
- All runtime behavior described in the user and administrator guides is the
  specified design, pending implementation.
