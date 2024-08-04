package buildinfo

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	// name is the application name, set at compile time.
	name string

	// version is the application version, set at compile time.
	version string

	// branch is the Git branch from which the application was built.
	branch string

	// commit is the Git commit hash of the built application.
	commit string

	// datetime is the build timestamp.
	datetime string
)

// AppName returns the application name. If not set at compile time,
// it derives the name from the executable filename.
func AppName() string {
	if name != "" {
		return name
	}
	exe := filepath.Base(os.Args[0])
	return strings.TrimSuffix(exe, ".exe")
}

// BuildString returns a formatted string containing all build information.
// This includes the application name, version, branch, commit hash,
// build datetime, and the Go runtime version.
func BuildString() string {
	return fmt.Sprintf("%s %s(%s: %s) built at %s by %s",
		AppName(), version, branch, commit, datetime, runtime.Version())
}

// PrintVersion outputs the full build information to stdout.
func PrintVersion() {
	fmt.Println(BuildString())
}

// Version returns the version string set at compile time.
func Version() string {
	return version
}
