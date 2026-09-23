// Package version exposes the build metadata for repo-ready.
//
// The Version variable is the single source of the tool's version string. It
// defaults to "dev" so that local builds are identifiable without any build
// configuration, and it is overridden at build time through the Go linker:
//
//	go build -ldflags "-X github.com/mcfuzzysquirrel/repo-ready/internal/version.Version=v1.2.3" ./cmd/repo-ready
package version

// Version is the repo-ready build version. It is "dev" unless a release build
// injects a value with -ldflags -X.
var Version = "dev"
