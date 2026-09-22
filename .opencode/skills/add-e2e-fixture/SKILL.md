---
name: add-e2e-fixture
description: "Add a committed end-to-end fixture repository under e2e/testdata and a companion assertion for repo-ready: exercise single-ecosystem, monorepo conflict, or empty-declaration behavior, assert exit codes and JSON structure, and require byte-identical repeated runs with no network or host-tool dependency."
---

# Skill: Add an End-to-End Fixture

Add a committed fixture repository under `e2e/testdata` and a companion assertion in
the `e2e` suite so repo-ready is validated against real inputs offline. Use this
when a new behavior needs end-to-end proof, a status or conflict outcome is not
covered, or a regression escaped the unit tests. Fixtures are committed source, not
generated at test time.

## Process

### Step 1: Choose the fixture case

Pick one of the covered cases: a single-ecosystem repo, a monorepo with multiple
components and a version conflict, or a repo declaring nothing. If the case
duplicates an existing fixture, then extend that fixture instead of adding a
parallel one.

### Step 2: Create committed fixture files

Add the declaring files under `e2e/testdata/<name>/`, for example
`node-rust/.nvmrc`, `monorepo/web/package.json`, and `monorepo/api/go.mod`. Keep
each fixture minimal but representative. Commit the files; never write them from
the test.

### Step 3: Invoke the real CLI surface

Run the built binary or the orchestrator so exit codes and the JSON encoder are
exercised. Load `references/fixture-matrix.md` when you need the exact cases, JSON fields, or golden-output rules. Assert the exit code, JSON structure, component grouping, and conflict flagging.

### Step 4: Keep the test offline and host-independent

Use fixture files and fake probe backends. If a fixture needs a host tool, then replace that dependency with a fake backend result. Never let an offline test depend on `node`, `go`, or `rustc` being present or on network access (CONST-08).

### Step 5: Assert deterministic repeatability

Run the same fixture twice and compare the output byte-for-byte. If the run is
flaky, then the defect is upstream ordering, not the test; fix ordering rather than
loosening the assertion (SCAN-FR-10).

### Step 6: Gate the network and performance cases

Add the opt-in network corpus and the performance budget test only behind an
explicit environment variable, and set the performance threshold to the documented
2-second budget. Both must skip by default so standard CI stays offline and
deterministic.

## Gotchas

- **Offline tests must not shell out to host tools.** Fake the probe backend instead of assuming a runtime is installed.
- **Network tests skip, not fail, by default.** A CI run without the env var must show `SKIP`, never a connection attempt.
- **Byte-identical means sorted output.** If two runs differ, the ordering bug is upstream; do not normalize in the test.
- **Fixtures are committed source.** Generating fixtures at test time reintroduces the nondeterminism the suite exists to catch.
- **The conflict fixture must declare one tool twice with different constraints.** A single declaration does not exercise PROBE-FR-06.
- **Assert exit codes explicitly.** A passing scan that returns the wrong code is a failure; 0/1/2 are part of the contract (REPORT-FR-02).
- **Keep the typical fixture small.** The performance assertion is environment-sensitive; a heavy fixture produces false failures.

## Validation

Run this self-check before committing:

- [ ] `go test ./e2e/... -run TestOffline` passes
- [ ] `gofmt -l .` prints nothing and `go vet ./e2e/...` is clean
- [ ] The single-ecosystem fixture produces the expected tools and statuses
- [ ] The monorepo fixture groups findings by component and flags the conflict
- [ ] The empty-declaration fixture produces an empty finding set with no error
- [ ] Two runs of the same fixture are byte-identical and the test makes no network calls

The reference `references/fixture-matrix.md` holds the case matrix and the JSON
assertion fields when you need to pin them.
