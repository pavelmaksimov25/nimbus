# Database Migrations

This document describes how to work with database migrations in the Nimbus project using the [golang-migrate](https://github.com/golang-migrate/migrate) tool.

## Overview

Nimbus uses [golang-migrate/migrate](https://github.com/golang-migrate/migrate/blob/v4.19.0/cmd/migrate/README.md) for managing database schema migrations. This tool provides a simple and reliable way to version and apply database changes.

## Installation

### Option 1: Using Docker (Recommended)

The project already includes the migrate tool in `docker-compose.yml`. No additional installation is required when using Docker Compose.

### Option 2: Install CLI Tool Locally

#### macOS

```bash
brew install golang-migrate
```

#### Linux

```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/
```

#### Windows

Download the binary from the [releases page](https://github.com/golang-migrate/migrate/releases/tag/v4.19.0) and add it to your PATH.


### Verify Installation

```bash
migrate -version
```

## Migration File Structure

Migration files are located in the `migrations/` directory and follow this naming convention:

```
{version}_{description}.up.sql    # Applied when migrating up
{version}_{description}.down.sql  # Applied when rolling back
```

**Example:**
```
000001_create_workflows_table.up.sql
000001_create_workflows_table.down.sql
000002_create_tasks_table.up.sql
000002_create_tasks_table.down.sql
```

## Creating New Migrations

### Using Docker Compose

```bash
# Create a new migration with sequential numbering
docker-compose run --rm migration create -ext sql -dir /migrations -seq <migration_name>
```

**Example:**
```bash
docker-compose run --rm migration create -ext sql -dir /migrations -seq add_user_table
```

This will create two files:
- `migrations/000003_add_user_table.up.sql`
- `migrations/000003_add_user_table.down.sql`

### Using Local CLI Tool

```bash
# Create a new migration
migrate create -ext sql -dir migrations -seq <migration_name>
```

**Example:**
```bash
migrate create -ext sql -dir migrations -seq add_user_table
```

### Migration File Guidelines

**UP Migration (`*.up.sql`):**
- Contains SQL statements to apply the change
- Should be idempotent when possible (use `IF NOT EXISTS`, etc.)
- Example:

```sql
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);
```

**DOWN Migration (`*.down.sql`):**
- Contains SQL statements to reverse the change
- Should cleanly undo what the UP migration does
- Example:

```sql
DROP TABLE IF EXISTS users;
```

## Running Migrations

### Using Docker Compose (Recommended)

The migrations run automatically when you start the services:

```bash
# Start all services (database + migrations)
docker-compose up -d

# View migration logs
docker-compose logs migration
```

### Using Local CLI Tool

First, ensure your PostgreSQL database is running. You can use the database from docker-compose:

```bash
# Start only the database
docker-compose up -d database
```

Then run migrations using the local tool:

```bash
# Set database connection string
export DATABASE_URL="postgres://user:password@localhost:5432/db_name?sslmode=disable"

# Run all pending migrations
migrate -path migrations -database "$DATABASE_URL" up 2

# Rollback last migration
migrate -path migrations -database "$DATABASE_URL" down 1

# Rollback all migrations
migrate -path migrations -database "$DATABASE_URL" down

# Migrate to specific version
migrate -path migrations -database "$DATABASE_URL" goto 2

# Check current migration version
migrate -path migrations -database "$DATABASE_URL" version

# Force set version (use with caution - only when migration state is corrupted)
migrate -path migrations -database "$DATABASE_URL" force 1
```

## Additional Resources

- [Official golang-migrate Documentation](https://github.com/golang-migrate/migrate/blob/v4.19.0/cmd/migrate/README.md)
- [golang-migrate GitHub Repository](https://github.com/golang-migrate/migrate)
- [PostgreSQL Documentation](https://www.postgresql.org/docs/)
- [Database Migration Best Practices](https://github.com/golang-migrate/migrate/blob/master/MIGRATIONS.md)

## Current Migrations

The project currently includes the following migrations:

1. **000001_create_workflows_table** - Creates the workflows table with UUID primary key
2. **000002_create_tasks_table** - Creates the tasks table with status enum and foreign key to workflows
