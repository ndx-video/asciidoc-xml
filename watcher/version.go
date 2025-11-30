package main

import ver "github.com/ndx-video/asciidoc-xml/internal/version"

// version is the version of the watcher daemon, read from VERSION file at compile time.
// This can still be overridden at build time using ldflags:
//   go build -ldflags "-X github.com/ndx-video/asciidoc-xml/watcher.version=1.0.0"
var version = ver.Version

