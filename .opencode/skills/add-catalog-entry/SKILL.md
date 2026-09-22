---
name: add-catalog-entry
description: "Add a curated tool entry to the embedded repo-ready install-metadata catalog: choose a stable id aligned with detector tool names, supply display name, category, an https documentation URL, and non-destructive Linux, macOS, and Windows install commands, then pass loader validation, the coverage test, and the human accuracy review."
---

# Skill: Add a Curated Catalog Entry

Add one tool to `internal/catalog/data/tools.json` so repo-ready can show a display
name, category, documentation link, and per-platform install commands for a required
tool. Use this when a detector emits a tool identifier the catalog does not yet
cover, or when a probe mapping needs a catalog anchor. Catalog entries are data,
not code; they are embedded in the binary and versioned with it.

## Process

### Step 1: Confirm the tool belongs in the catalog

Check the identifier is not already present and the tool fits a supported category.
If the same tool already exists under another id, then reuse the existing id
instead of adding a duplicate. Keep the catalog between 25 and 40 entries
(CATALOG-FR-01).

### Step 2: Choose a stable id aligned with detectors

Use a lowercase, canonical identifier that matches the tool name detectors emit
(`node`, `go`, `python`, `terraform`). If a catalog id and a detector tool name
disagree, then lookups silently miss; align them before adding the entry.

### Step 3: Gather official metadata

Add `id`, `name`, `category`, an `https` `docsURL`, and non-empty `linux`, `macos`,
and `windows` install commands. Commands must be the vendor's official, default
copy-pasteable installer and must be non-destructive: they install, never remove
or overwrite user data. If no official installer exists for a platform, then pick
a different tool rather than inventing a command.

### Step 4: Edit the embedded data file

Add the object to `internal/catalog/data/tools.json` and keep the file valid JSON.
Load the catalog entry schema from `references/entry-schema.md` when you need the
exact field names or the allowed category set. Do not fetch or mutate catalog data
at runtime; `go:embed` bundles it and the catalog version string versions it with
the binary (CONST-03).

### Step 5: Validate against the loader and coverage test

Run the catalog tests by default for every new entry. Validation requires unique
ids, non-empty required fields, an https docs URL, and at least one non-empty
command per platform (CATALOG-FR-02). The coverage test asserts the entry count is
in range and that every entry has all three platform commands.

### Step 6: Queue the human accuracy review

Add the new tool to the per-category sample prepared by `record-human-review` for
`docs/reviews/catalog-accuracy.json`. An agent cannot approve this gate; record the
evidence only and let a human make the decision.

## Gotchas

- **`go:embed` paths are relative to the package.** Keep the data file at `internal/catalog/data/tools.json` or the build fails before any test runs.
- **An unknown tool is not an error.** A tool absent from the catalog is still reported with its declared constraint and source (CATALOG-FR-03); never make lookup failure fatal.
- **https only.** An `http` docs URL fails catalog load; validation rejects it at startup.
- **All three platform commands are mandatory.** A missing Windows command fails coverage even when the JSON parses.
- **Ids are the join key.** When the catalog id and the detector tool name drift apart, probes and install commands silently stop resolving.
- **25 to 40 is a range, not a target to exceed.** More entries add review burden without v1 product value.
- **Commands are for the user to run.** The binary only prints them; adding an auto-install path violates CONST-01.

## Validation

Run this self-check before committing:

- [ ] `go test ./internal/catalog/...` passes, including `TestCatalogCoverage`
- [ ] `gofmt -l .` prints nothing and `go vet ./internal/catalog/...` is clean
- [ ] The new id is unique and matches the detector tool name it serves
- [ ] The docs URL is https and all three platform commands are non-empty and non-destructive
- [ ] The entry count remains between 25 and 40
- [ ] The tool is added to the human-review sample rather than marked approved

The reference `references/entry-schema.md` holds the field list and category
values when you need to pin them.
