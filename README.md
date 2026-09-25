# Conspiracy Board

A directed graph visualization system for mapping relationships between people. Build conspiracy theories, social networks, or any other relational data with typed connections.

## Quick Start

### Start Backend
```bash
cd backend
go run .
# Server runs on http://localhost:8000 (database: ./conspiracy-board.db, override with DB_PATH)
```

### Start Frontend
```bash
cd frontend
npm install
npm run dev
# App runs on http://localhost:3000
```

Requires Go 1.27+ and Node.js 20.19+ or 22.12+. For a single-container deployment see
[DOCKER.md](DOCKER.md).

## Features

- **Directed Graph**: Create people (nodes) and connections (edges) between them
- **Typed Connections**: Define custom connection types (e.g., "knows", "works-with", "investigated-by")
- **REST API**: Create, list and delete people and connections; full CRUD for connection types
- **SQLite Database**: Persistent storage with modernc.org/sqlite driver
- **Interactive Visualization**: Modern React UI with vis-network graph visualization
- **Automatic Updates**: Graph and forms reload after every change

## Project Structure

```
conspiracy-board/
├── backend/              # Go REST API backend
│   ├── main.go           # Router setup and server start
│   ├── models.go
│   ├── database.go
│   ├── migrate.go        # Migration runner
│   ├── migrations/       # Numbered SQL migrations (embedded, applied on startup)
│   ├── handlers.go
│   ├── *_test.go         # API and migration tests
│   └── go.mod
├── frontend/             # React Vite frontend
│   ├── src/
│   │   ├── components/   # React components
│   │   ├── api.js       # API client
│   │   ├── App.jsx
│   │   └── main.jsx
│   ├── package.json
│   ├── vite.config.js
│   └── index.html
└── README.md
```

## Database Migrations

The schema is defined by the SQL files in `backend/migrations/`, which are embedded in the
binary. On startup the backend applies any migration not yet listed in the `schema_migrations`
table, in version order, inside a single transaction (all pending migrations succeed or none do).

To change the schema, add a new file with the next number, e.g.
`backend/migrations/0002_add_person_notes.sql`. Rules:

- File names must match `NNNN_name.sql` (lowercase name), with versions contiguous from `0001`.
- Never edit or delete a migration that has already been released — add a new one instead.
- The backend refuses to start if the database is at a newer version than the binary knows.
- `PRAGMA foreign_keys` has no effect inside a transaction, so migrations that rebuild tables
  must not rely on toggling it.

Run the tests with `cd backend && go test ./...`.

## Backend API

All endpoints are served by the backend on port 8000 (behind nginx on port 8080 in Docker).

- Request and response bodies are JSON.
- List endpoints always return an array (`[]` when empty), never `null`.
- Errors are returned as `{"error": "<message>"}` with status `400` (invalid input, e.g.
  `"name is required"`), `404` (not found), `409` (duplicate name/connection, or a connection
  type still in use) or `500`.

### Health Check
```
GET    /health
```

### Connection Types
```
POST   /api/connection-types      # Create type {name, description?, color?}
GET    /api/connection-types      # List types
PUT    /api/connection-types/:id  # Update type {name, description?, color?}
DELETE /api/connection-types/:id  # Delete type (409 while connections use it)
```

### People
```
POST   /api/people              # Create person {name, description?}
GET    /api/people              # List people
GET    /api/people/:id          # Get person
DELETE /api/people/:id          # Delete person and their connections
```

### Connections
```
POST   /api/connections             # Create {from_person_id, to_person_id, type_id, description?}
GET    /api/connections             # List connections
GET    /api/people/:id/connections  # Get person's connections
DELETE /api/connections/:id         # Delete connection
```

### Graph
```
GET    /api/graph               # {nodes: [{person, connections}], types}
```

## Frontend

See `frontend/README.md` for detailed frontend documentation.

### Features
- Create, edit and delete connection types
- Add and delete people
- Create and delete connections between people
- Interactive graph visualization
- Automatic refresh after changes
