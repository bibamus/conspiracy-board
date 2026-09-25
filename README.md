# Conspiracy Board

A directed graph visualization system for mapping relationships between people. Build conspiracy theories, social networks, or any other relational data with typed connections.

## Quick Start

### Start Backend
```bash
cd backend
go build -o conspiracy-board.exe
./conspiracy-board.exe
# Server runs on http://localhost:8080
```

### Start Frontend
```bash
cd frontend
npm install
npm run dev
# App runs on http://localhost:3000
```

## Features

- **Directed Graph**: Create people (nodes) and connections (edges) between them
- **Typed Connections**: Define custom connection types (e.g., "knows", "works-with", "investigated-by")
- **REST API**: Full CRUD operations for managing the graph
- **SQLite Database**: Persistent storage with modernc.org/sqlite driver
- **Interactive Visualization**: Modern React UI with vis-network graph visualization
- **Real-time Updates**: Graph updates automatically when data changes

## Project Structure

```
conspiracy-board/
├── backend/              # Go REST API backend
│   ├── main.go
│   ├── models.go
│   ├── database.go
│   ├── migrate.go        # Migration runner
│   ├── migrations/       # Numbered SQL migrations (embedded, applied on startup)
│   ├── handlers.go
│   ├── go.mod
│   ├── conspiracy-board.exe
│   └── conspiracy-board.db
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

Full REST API documentation available at `backend/README.md` or checkout these endpoints:

### Health Check
```
GET http://localhost:8080/health
```

### Connection Types
```
POST   /api/connection-types    # Create type
GET    /api/connection-types    # List types
```

### People
```
POST   /api/people              # Create person
GET    /api/people              # List people
GET    /api/people/:id          # Get person
```

### Connections
```
POST   /api/connections         # Create connection
GET    /api/connections         # List connections
GET    /api/people/:id/connections  # Get person's connections
```

### Graph
```
GET    /api/graph               # Get full graph structure
```

## Frontend

See `frontend/README.md` for detailed frontend documentation.

### Features
- Add connection types
- Add people to the graph
- Create connections between people
- Interactive graph visualization
- Real-time updates
- Responsive design
