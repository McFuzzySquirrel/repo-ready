---
name: add-probe-command
description: "Add or adjust a tool installed-version probe mapping for repo-ready: validate the tool identifier against the catalog, add a known safe shell-free version-query invocation under the 5 second timeout, and cover version-output parsing with table-driven tests using a fake backend."
---

# Skill: Add a Probe Command Mapping

Add or adjust a tool's installed-version probe mapping in `internal/probe` so
repo-ready can discover an installed version without a shell. Use this when a
catalog tool is not yet probeable, a binary needs a non-default version flag, or
version-output parsing fails for a tool. Probing runs binaries directly, never
through a shell, and only with known safe version-query invocations.

## Process

### Step 1: Validate the tool identifier

Confirm the tool id exists in the catalog. If it does not, then add the catalog
entry first; only catalog-recognized identifiers may run a probe, and the id is
never interpolated into a shell string (CONST-09).

### Step 2: Choose a known safe version-query invocation

Use `--version` as the recommended default and override per tool only when the
binary documents another flag, such as `go version` or `python3 --version`. The
invocation must be a known safe version query, never a repository-provided or
arbitrary command.

### Step 3: Add the mapping explicitly

Add the tool-to-command entry to the explicit map in `internal/probe/path.go`.
Locate the binary with `exec.LookPath` and run it with `exec.CommandContext`. Load `references/command-map.md` when you need the accepted invocation forms, the 5 second timeout, or the 8-worker concurrency limit.

### Step 4: Parse the first version-like token

Parse the first version-like token from combined output and ignore banners. If the
binary is missing, then return a missing result rather than an error. If the output
contains no recognizable version, then return an unknown version. Never let either
case propagate as a fatal error.

### Step 5: Cover parsing with table-driven tests

Add rows for each output variant the tool emits: bare version, prefixed banner,
multi-line output, and no-version garbage. Use a fake backend for scheduler tests
so the suite never depends on a tool installed on the host.

### Step 6: Run the quality gate

Run the gate in `run-quality-gates` by default and confirm the invocation is
shell-free and honors the timeout before committing (PROBE-FR-02, CONST-09).

## Gotchas

- **Missing is not an error.** An absent binary yields a `missing` finding; only a failed probe setup is an error (PROBE-FR-02).
- **No shell, ever.** Use `exec.CommandContext(ctx, tool, args...)`; never interpolate a tool name into a shell string (CONST-09).
- **The 5 second timeout is a hard limit.** A hung binary must not stall the scan (CONST-11).
- **8 workers maximum, de-duplicated by tool.** Probing the same tool twice wastes the budget.
- **A non-default flag must be documented by the vendor.** Do not guess flags; an unrecognized invocation returns `unknown`, not a crash.
- **Only catalog identifiers are probeable.** An arbitrary tool name from a repository is never executed.
- **Known-safe commands only.** Never run a repository-provided script as a probe.

## Validation

Run this self-check before committing:

- [ ] `go test ./internal/probe/... ./internal/match/...` passes
- [ ] `gofmt -l .` prints nothing and `go vet ./internal/probe/...` is clean
- [ ] A missing binary yields a missing result, not an error
- [ ] The invocation uses no shell and honors the 5 second timeout
- [ ] Parsing variants are covered by table-driven tests with a fake backend
- [ ] A tool present twice is probed once and no more than 8 probes run concurrently

The reference `references/command-map.md` holds the invocation forms and limits
when you need to pin them.
