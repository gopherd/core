/*
Package builder provides utilities for managing and displaying build information.
It allows for compile-time injection of version details and offers methods
to retrieve and print this information.

Usage:

To use this package, import it in your main application:

	import "github.com/gopherd/core/builder"

You can then use the provided functions to access build information:

	fmt.Println(builder.AppName())     // Prints the application name
	fmt.Println(builder.Version())     // Prints the version string
	fmt.Println(builder.Info())        // Prints a string with all build details
	builder.PrintInfo()                // Prints the full build information to stdout

Compile-time configuration:

To set the build information at compile time, use the -ldflags option with go build
or go install. Here's an example:

	go build -ldflags "\
	-X github.com/gopherd/core/builder.name=myapp \
	-X github.com/gopherd/core/builder.version=v1.0.0 \
	-X github.com/gopherd/core/builder.branch=main \
	-X github.com/gopherd/core/builder.commit=abc123 \
	-X github.com/gopherd/core/builder.datetime=2023-08-04T12:00:00Z" \
	./cmd/myapp

For convenience, you can use a Makefile to automate this process:

	NAME := myapp
	VERSION := $(shell git describe --tags --always --dirty)
	BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
	COMMIT := $(shell git rev-parse --short HEAD)
	DATETIME := $(shell date +"%Y-%m-%dT%H:%M:%S%z")

	build:
	    go build -ldflags "\
	    -X github.com/gopherd/core/builder.name=$(NAME) \
	    -X github.com/gopherd/core/builder.version=$(VERSION) \
	    -X github.com/gopherd/core/builder.branch=$(BRANCH) \
	    -X github.com/gopherd/core/builder.commit=$(COMMIT) \
	    -X github.com/gopherd/core/builder.datetime=$(DATETIME)" \
	    ./cmd/myapp

This setup allows for flexible and automated injection of build information
without modifying the source code.
*/
package builder
