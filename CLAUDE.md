# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands
- Run server: `go run main.go serve`
- Run tests: `go test ./...`
- Run specific test: `go test ./auth/pwdless -run TestAuthResource_login`
- Lint code: `go vet ./...`
- Generate API docs: `go run main.go gendoc`
- Run migrations: `go run main.go migrate`

## Code Style Guidelines
- File organization: Package-based structure with resource separation
- Imports: Standard library first, third-party second, project imports last
- Naming: CamelCase for exported items, package names match directories
- Error handling: Custom error types in separate error.go files per package
- Testing: Table-driven tests with descriptive names
- Interfaces: Define before implementation, use mocks for testing
- Documentation: Godoc style comments for exported functions and types
- Dependencies: Use Chi router, Cobra CLI, JWT auth, Bun ORM with PostgreSQL