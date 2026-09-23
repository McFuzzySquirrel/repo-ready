// Command repo-ready reports which developer tools a repository requires and
// whether the current environment satisfies them.
//
// This entrypoint is intentionally thin: it only reports build metadata today,
// and scanning behavior will be delegated to internal/app as that pipeline
// lands.
package main

import (
	"flag"
	"fmt"

	"github.com/mcfuzzysquirrel/repo-ready/internal/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print version information and exit")
	flag.Parse()

	if *showVersion {
		fmt.Println(version.Version)
		return
	}

	fmt.Printf("repo-ready %s\n", version.Version)
}
