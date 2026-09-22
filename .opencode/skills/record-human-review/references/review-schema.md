# Human-review file schema

Two v1 gates share this evidence shape. Authoritative sources:
`docs/features/tool-catalog.md` CATALOG-3 and `docs/features/reporting.md`
REPORT-6.

## Gates

| Gate id | File | Reviewer evaluates |
|---------|------|--------------------|
| `catalog-accuracy` | `docs/reviews/catalog-accuracy.json` | Docs URLs and per-platform install commands are correct, current, non-destructive |
| `tui-legibility` | `docs/reviews/tui-legibility.json` | Status reads without color, declared/installed are unambiguous, detail is actionable, keys discoverable |

## Common envelope

```json
{
  "gate": "catalog-accuracy",
  "schemaVersion": 1,
  "reviewer": "human-name-or-handle",
  "reviewedAt": "2026-09-22T00:00:00Z",
  "decision": "pending",
  "sample": [],
  "observations": [],
  "corrections": []
}
```

`decision` is one of `pending`, `approved`, or `changes-requested`. Only a human
reviewer may set it away from `pending`.

## Catalog observation

```json
{
  "id": "node",
  "category": "language runtime",
  "docsURL": { "url": "https://nodejs.org/en/docs", "resolves": true },
  "commands": {
    "linux": { "command": "apt-get install -y nodejs", "observed": "official" },
    "macos": { "command": "brew install node", "observed": "official" },
    "windows": { "command": "winget install OpenJS.NodeJS.LTS", "observed": "official" }
  }
}
```

## TUI observation

```json
{
  "mode": "plain",
  "transcript": "path or inline text",
  "statusReadsWithoutColor": true,
  "declaredAndInstalledUnambiguous": true,
  "detailActionable": true,
  "keysDiscoverable": true
}
```

## Correction and disposition

```json
{
  "target": "node.install.windows",
  "severity": "medium",
  "issue": "command installs a non-LTS channel",
  "disposition": "applied",
  "rationale": "switched to the LTS package id"
}
```

`disposition` is one of `applied`, `accepted-with-rationale`, or `deferred`. Every
correction must carry one; an omitted disposition blocks the gate.
