# MovieLand

A Go application for managing movie data with PostgreSQL database support.

## Features

- Database abstraction layer using pgx
- Database migrations with Goose
- Environment-based configuration
- Separate development and test databases
- Genre management system

## Prerequisites

- Go 1.25.1 or later
- PostgreSQL
- Docker & Docker Compose (optional)
- Goose migration tool

## Setup

### 1. Environment Variables

Create the required environment variables:

```bash
export DATABASE_URL=postgres://postgres:secret@localhost:5432/movies
export DATABASE_TEST_URL=postgres://postgres:secret@localhost:5432/movies_test
```

### 2. Database Setup

Using Docker Compose (recommended):

```bash
docker-compose up -d
```

Create the databases:

```sql
CREATE DATABASE movies;
CREATE DATABASE movies_test;
```

### 3. Run Migrations

```bash
# Development database
make migrate-up

# Test database
make migrate-up-test
```

## Development

### Available Commands

```bash
# Build the application
make build

# Run the application
make run

# Run tests
make test

# Run benchmarks
make bench

# Format code and tidy modules
make tidy

# Create a new migration
make create-migration name=your_migration_name

# Migrate up/down
make migrate-up
make migrate-down

# Test database migrations
make migrate-up-test
make migrate-down-test

# Show all available commands
make help
```

### Project Structure

## Dependencies

- [pgx/v5](https://github.com/jackc/pgx) - PostgreSQL driver and toolkit
- [env/v11](https://github.com/caarlos0/env) - Environment variable parsing
- [Goose](https://github.com/pressly/goose) - Database migration tool

## License

This project is licensed under the terms specified in the LICENSE file.

