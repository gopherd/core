# gopherd/core

[![Go Reference](https://pkg.go.dev/badge/github.com/gopherd/core.svg)](https://pkg.go.dev/github.com/gopherd/core)
[![Go Report Card](https://goreportcard.com/badge/github.com/gopherd/core)](https://goreportcard.com/report/github.com/gopherd/core)
[![codecov](https://codecov.io/gh/gopherd/core/branch/main/graph/badge.svg)](https://codecov.io/gh/gopherd/core)
[![Build Status](https://github.com/gopherd/core/workflows/Go/badge.svg)](https://github.com/gopherd/core/actions)
[![License](https://img.shields.io/github/license/gopherd/core.svg)](https://github.com/gopherd/core/blob/main/LICENSE)

`gopherd/core` is a powerful Go library that provides a component-based development framework for building backend services, along with a set of essential utility functions.

## Overview

This library offers a robust mechanism for component-based development, enabling Go developers to create highly modular and maintainable backend services. By leveraging `gopherd/core`, developers can:

- Easily create and manage components for various functionalities (e.g., database connections, caching, authentication)
- Implement a plugin-like architecture for extensible services
- Utilize a set of fundamental helper functions to streamline common tasks

The component-based approach allows for better organization, reusability, and scalability of your Go backend services.

## Key Features

- Component-based architecture for modular service development
- Simplified configuration management for components
- Easy integration of custom components
- Helper functions for various low-level operations

## Installation

To use `gopherd/core` in your Go project, install it using `go get`:

```bash
go get github.com/gopherd/core
```

## Usage Example

Here's a simplified example of how to use `gopherd/core` to create a modular backend service:

```go
package main

import (
    "github.com/gopherd/core/component"
    "github.com/gopherd/core/config"
    "github.com/gopherd/core/raw"
    "github.com/gopherd/core/service"

    // your components
    "github.com/your/components/logger"
    "github.com/your/components/db"
    "github.com/your/components/httpserver"
    "github.com/your/components/blockexit"
)

func main() {
    // Define context for your service
    var context struct {
        Name       string
        ID         int
    }
    context.Name = "MyService"
    context.ID = 1

    // Run the service with components
    service.Run(service.NewBaseService(config.NewBaseConfig(
        context,
        []component.Config{
            {
                Name:    logger.Name,
                Options: raw.MustJSON(logger.DefaultOptions(nil)),
            },
            {
                Name: db.Name,
                Options: raw.MustJSON(db.Options{
                    Driver: "postgres",
                    DSN:    "host=localhost user=myuser password=mypassword dbname=mydb sslmode=disable",
                }),
            },
            {
                Name: httpserver.Name,
                Options: raw.MustJSON(httpserver.Options{
                    Addr: ":8080",
                }),
            },

            // Add more custom components to handle your business logic, such as:
            // - Auth component for processing authorization
            // - Users component for managing application users
            // ...

            {
                Name: blockexit.Name,
            },
        },
    )))
}
```

This example demonstrates:

1. Configuring multiple components (logger, database, HTTP server, ...)
2. Running the service with the configured components

Remember to implement your custom business logic within these components or create additional components as needed for your specific use case.

## Documentation

For detailed documentation of each package and component, please refer to the GoDoc:

[https://pkg.go.dev/github.com/gopherd/core](https://pkg.go.dev/github.com/gopherd/core)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the [MIT License](LICENSE).

## Support

If you encounter any problems or have any questions, please open an issue in this repository.

---

We hope you find `gopherd/core` valuable for your Go backend projects! Whether you're building a small microservice or a complex distributed system, `gopherd/core` provides the foundation for creating modular, maintainable, and efficient backend services.