package main

import ver "github.com/ndx-video/asciidoc-xml/internal/version"

// version is the version of the CLI tool, read from VERSION file at compile time.
// This can still be overridden at build time using ldflags:
//   go build -ldflags "-X github.com/ndx-video/asciidoc-xml/cli.version=1.0.0"
var version = ver.Version

