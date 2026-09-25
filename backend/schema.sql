-- Single source of truth for the database schema.
-- Embedded into the backend binary and applied on startup (see database.go),
-- so every statement must be idempotent.

-- Connection types
CREATE TABLE IF NOT EXISTS connection_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT,
    color TEXT DEFAULT '#666',
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
    FOREIGN KEY (type_id) REFERENCES connection_types(id) ON DELETE RESTRICT
);

-- At most one connection per direction (A->B), regardless of type.
-- An index (not a table constraint) so it also applies to existing databases.
CREATE UNIQUE INDEX IF NOT EXISTS idx_connections_pair ON connections(from_person_id, to_person_id);

-- Indexes for performance
CREATE INDEX IF NOT EXISTS idx_connections_from ON connections(from_person_id);
CREATE INDEX IF NOT EXISTS idx_connections_to ON connections(to_person_id);
CREATE INDEX IF NOT EXISTS idx_connections_type ON connections(type_id);

-- Keep updated_at current on every UPDATE. The WHEN guard lets callers set
-- updated_at explicitly and stops the trigger's own UPDATE from re-firing it.
CREATE TRIGGER IF NOT EXISTS trg_people_updated_at
AFTER UPDATE ON people
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE people SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS trg_connections_updated_at
AFTER UPDATE ON connections
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE connections SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
