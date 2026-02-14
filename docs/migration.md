# Database Migrations

This document describes how to work with database migrations in the Nimbus project using [goose](https://github.com/pressly/goose).

## Overview

Nimbus uses [pressly/goose](https://github.com/pressly/goose) for managing database schema migrations. Goose provides a simple CLI for versioning and applying database changes using annotated SQL files.

## Installation

### Option 1: Using Docker (Recommended)

The project already includes goose in `docker-compose.yml`. No additional installation is required when using Docker Compose.

### Option 2: Using Make

```bash
make install
```

This installs goose along with other project dependencies.

### Option 3: Install CLI Tool Locally

#### Using Go

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

#### macOS

```bash
brew install goose
```

### Verify Installation

```bash
goose --version
```

## Migration File Structure

Migration files are located in the `migrations/` directory. Each migration is a **single SQL file** with `-- +goose Up` and `-- +goose Down` annotations:

```
{version}_{description}.sql
```

**Example:**
```
00001_create_workflows_table.sql
00002_create_tasks_table.sql
```

**File format:**
```sql
-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS users;
```

## Creating New Migrations

```bash
make migrate-create name=add_user_table
```

This creates a new file: `migrations/00003_add_user_table.sql`

Edit the file and add your `-- +goose Up` and `-- +goose Down` sections.

### Migration File Guidelines

**Up section (`-- +goose Up`):**
- Contains SQL statements to apply the change
- Should be idempotent when possible (use `IF NOT EXISTS`, etc.)

**Down section (`-- +goose Down`):**
- Contains SQL statements to reverse the change
- Should cleanly undo what the Up section does

## Running Migrations

### Using Docker Compose (Recommended)

Migrations run automatically when you start the services:

```bash
# Start all services (database + migrations)
docker-compose up -d

# View migration logs
docker-compose logs migration
```

### Using Make Commands

First, ensure your PostgreSQL database is running:

```bash
docker-compose up -d database
```

Then run migrations:

```bash
# Run all pending migrations
make migrate

# Rollback last migration
make migrate-down

# Check migration status
make migrate-status

# Re-run the latest migration
make migrate-redo

# Create a new migration
make migrate-create name=add_user_table
```

### Using goose CLI Directly

```bash
# Set database connection string
export DB_DSN="postgres://nimbus:nimbus_password@localhost:5432/nimbus_db?sslmode=disable"

# Run all pending migrations
goose -dir migrations postgres "$DB_DSN" up

# Rollback last migration
goose -dir migrations postgres "$DB_DSN" down

# Rollback to specific version
goose -dir migrations postgres "$DB_DSN" down-to 1

# Check migration status
goose -dir migrations postgres "$DB_DSN" status

# Re-run latest migration
goose -dir migrations postgres "$DB_DSN" redo

# Reset all migrations
goose -dir migrations postgres "$DB_DSN" reset
```

## Additional Resources

- [Goose Documentation](https://pressly.github.io/goose/)
- [Goose GitHub Repository](https://github.com/pressly/goose)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)

## Current Migrations

The project currently includes the following migrations:

1. **00001_create_workflows_table** - Creates the workflows table with UUID primary key
2. **00002_create_tasks_table** - Creates the tasks table with status enum and foreign key to workflows
