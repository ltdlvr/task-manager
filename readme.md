# Task Manager API

## Overview

Task Manager API is a REST service for managing kanban-style boards, columns and tasks with JWT authentication.

## Features

* gap-based ordering for columns and tasks positions with automatic rebalancing
* access control with owner/member roles
* JWT authentication
* smoke test script (requires curl and jq)

## Tech Stack

* Go 1.26.1
* Fiber v3
* PostgreSQL
* pgx
* goose
* JWT
* Docker Compose

## Getting Started

1. Clone the repository:

```bash
git clone https://github.com/ltdlvr/task-manager.git
cd task-manager
```

2. Create a local environment file:

```bash
cp .env.dev.example .env.dev
```

3. Fill in .env.dev

4. Start PostgreSQL:

```bash
make db-up
```

5. Apply migrations:

```bash
make db-migrate
```

6. Run the REST server:

```bash
make run-rest
```

The API will be available at:

```text
http://localhost:6969/api/v1
```

7. Run the API smoke test:

```bash
make test-api
```

## API Endpoints

Public endpoints:

| Method | Path | Description |
| --- | --- | --- |
| `GET` | `/api/v1/healthcheck` | Health check |
| `POST` | `/api/v1/register` | Register a user |
| `POST` | `/api/v1/login` | Log in and receive JWT |

Private endpoints require the `Authorization: Bearer <token>` header.

| Method | Path | Description |
| --- | --- | --- |
| `POST` | `/api/v1/boards` | Create a board |
| `GET` | `/api/v1/boards/:id` | Get a board by ID |
| `DELETE` | `/api/v1/boards/:id` | Delete a board by ID |
| `POST` | `/api/v1/boards/:boardId/columns` | Create a column in a board |
| `GET` | `/api/v1/boards/:boardId/columns` | Get all board columns |
| `DELETE` | `/api/v1/columns/:id` | Delete a column by ID |
| `PATCH` | `/api/v1/columns/:id/move` | Move a column to another position |
| `POST` | `/api/v1/columns/:columnId/tasks` | Create a task in a column |
| `GET` | `/api/v1/columns/:columnId/tasks` | Get all column tasks |
| `DELETE` | `/api/v1/tasks/:id` | Delete a task by ID |
| `PATCH` | `/api/v1/tasks/:id/move` | Move a task to another column or position |

## Project Structure

```text
cmd/
  db-migrate/        Database migration entrypoint
  rest/              REST API entrypoint
internal/
  config/            Environment configuration
  core/              Domain models, services and ports
  infra/             PostgreSQL client and repositories
  tool/              Shared infrastructure helpers
  transport/rest/    Fiber handlers and middleware
migrations/          goose SQL migrations
scripts/             API smoke test scripts
```
