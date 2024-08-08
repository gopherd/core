# 🚀 gopherd/core

[![Go Reference](https://pkg.go.dev/badge/github.com/gopherd/core.svg)](https://pkg.go.dev/github.com/gopherd/core)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopherd/core)](https://goreportcard.com/report/github.com/gopherd/core)
[![codecov](https://codecov.io/gh/gopherd/core/branch/main/graph/badge.svg)](https://codecov.io/gh/gopherd/core)
[![Build Status](https://github.com/gopherd/core/workflows/Go/badge.svg)](https://github.com/gopherd/core/actions)
[![License](https://img.shields.io/github/license/gopherd/core.svg)](https://github.com/gopherd/core/blob/main/LICENSE)

`gopherd/core` is a Go library that provides a component-based development framework for building backend services, leveraging the power of Go's generics. It's a modern, type-safe approach to creating scalable applications! 🌟

## 🌟 Overview

This library offers a state-of-the-art mechanism for component-based development, enabling Go developers to create highly modular and maintainable backend services. By harnessing the power of `gopherd/core` and Go's generics, developers can:

- 🧩 Easily create and manage type-safe components for various functionalities (e.g., database connections, caching, authentication)
- 🔌 Implement a plugin-like architecture for extensible services with compile-time type checking
- 🛠️ Utilize a set of fundamental helper functions to streamline common tasks, all with the benefits of generics

The component-based approach, combined with Go's generics, allows for better organization, reusability, and scalability of your Go backend services. It's like LEGO for your code, but with perfect fit guaranteed by the type system! 🧱✨

## 🔥 Key Features

- **Modern, generic-based architecture**: Leverage Go's generics for type-safe component creation and management
- **Flexible configuration**: Load configurations from files, URLs, or standard input with type safety
- **Template processing**: Use Go templates in your component configurations for dynamic setups

## 📦 Installation

To use `gopherd/core` in your Go project, install it using `go get`:

```bash
go get github.com/gopherd/core
```

## ⚡ Quick Start

Here's a simple example showcasing the power of generics in our library:

```go
// demo/main.go
package main

import (
	"context"
	"fmt"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/service"
)

// helloComponent demonstrates the use of generics for type-safe options.
type helloComponent struct {
	component.BaseComponent[struct {
		Message string
	}]
}

func (c *helloComponent) Init(ctx context.Context) error {
	fmt.Println("Hello, " + c.Options().Message + "!")
	return nil
}

func init() {
	component.Register("hello", func() component.Component { return &helloComponent{} })
}

func main() {
	service.Run()
}
```

Yes, it's that simple and type-safe! 😮 With just these few lines of code, you can leverage the power of our generic-based, component-driven architecture. The simplicity of this example demonstrates how our library abstracts away the complexities of component management while maintaining type safety, allowing you to focus on building your application logic. Modern magic, right? ✨🔮

### Basic Usage

Run your application with a configuration file:

```
./demo app.json
```

> Here's an example `app.json`

```json
{
	"Components": [
		{
			"Name": "hello",
			"Options": {
				"Message": "world"
			}
		}
	]
}
```

### Load Configuration from Different Sources

- From a file: `./demo app.json` 📄
- From a URL: `./demo http://example.com/config/app.json` 🌐
- From stdin: `echo '{"Components":[...]}' | ./demo -` ⌨️

### Command-line Options

- `-v`: Print version information 🏷️
- `-p`: Print loaded configuration 🖨️
- `-t`: Test the configuration for validity ✅
- `-T`: Enable template processing for component configurations 🧩

### Example with Template Processing

```json
{
	"Context": {
		"Name": "world"
	},
	"Components": [
		{
			"Name": "hello",
			"Options": {
				"Message": "{{.Name}}"
			}
		}
	]
}
```

This example demonstrates how to use template processing in your component configurations. Cool, huh? 😎

### Help

For a full list of options and usage examples, run:

```
./demo -h
```

This will display the following help information:

```
Usage: ./demo [OPTIONS] <config>
       ./demo <path/to/file>   (Read config from file)
       ./demo <url>            (Read config from http)
       ./demo -                (Read config from stdin)
       ./demo -v               (Print version information)
       ./demo -p               (Print loaded config)
       ./demo -t               (Test the config for validity)
       ./demo -T               (Enable template processing for components config)

Examples:
       ./demo app.json
       ./demo http://example.com/app.json
       echo '{"Components":[{"Name":"$hello","Options":{"Message":"world"}}]}' | ./demo -
       ./demo -p app.json
       ./demo -t app.json
       ./demo -T app.json
       ./demo -p -T app.json
       ./demo -t -T app.json
```

## 📚 Documentation

For detailed documentation of each package and component, please refer to the GoDoc:

[https://pkg.go.dev/github.com/gopherd/core](https://pkg.go.dev/github.com/gopherd/core)

## 👥 Contributing

Contributions are welcome! Please feel free to submit a Pull Request. Let's make this library even more awesome together! 🤝

## 📜 License

This project is licensed under the [MIT License](LICENSE).

## 🆘 Support

If you encounter any problems or have any questions, please open an issue in this repository. We're here to help! 💪

---

We hope you find `gopherd/core` valuable for your modern Go backend projects! Whether you're building a small microservice or a complex distributed system, `gopherd/core` provides the foundation for creating modular, maintainable, and efficient backend services with the power of generics. Welcome to the future of Go development! 🚀🎉
