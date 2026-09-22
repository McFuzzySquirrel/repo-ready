# Quality gate reference

Authoritative sources: `docs/features/release-and-distribution.md` RELEASE-FR-01
and RELEASE-1, `docs/features/project-foundation.md` FOUND-1, and `docs/PRD.md`
CONST-03, CONST-08, CONST-11.

## Commands

```bash
gofmt -l .                                   # must print nothing
go vet ./...                                 # must be clean
go test ./...                                # must pass
go build ./...
go build -ldflags "-X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=test" ./...
make build
bash -n scripts/install.sh                   # installer changes only
```

`make build`, `make test`, `make vet`, and `make fmt` wrap the same Go toolchain.

## CI matrix

`.github/workflows/ci.yml` runs gofmt no-diff, `go vet ./...`, and `go test ./...`
on `ubuntu`, `macos`, and `windows` for every push and pull request. Any formatting
diff, vet finding, or failing test fails the job.

## Offline and opt-in test contract

| Test class | Default behavior | Gate |
|------------|------------------|------|
| Unit and offline e2e | Run, no network | Must pass (CONST-08) |
| Network corpus | Skip unless the opt-in env var is set | Must show `SKIP` by default |
| Performance budget | Run against a local fixture, assert 2s budget | Must pass (CONST-11) |

The network corpus must call `t.Skip` before any setup when its env var is absent,
so a plain `go test ./...` never dials out.

## Exit-code contract

The built CLI returns:

| Code | Meaning |
|------|---------|
| 0 | Successful scan |
| 1 | `--strict` set and at least one finding is missing or outdated |
| 2 | Fatal input, clone, or usage error |

## Version injection

Build metadata is injected through
`github.com/mcfuzzysquirrel/repo-ready/internal/version.Version` and defaults to
`dev`. Verify the injected value reaches `--version`; a typo in the target path
leaves `dev` in place without failing the build.
