# Probe command map reference

Authoritative sources: `docs/features/probing-and-matching.md` PROBE-FR-01..03 and
PROBE-1/2, `docs/PRD.md` CONST-09 (no-shell probing), CONST-11 (5s timeout, 8
workers), and `internal/probe` (owned by probe-engineer).

## Map shape

```go
var commandMap = map[string][]string{
    "go":     {"version"},
    "node":   {"--version"},
    "python": {"--version"},
    "rust":   {"--version"},
    "jq":     {"--version"},
}
```

The default invocation is `--version`. Override only when the vendor documents a
different version flag, as `go version` does. The map key must match a catalog id.

## Execution rules

- Resolve the binary with `exec.LookPath`; if it is absent, return a missing result.
- Run with `exec.CommandContext(ctx, binary, args...)`; never build a shell string.
- Apply a 5 second context timeout per command; cancel on scheduler shutdown.
- Schedule at most 8 probes concurrently and de-duplicate identical tool ids.
- Capture combined output and parse the first version-like token.

## Parsing table

| Tool output | Expected version |
|-------------|------------------|
| `v20.11.1` | `20.11.1` |
| `go version go1.27.1 linux/amd64` | `1.27.1` |
| `Python 3.12.4` | `3.12.4` |
| `jq-1.7.1` | `1.7.1` |
| `not a version` | unknown |

Stripping a leading `v` and the `go`/`Python` banner words is part of parsing, not
part of the tool name. The tool identifier stays exactly the catalog id.

## Failure mapping

| Situation | Result |
|-----------|--------|
| Binary not on PATH | `missing` finding, no error |
| Non-zero exit | `missing` or `unknown`, never fatal |
| Timeout exceeded | `unknown`, never fatal |
| No version-like token | `unknown`, never fatal |
| Catalog id unknown | Do not probe; the tool is reported without catalog metadata |
