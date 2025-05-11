# gopher-tasks

> A simple Task Management Service (API + CLI) in Go, built with Clean Architecture and SOLID principles.

## Table of Contents

- [gopher-tasks](#gopher-tasks)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Features](#features)
  - [Tech Stack](#tech-stack)
  - [Prerequisites](#prerequisites)
  - [Installation](#installation)
  - [Configuration](#configuration)
  - [Usage](#usage)
    - [Run the API Server](#run-the-api-server)
    - [Use the CLI](#use-the-cli)
  - [Project Structure](#project-structure)
  - [Testing](#testing)
  - [Contributing](#contributing)
  - [License](#license)

---

## Overview

**gopher-tasks** is a learning project that exposes:

- A RESTful HTTP API to create, read, update and delete tasks.  
- A Command-Line Interface (CLI) powered by [cobra](https://github.com/spf13/cobra) to manage your tasks from the terminal.  

It’s organized following Clean Architecture:

- **Domain**: Entities & interfaces.  
- **Use Cases**: Business logic.  
- **Delivery**: HTTP handlers & CLI commands.  
- **Infrastructure**: Concrete implementations (PostgreSQL, JWT, config).

---

## Features

- User registration & login with JWT authentication  
- CRUD operations for tasks (title, description, due date, status)  
- Filtering & pagination of tasks  
- CLI commands:  
  - `gopher-tasks login`  
  - `gopher-tasks task create "Buy milk"`  
  - `gopher-tasks task list --status=pending`  
  - `gopher-tasks task complete <id>`  
- Configurable via YAML or environment variables  
- Unit & integration tests  

---

## Tech Stack

- **Language**: Go 1.21+  
- **HTTP Router**: [chi](https://github.com/go-chi/chi) or `net/http` + middleware  
- **CLI**: [cobra](https://github.com/spf13/cobra)  
- **Database**: PostgreSQL (or SQLite for quick start)  
- **Migrations**: [golang-migrate](https://github.com/golang-migrate/migrate)  
- **Logger**: [logrus](https://github.com/sirupsen/logrus) or [zap](https://github.com/uber-go/zap)  
- **Testing**: `testing`, [testify](https://github.com/stretchr/testify), [testcontainers-go](https://github.com/testcontainers/testcontainers-go)  

---

## Prerequisites

- Go 1.21+ installed  
- PostgreSQL database (or SQLite)  
- `migrate` CLI for database migrations (optional)  
- [Make](https://www.gnu.org/software/make/) (optional, for convenience)  

---

## Installation

Clone the repo and download dependencies:

```bash
git clone https://github.com/seu-usuario/gopher-tasks.git
cd gopher-tasks
go mod download
```

---

## Configuration

Copy the example config and adjust:

```bash
cp configs/config.example.yaml configs/config.yaml
```

Edit `configs/config.yaml`:

```yaml
server:
  port: 8080

database:
  driver: postgres
  dsn: "postgres://user:pass@localhost:5432/gophertasks?sslmode=disable"

auth:
  jwt_secret: "your-secret-key"
  token_expiry_minutes: 60
```

You may also override via environment variables (`SERVER_PORT`, `DATABASE_DSN`, `AUTH_JWT_SECRET`, …).

---

## Usage

### Run the API Server

```bash
cd cmd/server
go run main.go
```

The API will be available at `http://localhost:8080`.

### Use the CLI

Build and install:

```bash
cd cmd/cli
go install
```

_(Make sure `$GOPATH/bin` is in your `PATH`.)_

Authenticate and manage tasks:

```bash
gopher-tasks login --email you@example.com --password secret
gopher-tasks task create "Write README"
gopher-tasks task list --status=pending
gopher-tasks task complete 1234-abcd
```

---

## Project Structure

```
gopher-tasks/
├── cmd/
│   ├── server/           # API entrypoint
│   │   └── main.go
│   └── cli/              # CLI entrypoint (cobra)
│       └── main.go
├── internal/
│   ├── domain/           # Entities & interfaces (User, Task, Repositories)
│   ├── usecase/          # Business logic (Register, Login, CreateTask…)
│   ├── delivery/         # External layers
│   │   ├── http/         # HTTP handlers & router
│   │   └── cli/          # cobra commands
│   └── infrastructure/   # Concrete impls
│       ├── persistence/  # SQL repos
│       ├── auth/         # JWT, password hashing
│       └── config/       # Viper-based config loader
├── pkg/                  # Reusable utilities (logger, errors)
├── api/                  # OpenAPI/Swagger definitions
├── configs/              # YAML config files
├── scripts/              # DB migrations, seeds
├── test/                 # Integration tests & fixtures
├── go.mod
└── README.md
```

---

## Testing

Run unit tests:

```bash
go test ./internal/usecase
```

Run all tests:

```bash
go test ./...
```

For integration tests (with Testcontainers):

```bash
# ensure Docker is running
go test -tags=integration ./...
```

---

## Contributing

1. Fork the repository  
2. Create a feature branch (`git checkout -b feat/your-feature`)  
3. Commit your changes (`git commit -m "feat: add …"`)  
4. Push to your branch (`git push origin feat/your-feature`)  
5. Open a Pull Request  

Please follow the [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) guidelines.

---

## License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.
