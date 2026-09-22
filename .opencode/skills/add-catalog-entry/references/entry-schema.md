# Catalog entry schema reference

Authoritative sources: `docs/features/tool-catalog.md` CATALOG-FR-01..03 and
CATALOG-1/2, `docs/PRD.md` CONST-03 and CONST-01, and `internal/catalog/catalog.go`
(owned by catalog-engineer).

## Entry fields

```json
{
  "id": "node",
  "name": "Node.js",
  "category": "language runtime",
  "docsURL": "https://nodejs.org/en/docs",
  "install": {
    "linux": "curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash - && sudo apt-get install -y nodejs",
    "macos": "brew install node",
    "windows": "winget install OpenJS.NodeJS.LTS"
  }
}
```

Field contract, matching the loader in CATALOG-1:

| Field | Required | Rule |
|-------|----------|------|
| `id` | yes | Unique, lowercase, matches the detector tool name |
| `name` | yes | Non-empty human display name |
| `category` | yes | One value from the category set below |
| `docsURL` | yes | Non-empty, `https` scheme, official documentation |
| `install.linux` | yes | Non-empty, official, non-destructive |
| `install.macos` | yes | Non-empty, official, non-destructive |
| `install.windows` | yes | Non-empty, official, non-destructive |

Include a catalog version string alongside the entries so the embedded data is
versioned with the binary.

## Category set

Keep category values stable; the coverage test asserts the intended spread and the
report groups and filters by them:

- language runtime
- package manager
- container
- orchestration
- cloud CLI
- database client
- system dependency

## Validation failure modes

The loader fails the whole catalog, with a clear error, when any of these hold:

- two entries share an `id`
- `name`, `category`, or `docsURL` is empty
- `docsURL` is not `https`
- a platform command is missing or empty

An unknown `id` at lookup time is **not** a load failure; it returns a not-found
result so the caller can report the tool without catalog metadata.
