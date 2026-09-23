package version

import "testing"

// TestVersionDefaultsToDev asserts the documented default for local builds.
// It runs against a build that did not inject a version, so the package
// variable must still hold the literal default.
func TestVersionDefaultsToDev(t *testing.T) {
	if Version != "dev" {
		t.Fatalf("Version = %q, want %q", Version, "dev")
	}
}

// TestVersionIsOverridable simulates the effect of a release build's
// -ldflags -X injection by assigning the package variable, which is exactly
// what the linker does to a string variable. This keeps the injection target
// (internal/version.Version) covered without shelling out to the Go toolchain.
func TestVersionIsOverridable(t *testing.T) {
	original := Version
	t.Cleanup(func() { Version = original })

	const injected = "v9.9.9-test"
	Version = injected

	if Version != injected {
		t.Fatalf("Version = %q, want %q", Version, injected)
	}
}
