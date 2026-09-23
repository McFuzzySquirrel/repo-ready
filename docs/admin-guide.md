# repo-ready Administrator Guide

> [!WARNING]
> **`repo-ready` has not been implemented or released.** This guide describes
> the specified (planned) build, release, and operation procedures. No commands
> below run today. See [`docs/releases/unreleased.md`](releases/unreleased.md)
> for validation status and [`docs/PRD.md`](PRD.md) for requirements.

## Responsibilities and architecture

`repo-ready` is a **local, single-user CLI**. It is not a hosted service: there
is no server, no accounts, no database, and no multi-tenant surface. The
"administrator" role therefore does not exist in the usual sense. The roles that
*do* exist are:

- **End user** — runs the binary against a repository. Covered by the
  [user guide](user-guide.md).
- **Maintainer / operator** — builds, tests, and releases the binary and keeps
  the embedded catalog current. Covered by this guide.
- **Human reviewer** — approves catalog install-command accuracy and TUI
  legibility gates. See [`docs/features/tool-catalog.md`](features/tool-catalog.md)
  and [`docs/features/reporting.md`](features/reporting.md).

The runtime pipeline is: input resolution → scan → detect → probe → match →
report. See [ADR-0001](adr/0001-single-static-go-binary.md) and the feature
specifications under [`docs/features/`](features/).

### Version source

There is **no released version and no git tag**. The authoritative version for
any future release will be the tag that drives the release workflow; build
metadata is injected via `-ldflags` into `internal/version` and defaults to the
string `dev` when nothing is injected (FOUND-FR-03). Until a tag exists, this
project is documented as `Unreleased`.

## Prerequisites

| Purpose | Requirement |
|---------|-------------|
| Build and test | Go 1.27.1 (module `github.com/mcfuzzysquirrel/repo-ready`) |
| Remote repository inputs | `git` on `PATH` |
| Optional enrichment | Outbound HTTPS access; only when `--enrich` is passed |
| Release automation | GitHub Actions and GoReleaser (planned) |

## Installation

### Build from source

```bash
git clone https://github.com/McFuzzySquirrel/repo-ready.git
cd repo-ready
make build
```

Makefile targets (Foundation scope): `build`, `test`, `vet`, and `fmt`. The
binary is written to `bin/repo-ready` with `CGO_ENABLED=0` (static, no runtime
dependencies).

### Version injection

```bash
go build -ldflags "-X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=<version>" ./cmd/repo-ready
```

A build with no injected value reports `dev`.

> [!NOTE]
> **Module path vs repository URL casing.** The specified Go module path is
> lowercase (`github.com/mcfuzzysquirrel/repo-ready`), while the configured git
> remote is `https://github.com/McFuzzySquirrel/repo-ready.git`. GitHub treats
> these as the same repository, but the exact module path in `go.mod` determines
> what `go install` must be given. Reconcile this casing before the first
> release.

## Configuration

There is **no server configuration and no config file**. Runtime behavior is
controlled entirely by the CLI flags documented in the [user guide](user-guide.md).
v1 deliberately ships no user config and no per-repo manifest (CONST-05);
`--enrich` is the only switch that enables network access.

## Identity, secrets, and TLS

Not applicable to the runtime: `repo-ready` has no accounts, no login, no
server, and no TLS endpoint, and it never reads, stores, or manages credentials.
Remote clones reuse the operator's existing git configuration and credential
helpers (CONST-02).

For release automation, standard repository-secrets hygiene applies: release
workflows should rely on the default `GITHUB_TOKEN` scoped to the repository,
and no long-lived credentials or signing keys are needed for v1 (signing and
notarization are explicitly out of scope).

## Storage and backups

Not applicable. `repo-ready` is a stateless binary: it writes no databases,
caches, or config files, and it cleans up any temporary clone before exiting
(CONST-13). There is nothing to back up or restore.

## Health checks and monitoring

Not applicable: there is no long-running process to probe. Observability is
limited to the process exit code (`0`, `1`, `2`) and the report itself. See
[ADR-0007](adr/0007-versioned-json-and-exit-codes.md).

## Upgrades and rollback

No release exists yet, so there is no upgrade path. The planned model is:

- **Upgrade** — download the newer verified archive or re-run `go install`.
  Catalog data is embedded in the binary, so upgrading the binary upgrades the
  catalog.
- **Rollback** — reinstall the previous verified archive. Because the tool is
  stateless and writes nothing, rollback has no migration or cleanup step.

## Troubleshooting

| Symptom | Likely cause | Next action |
|---------|--------------|-------------|
| `make build` fails | Go 1.27.1 not installed, or module path/casing mismatch | Verify `go version`; reconcile the module path noted above. |
| Remote scan fails in CI | No `git` or no credentials configured for the runner | Configure the runner's git access; the tool surfaces git's stderr. |
| Build reports `dev` after release build | The `-ldflags` target path did not match `internal/version.Version` | Re-check the injection path against the module path. |
| `--json` consumers break | Emitted fields changed without a schema-version bump | Treat the JSON schema as versioned and stable; file a bug if it changes. |
| Enrichment unavailable | Feeds unreachable from the network | Expected; enrichment degrades gracefully and never changes the report. |

## Security and privacy checklist

- [ ] Confirm the binary is built statically (`CGO_ENABLED=0`) with no runtime
      dependencies (CONST-03).
- [ ] Confirm the default invocation performs **no network I/O**; `--enrich` is
      the only path that dials out (CONST-08).
- [ ] Confirm probing is shell-free, uses only known safe version queries, and
      enforces the 5-second timeout and 8-worker bound (CONST-09, CONST-11).
- [ ] Confirm remote clones use `--depth 1` into a `0700` temp dir removed on
      every exit path, including signals (CONST-13).
- [ ] Confirm no credentials, tokens, or private URLs are read, stored, logged,
      or committed (CONST-02).
- [ ] Confirm releases publish archives and a checksums file so users can verify
      downloads.
- [ ] Confirm the catalog accuracy and TUI legibility human-review gates are
      recorded and approved before release.

## Recovery and support

- **Recovery** — there is no state to recover; reinstall the verified binary.
- **Support** — file issues at
  <https://github.com/McFuzzySquirrel/repo-ready/issues>.
- **Documentation** — [README](../README.md), [user guide](user-guide.md),
  [changelog](../CHANGELOG.md), [ADRs](adr/README.md), and
  [release notes](releases/README.md).
