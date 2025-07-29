# ☄️ Comet - Blazing-Fast, Schema-First ORM for Go

<div align="center">

![☄️ Comet Logo](https://images.unsplash.com/photo-1419242902214-272b3f66ee7a?w=800&h=400&fit=crop&crop=center&auto=format&q=80)

*A lightweight, Prisma-inspired ORM for Go that prioritizes simplicity, speed, and type-safety.*

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![License](https://img.shields.io/badge/License-MIT-blue?style=for-the-badge)](LICENSE)
[![Build Status](https://img.shields.io/badge/Build-Passing-brightgreen?style=for-the-badge)](https://github.com/nitrix4ly/comet)

</div>

---

## Features

- **Schema-first development** with `.cmt` files
- **Fast, type-safe query building** with fluent API
- **Built-in CLI** for codegen, migrations, and seeding
- **Multiple database backends** (PostgreSQL, MySQL, SQLite)
- **Minimal dependencies**, perfect for microservices
- **Goroutine-optimized** for async patterns

<div align="center">

![Database Architecture](https://images.unsplash.com/photo-1558494949-ef010cbdcc31?w=600&h=300&fit=crop&crop=center&auto=format&q=80)

</div>

## Quick Start

```bash
# Install ☄️ Comet CLI
go install github.com/nitrix4ly/comet/cli@latest

# Initialize project
mkdir myapp && cd myapp
go mod init myapp

# Create schema
cat > schema/schema.cmt << EOF
model User {
  id        Int      @id @auto
  email     String   @unique
  name      String
  createdAt DateTime @default(now())
}
EOF

# Generate models
comet gen

# Run migrations
comet migrate
```

## Documentation

See [docs/README.md](docs/README.md) for complete setup and usage instructions.

<div align="center">

![Code Example](https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=600&h=300&fit=crop&crop=center&auto=format&q=80)

</div>

## Example Usage

```go
import "myapp/models"

// Find users
users, err := models.User.Find().
    Where("email", "=", "hello@comet.dev").
    All(ctx)

// Create user
user := &models.User{
    Email: "john@example.com",
    Name:  "John Doe",
}
err = user.Save(ctx)
```

## Architecture

<div align="center">

![Architecture Diagram](https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600&h=300&fit=crop&crop=center&auto=format&q=80)

</div>

- **cli/** - Command-line interface
- **core/** - Core types and query builder
- **gen/** - Code generation from schemas
- **drivers/** - Database drivers
- **schema/** - Schema definitions
- **testz/** - Example applications

## Contributing

1. Fork the repository
2. Create your feature branch
3. Add tests for new functionality
4. Submit a pull request

<div align="center">

![Contributing](https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

## License

MIT License - see LICENSE file for details.

---

Built with ❤️ by [Nitrix](https://github.com/nitrix4ly)
