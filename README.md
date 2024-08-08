# 🚀 gopherd/core

[![Go Reference](https://pkg.go.dev/badge/github.com/gopherd/core.svg)](https://pkg.go.dev/github.com/gopherd/core)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopherd/core)](https://goreportcard.com/report/github.com/gopherd/core)
[![codecov](https://codecov.io/gh/gopherd/core/branch/main/graph/badge.svg)](https://codecov.io/gh/gopherd/core)
[![Build Status](https://github.com/gopherd/core/workflows/Go/badge.svg)](https://github.com/gopherd/core/actions)
[![License](https://img.shields.io/github/license/gopherd/core.svg)](https://github.com/gopherd/core/blob/main/LICENSE)

`gopherd/core` is a powerful Go library that provides a component-based development framework for building backend services, along with a set of essential utility functions. Let's dive in! 💡

## 🌟 Overview

This library offers a mechanism for component-based development, enabling Go developers to create highly modular and maintainable backend services. By leveraging `gopherd/core`, developers can:

- 🧩 Easily create and manage components for various functionalities (e.g., database connections, caching, authentication)
- 🔌 Implement a plugin-like architecture for extensible services
- 🛠️ Utilize a set of fundamental helper functions to streamline common tasks

The component-based approach allows for better organization, reusability, and scalability of your Go backend services. It's like LEGO for your code! 🧱

## 🔥 Key Features

- **Component-based architecture**: Easily create and manage reusable components
- **Flexible configuration**: Load configurations from files, URLs, or standard input
- **Template processing**: Use Go templates in your component configurations

## 📦 Installation

To use `gopherd/core` in your Go project, install it using `go get`:

```bash
go get github.com/gopherd/core
```

## ⚡ Quick Start

Here's a simple example of how to use the library:

```go
// demo/main.go
package main

import (
	"context"
	"fmt"

	"github.com/gopherd/core/component"
	"github.com/gopherd/core/service"
)

// helloComponent is a simple example of a component.
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

Yes, it's that simple! 😮 With just these few lines of code, you can leverage the power of our component-based architecture. The simplicity of this example demonstrates how our library abstracts away the complexities of component management, allowing you to focus on building your application logic. Magic, right? ✨

### Basic Usage

Run your application with a configuration file:

```
./demo app.json
```

> Here's an example `app.json`

```json
{
	"Context": {
		"Name": "world"
	},
	"Components": [
		{
			"Name": "hello",
			"Options": {
				"Msg": "{{.Name}}"
			}
		}
	]
}
```


### Load Configuration from Different Sources

- From a file: `./demo app.json` 📄
- From a URL: `./demo http://example.com/app.json` 🌐
- From stdin: `echo '{"Components":[...]}' | ./demo -` ⌨️

### Command-line Options

- `-v`: Print version information 🏷️
- `-p`: Print loaded configuration 🖨️
- `-t`: Test the configuration for validity ✅
- `-T`: Enable template processing for component configurations 🧩

### Example with Template Processing

```
echo '{"Context":{"Name":"world"},"Components":[{"Name":"hello","Options":{"Msg":"{{.Name}}"}}]}' | go run . -T -
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

We hope you find `gopherd/core` valuable for your Go backend projects! Whether you're building a small microservice or a complex distributed system, `gopherd/core` provides the foundation for creating modular, maintainable, and efficient backend services. Happy coding! 🎉