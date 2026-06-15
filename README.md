# UserAgeAPI

UserAgeAPI is a Go REST API for managing users and returning each user's age dynamically from their date of birth. It exposes CRUD endpoints for users, stores user records in PostgreSQL, uses SQLC for type-safe database access, and follows a handler-service-repository structure.

## Features

- Create a user with a name and date of birth.
- Get a user by ID.
- Update a user's name and date of birth.
- Delete a user.
- List users with pagination.
- Calculate age dynamically from the stored date of birth.
- Validate request payloads before they reach the service layer.
- Persist users in PostgreSQL.
- Generate type-safe database code with SQLC.
- Run locally or with Docker Compose.
- Log service events and internal errors with Zap.

## Tech Stack

- Go `1.26.0`
- Fiber `v2`
- PostgreSQL
- pgx `v5`
- SQLC
- golang-migrate migration files
- go-playground/validator
- Zap logger
- godotenv
- Docker
- Docker Compose

## Project Structure

```text
.
|-- cmd/
|   `-- server/
|       `-- main.go                 # Application entry point and dependency wiring
|-- config/
|   `-- database.go                 # PostgreSQL connection setup
|-- db/
|   |-- migrations/                 # Database migration files
|   `-- sqlc/
|       |-- query.sql               # SQLC query definitions
|       |-- schema.sql              # SQLC schema input
|       `-- generated/              # Generated SQLC Go code
|-- internal/
|   |-- errors/                     # Domain error values
|   |-- handler/                    # HTTP handlers and request validation
|   |-- logger/                     # Zap logger setup
|   |-- models/                     # Request and response models
|   |-- repository/                 # Database access wrapper around SQLC
|   |-- routes/                     # HTTP route registration
|   `-- service/                    # Business logic and age calculation
|-- Dockerfile                      # Multi-stage application image build
|-- docker-compose.yml              # PostgreSQL, migration, and app services
|-- go.mod
|-- go.sum
`-- sqlc.yaml                       # SQLC configuration
```

## Prerequisites

- Go `1.26.0` or a compatible Go version
- PostgreSQL
- SQLC, for regenerating database code
- golang-migrate CLI, for running migrations locally
- Docker and Docker Compose, if running with containers

## Environment Variables

| Variable | Required | Description | Example |
| --- | --- | --- | --- |
| `DATABASE_URL` | Yes | PostgreSQL connection string used by `config.ConnectDB`. | `postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable` |


## Setup and Run

### 1. Clone the repository

```bash
git clone https://github.com/aditya2k413/user-management-api
cd user-management-api
```

### 2. Create a `.env` file

```env
DATABASE_URL=postgres://postgres:postgres@localhost:5432/userdb?sslmode=disable
```

### 3. Install dependencies

```bash
go mod download
```
### 4. Create the database

```sql
CREATE DATABASE userdb;
```

### 5. Run migrations

```bash
migrate -path db/migrations -database "$DATABASE_URL" up
```

PowerShell:

```powershell
migrate -path db/migrations -database $env:DATABASE_URL up
```



### 6. Start the application

```bash
go run ./cmd/server
```

The API will be available at:

```text
http://localhost:3000
```

### 7. Run tests

```bash
go test ./...
```

## Docker Setup

### Build the image

```bash
docker build -t user-age-api .
```

### Run with Docker Compose

```bash
docker compose up --build
```

The API will be available at:

```text
http://localhost:3000
```

## API Documentation

Base URL:

```text
http://localhost:3000
```


### Create User

Method: `POST`

Route: `/users`

Request:

```json
{
  "name": "Alice Johnson",
  "dob": "2000-06-15"
}
```

Response: `201 Created`

```json
{
  "id": 1,
  "name": "Alice Johnson",
  "dob": "2000-06-15"
}
```

### Get User By ID

Method: `GET`

Route: `/users/:id`

Example:

```bash
curl http://localhost:3000/users/1
```

Response: `200 OK`

```json
{
  "id": 1,
  "name": "Alice Johnson",
  "dob": "2000-06-15",
  "age": 26
}
```

### Update User

Method: `PUT`

Route: `/users/:id`

Request:

```json
{
  "name": "Alice Smith",
  "dob": "1999-12-20"
}
```

Response: `200 OK`

```json
{
  "id": 1,
  "name": "Alice Smith",
  "dob": "1999-12-20"
}
```

### Delete User

Method: `DELETE`

Route: `/users/:id`

Example:

```bash
curl -X DELETE http://localhost:3000/users/1
```

Response: `204 No Content`

The response body is empty.

### List Users

Method: `GET`

Route: `/users`

Query parameters:

| Parameter | Default | Behavior |
| --- | --- | --- |
| `page` | `1` | Values below `1` are treated as `1`. |
| `limit` | `10` | Values below `1` are treated as `10`; values above `100` are capped at `100`. |

Example:

```bash
curl "http://localhost:3000/users?page=1&limit=10"
```

Response: `200 OK`

```json
[
  {
    "id": 1,
    "name": "Alice Johnson",
    "dob": "2000-06-15",
    "age": 26
  },
  {
    "id": 2,
    "name": "Bob Lee",
    "dob": "1995-03-10",
    "age": 31
  }
]
```

## Validation Rules

Create and update requests use the same validation rules:

| Field | Type | Rules |
| --- | --- | --- |
| `name` | string | Required, minimum length `2`, maximum length `100`. |
| `dob` | string | Required, must match date format `YYYY-MM-DD`. |

Route parameter validation:

| Parameter | Rules |
| --- | --- |
| `id` | Must parse as a base-10 integer that fits into `int32`. |

List query behavior:

| Query Parameter | Rules |
| --- | --- |
| `page` | Defaults to `1`; values less than `1` are normalized to `1`. |
| `limit` | Defaults to `10`; values less than `1` are normalized to `10`; values greater than `100` are capped at `100`. |


## Error Handling

Common error responses:

| Scenario | Status | Response |
| --- | --- | --- |
| Invalid JSON/request body | `400 Bad Request` | `{"error":"invalid request body"}` |
| Validation failure | `400 Bad Request` | `{"error":"validation failed: invalid input format"}` |
| Invalid user ID path parameter | `400 Bad Request` | `{"error":"invalid user id"}` |
| Invalid date parsed by service | `400 Bad Request` | `{"error":"invalid date format, use YYYY-MM-DD"}` |
| User not found | `404 Not Found` | `{"error":"user not found"}` |
| Unexpected error | `500 Internal Server Error` | `{"error":"internal server error"}` |

## Development Notes

- Handler layer: parses HTTP requests, validates payloads, converts route/query parameters, calls the service layer, and returns HTTP responses.
- Service layer: owns business logic, age calculation, date parsing, logging, pagination normalization, and domain error mapping.
- Repository layer: wraps SQLC-generated query methods so the service layer does not depend directly on every generated method call.
- SQLC usage: SQL queries live in `db/sqlc/query.sql`; SQLC generates typed Go methods in `db/sqlc/generated`.
- Dependency injection: `cmd/server/main.go` wires the database pool, SQLC queries, repository, service, handler, and routes explicitly.
- Database schema: the active schema has a single `users` table with `id`, `name`, and `dob`.
- Configuration: only `DATABASE_URL` is environment-driven; the HTTP port is hardcoded.


## Future Improvements

- Add API versioning, such as `/api/v1/users`.
- Add response metadata for paginated list endpoints.
- Add validation to reject future dates of birth.
- Add database constraints or service checks for duplicate users if required by the domain.
- Add integration tests for handlers and repository behavior.
- Add graceful shutdown for the Fiber server and database pool.
- Add structured request logging middleware and request IDs.
- Add OpenAPI/Swagger documentation.
- Make the server port configurable with an environment variable.
- Make Docker Compose wait for migrations before starting the app.
