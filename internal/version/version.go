package version

import (
	_ "embed"
	"strings"
)

//go:embed VERSION
var versionFile string

// Version is the version read from the VERSION file at compile time.
var Version = strings.TrimSpace(versionFile)

