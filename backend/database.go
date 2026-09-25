package main

import (
	"database/sql"
	"errors"
	"fmt"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase(path string) (*Database, error) {
	// Pragmas are applied by modernc.org/sqlite to every pooled connection.
	dsn := fmt.Sprintf("file:%s?mode=rwc&_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)", path)
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
	    FOREIGN KEY (type_id) REFERENCES connection_types(id) ON DELETE RESTRICT
	);

	-- At most one connection per direction (A->B), regardless of type.
	-- An index (not a table constraint) so it also applies to existing databases.
	CREATE UNIQUE INDEX IF NOT EXISTS idx_connections_pair ON connections(from_person_id, to_person_id);
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

// querier is satisfied by both *sql.DB and *sql.Tx.
type querier interface {
	Query(query string, args ...any) (*sql.Rows, error)
}

func (db *Database) GetAllConnectionTypes() ([]ConnectionType, error) {
	return queryConnectionTypes(db.conn)
}

func queryConnectionTypes(q querier) ([]ConnectionType, error) {
	rows, err := q.Query("SELECT id, name, description, color, created_at FROM connection_types ORDER BY name")
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
	return queryPeople(db.conn)
}

func queryPeople(q querier) ([]Person, error) {
	rows, err := q.Query("SELECT id, name, description, created_at, updated_at FROM people ORDER BY name")
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

// ErrConnectionExists is returned when a connection from A to B already exists.
var ErrConnectionExists = errors.New("a connection from this person to the other person already exists")

func (db *Database) CreateConnection(fromID, toID, typeID int, description string, weight float64) (*Connection, error) {
	if weight <= 0 {
		weight = 1.0
	}

	result, err := db.conn.Exec(
		"INSERT INTO connections (from_person_id, to_person_id, type_id, description, weight) VALUES (?, ?, ?, ?, ?)",
		fromID, toID, typeID, description, weight,
	)
	if err != nil {
		var sqliteErr *sqlite.Error
		if errors.As(err, &sqliteErr) && sqliteErr.Code() == sqlite3.SQLITE_CONSTRAINT_UNIQUE {
			return nil, ErrConnectionExists
		}
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

const connectionColumns = "SELECT id, from_person_id, to_person_id, type_id, description, weight, created_at, updated_at FROM connections"

func (db *Database) GetConnectionsByPerson(personID int) ([]Connection, error) {
	return queryConnections(db.conn,
		connectionColumns+" WHERE from_person_id = ? OR to_person_id = ? ORDER BY created_at DESC",
		personID, personID,
	)
}

func (db *Database) GetAllConnections() ([]Connection, error) {
	return queryConnections(db.conn, connectionColumns+" ORDER BY created_at DESC")
}

func queryConnections(q querier, query string, args ...any) ([]Connection, error) {
	rows, err := q.Query(query, args...)
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

// GetGraph returns the full graph using a fixed number of queries inside a
// single transaction, so people, types and connections form a consistent snapshot.
func (db *Database) GetGraph() (*Graph, error) {
	tx, err := db.conn.Begin()
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	people, err := queryPeople(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch people: %w", err)
	}

	types, err := queryConnectionTypes(tx)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch connection types: %w", err)
	}

	connections, err := queryConnections(tx, connectionColumns+" ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch connections: %w", err)
	}

	graph := &Graph{
		Nodes: make([]*GraphNode, 0, len(people)),
		Types: types,
	}

	byPerson := make(map[int]*GraphNode, len(people))
	for _, p := range people {
		node := &GraphNode{Person: p}
		byPerson[p.ID] = node
		graph.Nodes = append(graph.Nodes, node)
	}

	// Each connection is listed under both endpoints (once for self-loops),
	// matching the previous per-person query semantics.
	for _, c := range connections {
		if node, ok := byPerson[c.FromID]; ok {
			node.Connections = append(node.Connections, c)
		}
		if c.ToID != c.FromID {
			if node, ok := byPerson[c.ToID]; ok {
				node.Connections = append(node.Connections, c)
			}
		}
	}

	return graph, nil
}
