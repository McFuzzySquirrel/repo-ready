---
name: run-quality-gates
description: "Run the repo-ready quality gate for any change: gofmt no-diff, go vet, go test ./..., a local build with version ldflag injection, and confirmation that offline tests make no network calls and opt-in tests skip by default."
---

# Skill: Run repo-ready Quality Gates

Run the standard gate for any change to repo-ready before committing. Every agent
that edits Go source, data, workflows, or fixtures uses the same gate so the result
is identical regardless of who ran it. The gate is deterministic and offline; it
never performs network I/O beyond the module cache.

## Process

### Step 1: Require a formatting no-diff

Run the formatter and require no output. If it lists files, then format them and
re-run until the output is empty. Use `gofmt -l` as the gate, never `go fmt`, which
rewrites files and always exits zero.

```bash
gofmt -l .
```

### Step 2: Run vet across the module

```bash
go vet ./...
```

### Step 3: Run the full test suite

```bash
go test ./...
```

Every package test must pass. The offline end-to-end suite is the default gate;
the network corpus and performance tests must skip cleanly unless their environment
variable is set.

### Step 4: Build with version injection

Confirm the binary builds and that the ldflag target path is correct:

```bash
go build ./...
go build -ldflags "-X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=test" ./...
make build
```

### Step 5: Confirm the offline and opt-in split

Verify no offline test performs network I/O and that no test invokes a tool
installed on the host. If a test reaches the network without its opt-in variable
set, then fix the test rather than relaxing the gate (CONST-08). Load `references/gates.md` when you need the exact commands, the CI matrix, or the documented exit codes.

### Step 6: Record the result

Report pass or the first failure with its command and output. Do not commit on a
failing gate. The recommended commit message cites the task or requirement id.

## Gotchas

- **`gofmt -l` is the only format gate.** `go fmt ./...` always exits successfully, so it can never fail a pipeline (RELEASE-FR-01).
- **Offline means no network I/O at all.** Enrichment runs only under `--enrich`; a default-path test that dials out violates CONST-08.
- **Opt-in tests must skip, not fail, by default.** A missing env var shows `SKIP`, never a connection attempt.
- **Static builds need `CGO_ENABLED=0`.** Otherwise the macOS binary may link libc and break the single-static-binary constraint (CONST-03).
- **Pinned versions are intentional.** Do not bump the pinned dependencies to clear a gate (docs/PRD.md §6.1).
- **A green local run is not a green CI run.** CI also checks Linux, macOS, and Windows; platform-specific failures belong to the owning agent.
- **Watch for the version default.** The binary reports `dev` unless the ldflag injects a value; a wrong target path silently keeps `dev`.

## Validation

Run this self-check and keep the output as evidence:

- [ ] `gofmt -l .` prints nothing
- [ ] `go vet ./...` is clean
- [ ] `go test ./...` passes
- [ ] Both `go build` invocations and `make build` succeed
- [ ] An injected version is visible and the default remains `dev` without it
- [ ] Offline tests make no network calls and opt-in tests skip by default

The reference `references/gates.md` holds the exact commands, the CI matrix, and
the exit-code contract when you need to pin them.
