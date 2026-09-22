# ADR-0002: Report-only — never install or execute repository code

- **Status:** Proposed
- **Date:** 2026-09-22
- **Decision owners:** cli-engineer, scan-engineer, probe-engineer

## Context

`repo-ready` inspects untrusted, often unfamiliar repositories. A tool that
installs prerequisites or executes repository scripts would become a
supply-chain risk and would need credential handling and elevated privileges.
The product's value is a fast, trustworthy answer, not automation of setup.

## Decision

The tool only **reports** requirements, install commands, and documentation
links for the user to run. It never installs tools and never executes
repository code or scripts (CONST-01). It never reads, stores, or manages
credentials; remote clones reuse the user's ambient git configuration and
credential helpers, and on failure it surfaces git's own error and exits
non-zero (CONST-02). Command probing runs binaries directly without a shell,
using only known safe version queries (CONST-09).

## Alternatives Considered

- **Offer an `--install` mode.** Rejected: it would require running package
  managers with user privileges and is explicitly out of scope for v1.
- **Execute repository setup scripts to infer requirements.** Rejected: it
  contradicts the report-only guarantee and the offline/no-execution posture.

## Consequences

- The blast radius of scanning an untrusted repo is limited to reading files,
  running `git clone` for remote inputs, and querying versions of known tools.
- Users copy-paste commands themselves, so install-command accuracy is a
  human-reviewed concern (see the catalog accuracy review gate).
- No credential store, secret scanning, or privileged operation is needed.

## Implementation References

- Planned: `internal/input/clone.go`, `internal/probe/path.go`,
  `internal/detect/*` (pure functions over files).
- Requirements: `docs/PRD.md#CONST-01`, `docs/PRD.md#CONST-02`,
  `docs/PRD.md#CONST-09`, `docs/features/scan-and-detect.md#SCAN-FR-04`.
