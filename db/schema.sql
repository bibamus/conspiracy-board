-- Connection types
CREATE TABLE IF NOT EXISTS connection_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- People/Nodes
CREATE TABLE IF NOT EXISTS people (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Connections/Edges
CREATE TABLE IF NOT EXISTS connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_person_id INTEGER NOT NULL,
    to_person_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    description TEXT,
    weight REAL DEFAULT 1.0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (from_person_id) REFERENCES people(id) ON DELETE CASCADE,
    FOREIGN KEY (to_person_id) REFERENCES people(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES connection_types(id) ON DELETE RESTRICT,
    UNIQUE(from_person_id, to_person_id, type_id)
);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_connections_from ON connections(from_person_id);
CREATE INDEX IF NOT EXISTS idx_connections_to ON connections(to_person_id);
CREATE INDEX IF NOT EXISTS idx_connections_type ON connections(type_id);
