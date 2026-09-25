package main

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared&mode=rwc&journal=wal", path)
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(5)

	db := &Database{conn: conn}
	if err := db.initSchema(); err != nil {
		return nil, err
	}

	return db, nil
}

func (db *Database) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS connection_types (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT UNIQUE NOT NULL,
	    description TEXT,
	    color TEXT DEFAULT '#666',
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS people (
	    id INTEGER PRIMARY KEY AUTOINCREMENT,
	    name TEXT UNIQUE NOT NULL,
	    description TEXT,
	    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
	    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

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

	CREATE INDEX IF NOT EXISTS idx_connections_from ON connections(from_person_id);
	CREATE INDEX IF NOT EXISTS idx_connections_to ON connections(to_person_id);
	CREATE INDEX IF NOT EXISTS idx_connections_type ON connections(type_id);
	`

	_, err := db.conn.Exec(schema)
	if err != nil {
		return err
	}

	// Migrate: Add color column if it doesn't exist
	db.conn.Exec("ALTER TABLE connection_types ADD COLUMN color TEXT DEFAULT '#666'")
	return nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

// ConnectionType operations
func (db *Database) CreateConnectionType(name, description, color string) (*ConnectionType, error) {
	result, err := db.conn.Exec(
		"INSERT INTO connection_types (name, description, color) VALUES (?, ?, ?)",
		name, description, color,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetConnectionType(int(id))
}

func (db *Database) GetConnectionType(id int) (*ConnectionType, error) {
	row := db.conn.QueryRow(
		"SELECT id, name, description, color, created_at FROM connection_types WHERE id = ?",
		id,
	)

	var ct ConnectionType
	err := row.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.Color, &ct.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &ct, nil
}

func (db *Database) GetAllConnectionTypes() ([]ConnectionType, error) {
	rows, err := db.conn.Query("SELECT id, name, description, color, created_at FROM connection_types ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var types []ConnectionType
	for rows.Next() {
		var ct ConnectionType
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Description, &ct.Color, &ct.CreatedAt); err != nil {
			return nil, err
		}
		types = append(types, ct)
	}

	return types, rows.Err()
}

func (db *Database) UpdateConnectionType(id int, name, description, color string) (*ConnectionType, error) {
	_, err := db.conn.Exec(
		"UPDATE connection_types SET name = ?, description = ?, color = ? WHERE id = ?",
		name, description, color, id,
	)
	if err != nil {
		return nil, err
	}

	return db.GetConnectionType(id)
}

func (db *Database) DeleteConnectionType(id int) error {
	_, err := db.conn.Exec("DELETE FROM connection_types WHERE id = ?", id)
	return err
}

// Person operations
func (db *Database) CreatePerson(name, description string) (*Person, error) {
	result, err := db.conn.Exec(
		"INSERT INTO people (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetPerson(int(id))
}

func (db *Database) GetPerson(id int) (*Person, error) {
	row := db.conn.QueryRow(
		"SELECT id, name, description, created_at, updated_at FROM people WHERE id = ?",
		id,
	)

	var p Person
	err := row.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (db *Database) GetAllPeople() ([]Person, error) {
	rows, err := db.conn.Query("SELECT id, name, description, created_at, updated_at FROM people ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []Person
	for rows.Next() {
		var p Person
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		people = append(people, p)
	}

	return people, rows.Err()
}

func (db *Database) DeletePerson(id int) error {
	_, err := db.conn.Exec("DELETE FROM people WHERE id = ?", id)
	return err
}

// Connection operations
func (db *Database) CreateConnection(fromID, toID, typeID int, description string, weight float64) (*Connection, error) {
	if weight <= 0 {
		weight = 1.0
	}

	result, err := db.conn.Exec(
		"INSERT INTO connections (from_person_id, to_person_id, type_id, description, weight) VALUES (?, ?, ?, ?, ?)",
		fromID, toID, typeID, description, weight,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetConnection(int(id))
}

func (db *Database) GetConnection(id int) (*Connection, error) {
	row := db.conn.QueryRow(
		"SELECT id, from_person_id, to_person_id, type_id, description, weight, created_at, updated_at FROM connections WHERE id = ?",
		id,
	)

	var c Connection
	err := row.Scan(&c.ID, &c.FromID, &c.ToID, &c.TypeID, &c.Description, &c.Weight, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func (db *Database) GetConnectionsByPerson(personID int) ([]Connection, error) {
	rows, err := db.conn.Query(
		"SELECT id, from_person_id, to_person_id, type_id, description, weight, created_at, updated_at FROM connections WHERE from_person_id = ? OR to_person_id = ? ORDER BY created_at DESC",
		personID, personID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []Connection
	for rows.Next() {
		var c Connection
		if err := rows.Scan(&c.ID, &c.FromID, &c.ToID, &c.TypeID, &c.Description, &c.Weight, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}

	return connections, rows.Err()
}

func (db *Database) GetAllConnections() ([]Connection, error) {
	rows, err := db.conn.Query(
		"SELECT id, from_person_id, to_person_id, type_id, description, weight, created_at, updated_at FROM connections ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var connections []Connection
	for rows.Next() {
		var c Connection
		if err := rows.Scan(&c.ID, &c.FromID, &c.ToID, &c.TypeID, &c.Description, &c.Weight, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}

	return connections, rows.Err()
}

func (db *Database) DeleteConnection(id int) error {
	_, err := db.conn.Exec("DELETE FROM connections WHERE id = ?", id)
	return err
}

// GetGraph returns the full graph
func (db *Database) GetGraph() (*Graph, error) {
	people, err := db.GetAllPeople()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch people: %w", err)
	}

	types, err := db.GetAllConnectionTypes()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch connection types: %w", err)
	}

	graph := &Graph{
		Nodes: make([]*GraphNode, 0, len(people)),
		Types: types,
	}

	for _, p := range people {
		connections, err := db.GetConnectionsByPerson(p.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch connections for person %d: %w", p.ID, err)
		}

		node := &GraphNode{
			Person:      p,
			Connections: connections,
		}
		graph.Nodes = append(graph.Nodes, node)
	}

	return graph, nil
}
