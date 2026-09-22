# End-to-end fixture matrix

Authoritative sources: `docs/features/e2e-validation.md` E2E-FR-01..03 and
E2E-1/2, plus `docs/features/reporting.md` REPORT-FR-02/03 for the exit codes and
JSON fields asserted here.

## Case matrix

| Case | Fixture path | Declares | Asserts |
|------|--------------|----------|---------|
| Single ecosystem | `e2e/testdata/node-rust/.nvmrc` | One runtime pin | Expected tools and statuses; exit 0 |
| Monorepo conflict | `e2e/testdata/monorepo/web/package.json`, `e2e/testdata/monorepo/api/go.mod` | Two components, same tool twice with different constraints | Component grouping and the conflict flag; exit 0 |
| Empty declaration | `e2e/testdata/empty/README.md` | Nothing | Empty finding set, no error; exit 0 |
| Strict gate | any missing/outdated fixture | A missing or outdated tool | Exit 1 only with `--strict` |

## JSON fields to assert

The `--json` document is a single versioned object. Assert at least:

- `schemaVersion` is present and stable
- `scan.skippedPaths` and `scan.depth` are present
- `components` group findings by component directory
- each finding carries `declared` and `installed` values together
- a conflict finding lists every declaration and sets the conflict flag
- the enrichment section is omitted when enrichment did not run

## Assertion shape

```go
func TestOfflineMonorepoConflict(t *testing.T) {
    out := runBinary(t, "e2e/testdata/monorepo", "--json")
    var report Report
    decode(t, out, &report)
    // assert two components, one conflict-flagged tool, both declarations present
}
```

Prefer invoking the built binary or the `internal/app` orchestrator over a private
helper so the real flag parsing and exit codes are covered.

## Opt-in cases

| Test | Gate | Default |
|------|------|---------|
| Network corpus | `REPO_READY_NETWORK_TESTS` (or the documented variable) | Skips before any setup |
| Performance budget | always runs against the local fixture | Fails above the 2s budget |

The network corpus targets two or three small public repositories configured in the
test. Never let it run during a default `go test ./...`.
