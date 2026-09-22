# ADR-0001: Ship as a single static Go binary with an embedded catalog

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** release-engineer, catalog-engineer

## Context

`repo-ready` is distributed to individual developer workstations on Linux,
macOS, and Windows and is expected to run immediately after download. It needs
per-tool install metadata (display name, category, documentation URL, and
per-platform install commands) but must work with no runtime dependencies and,
by default, no network access.

## Decision

Build the product as one static Go binary (`CGO_ENABLED=0`) with no runtime
dependencies. Embed the curated install-metadata catalog with `go:embed` so the
catalog ships inside the binary and is versioned alongside it (CONST-03). Treat
Linux, macOS, and Windows as first-class targets for both runtime behavior and
release artifacts (CONST-04).

## Alternatives Considered

- **Per-platform installers / managed runtime.** Rejected: a Go single binary
  already runs everywhere and adds no interpreter or package-manager
  prerequisite.
- **Download the catalog at runtime.** Rejected: it would make the default,
  offline path depend on the network and break CONST-08.
- **Ship catalog as a sidecar file.** Rejected: a sidecar can drift from the
  binary and complicates checksum verification.

## Consequences

- Releases are cross-compiled archives for `amd64` and `arm64` per OS with
  checksums; the catalog is frozen with each binary.
- Catalog updates require a new binary release, which is acceptable because the
  product already prefers correctness and reproducibility over live data.
- Static builds must set `CGO_ENABLED=0` so macOS builds do not link libc.

## Implementation References

- Planned: `go.mod` (module `github.com/mcfuzzysquirrel/repo-ready`),
  `internal/catalog/catalog.go` with `go:embed data/tools.json`,
  `internal/catalog/data/tools.json`, `.goreleaser.yaml`,
  `.github/workflows/release.yml`.
- Requirements: `docs/PRD.md#CONST-03`, `docs/PRD.md#CONST-04`,
  `docs/features/tool-catalog.md#CATALOG-FR-01`.
