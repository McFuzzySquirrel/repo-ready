# Detector fixture shape

Fixtures are committed source under `testdata/detect/<kind>/`. The test reads the
fixture through the registry, not through a private parser, so the fixture also
proves registration.

## Directory layout

```
testdata/detect/
  toolversions/.tool-versions
  nvmrc/.nvmrc
  packagejson/package.json
  gomod/go.mod
  readme/README.md
  malformed/package.json
```

One directory per file kind. Add a `malformed/` variant when the detector must
return an error, and keep it separate so the error case cannot silently succeed.

## Table-driven test shape

```go
func TestDotfileDetectors(t *testing.T) {
    cases := []struct {
        name    string
        dir     string
        want    []model.Requirement
        wantErr bool
    }{
        // one row per format and per boundary (plain pin, major-only, multi-tool)
    }
    for _, tc := range cases {
        t.Run(tc.name, func(t *testing.T) { /* run registry, compare */ })
    }
}
```

Each family test is named exactly once: `TestDotfileDetectors`,
`TestManifestDetectors`, `TestTextDetectors`. A new file kind adds rows, not a
new top-level test.

## Coverage checklist per fixture

- A plain pin, a major-only pin, and a multi-tool file are each represented.
- The malformed case asserts `wantErr` and an empty requirement slice.
- The README fixture includes both a fenced prerequisite block and a prose-only
  mention, and asserts the prose yields nothing.
- Requirements are compared with ordering, so a sort regression fails the test.
