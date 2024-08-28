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
- **Multiple format support**: Handle JSON, TOML, YAML, and other arbitrary configuration formats through encoders and decoders
- **Automatic dependency injection**: Simplify component integration with built-in dependency resolution and injection

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

```sh
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

- `-p`: Print the configuration 🖨️
- `-t`: Test the configuration for validity ✅
- `-T`: Enable template processing for component configurations 🧩

## 🎓 Example Project

For a more comprehensive example of how to use `gopherd/core` in a real-world scenario, check out our example project:

[https://github.com/gopherd/example](https://github.com/gopherd/example)

This project demonstrates how to build a modular backend service using `gopherd/core`, including:

- Setting up multiple components
- Configuring dependencies between components
- Using the event system
- Implementing authentication and user management

It's a great resource for understanding how all the pieces fit together in a larger application! 🧩

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