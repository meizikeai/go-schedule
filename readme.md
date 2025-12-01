# go-schedule

A practice project built with Go (Golang), following clean architecture principles and common enterprise-level project structure.

## Tech Stack

- Web framework: [Gin](https://github.com/gin-gonic/gin)
- Database: Supports [MySQL](https://github.com/go-sql-driver/mysql)
- Cache: Supports [Redis](https://github.com/go-redis/redis)
- Runtime: Requires [Go](https://go.dev/dl/) environment

## Project Structure

```bash
go-schedule
|-- cmd
|   |-- main.go             → Application entry point
|-- internal
|   |-- app                 → Dependency injection module
|   |   |-- repository      → Data access layer
|   |   |-- task            → Business logic layer
|   |   |-- app.go
|-- |-- config              → Configuration files and settings
|-- |-- dto                 → Data Transfer Objects (request/response structs)
|-- |-- model               → Domain models and entity definitions
|-- |-- pkg                 → Utility packages (common reusable functions)
|-- go.mod                  → Go module and dependency management
|-- README.md               → Project documentation
```

## Quick Start

```bash
# Clone the repository
git clone https://github.com/meizikeai/go-schedule.git
cd go-schedule

# Download dependencies
go mod tidy

# Run the application
cd cmd && go run .
```

## Recommended Development Environment

- Editor: [Visual Studio Code](https://code.visualstudio.com)
- Go extension and tools: Refer to the official [Go tools for VS Code](https://code.visualstudio.com/docs/languages/go)

Please make sure all Go tools are properly installed before starting development.

## Helpful Resources

- [The Go Programming Language](https://go.dev)
- [Go Standard Library](https://pkg.go.dev/std)
- [Effective Go](https://go.dev/doc/effective_go)
- [Go by Example](https://gobyexample.com)
- [Go Web Examples](https://gowebexamples.com)

Feel free to explore, modify, and use this project as a reference for building production-ready Go applications!
