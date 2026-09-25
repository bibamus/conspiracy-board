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
- **Weight Support**: Connections can have weights (0-1) to indicate relationship strength
- **Interactive Visualization**: Modern React UI with vis-network graph visualization
- **Real-time Updates**: Graph updates automatically when data changes

## Project Structure

```
conspiracy-board/
├── backend/              # Go REST API backend
│   ├── main.go
│   ├── models.go
│   ├── database.go
│   ├── handlers.go
│   ├── schema.sql        # SQL schema (embedded and applied on startup)
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
- Create connections between people with weights
- Interactive graph visualization
- Real-time updates
- Responsive design
