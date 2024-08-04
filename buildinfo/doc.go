/*
Package buildinfo provides utilities for managing and displaying build information.
It allows for compile-time injection of version details and offers methods
to retrieve and print this information.

Usage:

To use this package, import it in your main application:

	import "github.com/gopherd/core/buildinfo"

You can then use the provided functions to access build information:

	fmt.Println(buildinfo.AppName())     // Prints the application name
	fmt.Println(buildinfo.Version())     // Prints the version string
	fmt.Println(buildinfo.BuildString()) // Prints a string with all build details
	buildinfo.PrintVersion()             // Prints the full build information to stdout

	info := buildinfo.BuildInfo()        // Returns a struct with all build information
	fmt.Printf("Built on branch: %s\n", info.Branch)

Compile-time configuration:

To set the build information at compile time, use the -ldflags option with go build
or go install. Here's an example:

	go build -ldflags "\
	-X github.com/gopherd/core/buildinfo.name=MyApp \
	-X github.com/gopherd/core/buildinfo.version=v1.0.0 \
	-X github.com/gopherd/core/buildinfo.branch=main \
	-X github.com/gopherd/core/buildinfo.commit=abc123 \
	-X github.com/gopherd/core/buildinfo.datetime=2023-08-04T12:00:00Z" \
	./cmd/myapp

For convenience, you can use a Makefile to automate this process:

	NAME := MyApp
	VERSION := $(shell git describe --tags --always --dirty)
	BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
	COMMIT := $(shell git rev-parse --short HEAD)
	DATETIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

	build:
	    go build -ldflags "\
	    -X github.com/gopherd/core/buildinfo.name=$(NAME) \
	    -X github.com/gopherd/core/buildinfo.version=$(VERSION) \
	    -X github.com/gopherd/core/buildinfo.branch=$(BRANCH) \
	    -X github.com/gopherd/core/buildinfo.commit=$(COMMIT) \
	    -X github.com/gopherd/core/buildinfo.datetime=$(DATETIME)" \
	    ./cmd/myapp

This setup allows for flexible and automated injection of build information
without modifying the source code.
*/
package buildinfo
