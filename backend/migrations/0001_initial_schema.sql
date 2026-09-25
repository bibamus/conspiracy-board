-- Connection types
CREATE TABLE connection_types (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    color TEXT NOT NULL DEFAULT '#666',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- People/Nodes
CREATE TABLE people (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT UNIQUE NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Connections/Edges. At most one connection per direction (A->B), regardless of type.
CREATE TABLE connections (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    from_person_id INTEGER NOT NULL,
    to_person_id INTEGER NOT NULL,
    type_id INTEGER NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (from_person_id, to_person_id),
    FOREIGN KEY (from_person_id) REFERENCES people(id) ON DELETE CASCADE,
    FOREIGN KEY (to_person_id) REFERENCES people(id) ON DELETE CASCADE,
    FOREIGN KEY (type_id) REFERENCES connection_types(id) ON DELETE RESTRICT
);

-- The UNIQUE constraint above already indexes from_person_id as its prefix.
CREATE INDEX idx_connections_to ON connections(to_person_id);
CREATE INDEX idx_connections_type ON connections(type_id);

-- Keep updated_at current on every UPDATE. The WHEN guard lets callers set
-- updated_at explicitly and stops the trigger's own UPDATE from re-firing it.
CREATE TRIGGER trg_people_updated_at
AFTER UPDATE ON people
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE people SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER trg_connections_updated_at
AFTER UPDATE ON connections
FOR EACH ROW WHEN NEW.updated_at IS OLD.updated_at
BEGIN
    UPDATE connections SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;
