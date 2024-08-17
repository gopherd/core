// Package builder provides utilities for managing and displaying build information.
//
// It allows for compile-time injection of version details and offers methods
// to retrieve and print this information.
//
// Usage:
//
// To use this package, import it in your main application:
//
//	import "github.com/gopherd/core/builder"
//
// You can then use the provided functions to access build information:
//
//	fmt.Println(builder.Info())     // Prints the version string
//	builder.PrintInfo()             // Prints a string with all build details
//
// Compile-time configuration:
//
// To set the build information at compile time, use the -ldflags option with go build
// or go install. Here's an example:
//
//	go build -ldflags "\
//	-X github.com/gopherd/core/builder.name=myapp \
//	-X github.com/gopherd/core/builder.version=v1.0.0 \
//	-X github.com/gopherd/core/builder.branch=main \
//	-X github.com/gopherd/core/builder.commit=abc123 \
//	-X github.com/gopherd/core/builder.datetime=2023-08-04T12:00:00Z" \
//	./cmd/myapp
//
// For convenience, you can use a Makefile to automate this process:
//
//	BUILD_NAME := myapp
//	BUILD_VERSION := $(shell git describe --tags --always --dirty)
//	BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
//	BUILD_COMMIT := $(shell git rev-parse --short HEAD)
//	BUILD_DATETIME := $(shell date +"%Y-%m-%dT%H:%M:%S%z")
//	BUILD_PKG := github.com/gopherd/core/builder
//
//	build:
//		go build -ldflags "\
//		-X $(BUILD_PKG).name=$(BUILD_NAME) \
//		-X $(BUILD_PKG).version=$(BUILD_VERSION) \
//		-X $(BUILD_PKG).branch=$(BUILD_BRANCH) \
//		-X $(BUILD_PKG).commit=$(BUILD_COMMIT) \
//		-X $(BUILD_PKG).datetime=$(BUILD_DATETIME)" \
//		./cmd/myapp
//
// This setup allows for flexible and automated injection of build information
// without modifying the source code.
package builder

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

var (
	name     string // Application name, set at compile time
	version  string // Application version, set at compile time
	branch   string // Git branch from which the application was built
	commit   string // Git commit hash of the built application
	datetime string // Build timestamp
)

// appName returns the application name. If not set at compile time,
// it derives the name from the executable filename.
func appName() string {
	if name != "" {
		return name
	}
	exe, _ := os.Executable()
	return strings.TrimSuffix(filepath.Base(exe), ".exe")
}

// buildInfo contains all build information.
type buildInfo struct {
	Name     string
	Version  string
	Branch   string
	Commit   string
	DateTime string
}

var runtimeVersion = runtime.Version

// String returns a formatted string containing all build information.
func (info buildInfo) String() string {
	return fmt.Sprintf("%s %s(%s: %s) built at %s by %s",
		info.Name, info.Version, info.Branch, info.Commit, info.DateTime, runtimeVersion())
}

// Info returns a struct containing the build information.
func Info() buildInfo {
	return buildInfo{
		Name:     appName(),
		Version:  version,
		Branch:   branch,
		Commit:   commit,
		DateTime: datetime,
	}
}

// PrintInfo outputs the full build information to stdout.
func PrintInfo() {
	fmt.Println(Info().String())
}
