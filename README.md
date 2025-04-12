```sh
go run .
```

run the seed command

```sh
go run . seed
```

# Webmane Go Music API

A Go application serving a GraphQL API for managing a music library stored in a PostgreSQL database. It also includes command-line interface (CLI) capabilities.

## Overview

This project provides a backend service with the following main components:

1.  **GraphQL API:** Exposes queries and mutations for interacting with music data.
    - Fetch lists of songs with pagination and search (`music` query).
    - Create or update song entries (`upsertSong` mutation).
    - (Planned: Fetch single song by ID - `song` query is defined but not implemented).
2.  **PostgreSQL Database:** Uses `pgx` to connect to and interact with a PostgreSQL database, storing song information in a `music` table.
3.  **Web Server:** Uses `chi` for routing.
    - Serves the GraphQL API at `/query`.
    - Provides a GraphQL Playground UI at `/` for easy testing.
    - Includes WebSocket support for potential GraphQL subscriptions (origin check currently restricted to `webmane.net`).
    - Configured with CORS for frontend development (allowing `localhost:5173`).
    - Includes a separate REST endpoint at `/music` (handled by the `music` package).
4.  **Command-Line Interface (CLI):** Uses `cobra` for defining CLI commands (likely for tasks like database management, seeding, etc.). Commands are defined in the `cmd/` directory.
5.  **Dependencies:** Key dependencies include `gqlgen`, `pgx`, `chi`, `cobra`, `cors`, `websocket`, and `ffmpeg-go` (suggesting potential audio/video processing capabilities).

## Technology Stack

- **Language:** Go (version 1.21+)
- **API:** GraphQL (`gqlgen`)
- **Database:** PostgreSQL (`pgx/v4`)
- **HTTP Router:** `chi`
- **WebSockets:** `gorilla/websocket`
- **CLI:** `cobra`
- **CORS:** `rs/cors`

## Getting Started

1.  **Prerequisites:**
    - Go (>= 1.21) installed.
    - A running PostgreSQL database instance.
    - (Potentially) FFmpeg installed if its features are used.
2.  **Clone:** Clone the repository.
3.  **Dependencies:** Install Go modules:
    ```bash
    go mod tidy
    ```
4.  **Database Setup:**
    - Ensure your PostgreSQL server is running.
    - Configure the database connection string. The application likely expects environment variables or a configuration file used by `db.ConnectToDb()` (check the `db/` package or environment setup).
    - Apply database schema/migrations. Check the `cmd/` directory for potential migration commands or apply the schema manually if needed. The primary table seems to be `music`.
5.  **Run the Server:**
    ```bash
    go run server.go
    ```
    The server will start, typically on port `8080`. You can access the GraphQL Playground at `http://localhost:8080/`.

## Development Commands

- **Run the Web Server:**
  ```bash
  go run server.go
  ```
- **Regenerate GraphQL Code (`gqlgen`):**
  After modifying the GraphQL schema (`graph/schema.graphqls`) or related configuration (`gqlgen.yml`), run this command to update the generated Go code (resolvers, models, etc.):
  ```bash
  go run github.com/99designs/gqlgen generate
  ```
- **Run CLI Commands:**
  The `server.go` entry point also executes the Cobra root command. To run specific CLI commands defined in the `cmd/` package:
  ```bash
  # Replace [command_name] and [args...] with actual commands/arguments
  go run server.go [command_name] [args...]
  ```
  _(Example: If there's a `migrate` command, it might be `go run server.go migrate up`)_
- **Build the Binary:**
  ```bash
  go build -o webmane_go server.go
  ```
  Then run the compiled application: `./webmane_go` (or `./webmane_go [command_name] [args...]` for CLI commands).

## API Endpoints

- `/`: GraphQL Playground UI.
- `/query`: GraphQL API endpoint (HTTP POST for queries/mutations, WebSocket for subscriptions).
- `/music`: REST endpoint (GET - see `music/` package for details).

## TODO / Areas to Check

- Implement the `song(id: ID!)` query resolver in `graph/schema.resolvers.go`.
- Review the implementation and purpose of the `/music` REST endpoint in the `music/` package.
- Investigate how the `ffmpeg-go` dependency is used.
- Update the WebSocket `CheckOrigin` function in `server.go` for deployment environments.
- Explore and document the available CLI commands in the `cmd/` directory.

## Overview

This project provides a backend service with the following main components:

1.  **GraphQL API:** Exposes queries and mutations for interacting with music data via `gqlgen`.
    - Fetch lists of songs with pagination and search (`music` query).
    - Create or update song entries (`upsertSong` mutation).
    - (Planned: Fetch single song by ID - `song` query is defined but not implemented).
2.  **Database Interaction (Ent):** Uses the **Ent framework** (`entgo.io/ent`) as an ORM.
    - Defines database schema in Go code (`ent/schema/`).
    - Generates a type-safe Go client for all database operations (CRUD, queries, graph traversals).
    - Handles database connections (`db/` package).
    - Provides schema migration capabilities (currently using `Schema.Create` for development).
3.  **PostgreSQL Database:** Uses `pgx` (underneath Ent) to connect to and interact with a PostgreSQL database. The schema (e.g., `music_ent` table) is managed by Ent.

```bash
go generate ./ent
```

4.  **Web Server:** Uses `chi` for routing.
    - Serves the GraphQL API at `/query`.
    - Provides a GraphQL Playground UI at `/` for easy testing.
    - Includes WebSocket support for potential GraphQL subscriptions (origin check currently restricted).
    - Configured with CORS for frontend development (allowing `localhost:5173`).
    - Includes a separate REST endpoint at `/music` (handled by the `music` package, also using the Ent client).
5.  **Command-Line Interface (CLI):** Uses `cobra` for defining CLI commands (e.g., `seed` command for populating the database). Commands are defined in the `cmd/` directory and receive the Ent client via the GraphQL resolver.
6.  **Dependencies:** Key dependencies include `gqlgen`, `entgo.io/ent`, `pgx`, `chi`, `cobra`, `cors`, `websocket`, and `ffmpeg-go`.

## Technology Stack

- **Language:** Go (version 1.21+)
- **API:** GraphQL (`gqlgen`)
- **ORM / Database:** Ent (`entgo.io/ent`) with PostgreSQL (`pgx/v4`)
- **HTTP Router:** `chi`
- **WebSockets:** `gorilla/websocket`
- **CLI:** `cobra`
- **CORS:** `rs/cors`

## Getting Started

1.  **Prerequisites:**
    - Go (>= 1.21) installed.
    - A running PostgreSQL database instance.
    - (Potentially) FFmpeg installed if its features are used (e.g., by the `seed` command).
2.  **Clone:** Clone the repository.
3.  **Dependencies:** Install Go modules:
    ```bash
    go mod tidy
    ```
4.  **Database Setup:**
    - Ensure your PostgreSQL server is running.
    - Configure the database connection string. The application expects the `DATABASE_URL` environment variable (see `db/db.go`). Example: `export DATABASE_URL="postgres://user:password@host:port/database_name?sslmode=disable"`
    - **Schema Management (Ent):**
      - The database schema is defined as Go code in the `ent/schema/` directory.
      - On startup (`server.go`), the application currently uses `entClient.Schema.Create` to check and potentially update the database schema based on these definitions. **Note:** This is suitable for development but **not recommended for production**. For production environments, use a dedicated migration tool like Atlas (which integrates well with Ent).
      - The generated, type-safe database client code resides in the `ent/` directory.
5.  **Run the Server:**
    ```bash
    go run .
    ```
    The server will start, typically on port `8080`. You can access the GraphQL Playground at `http://localhost:8080/`.

## Development Commands

- **Run the Web Server:**
  ```bash
  go run .
  # or: go run server.go
  ```
