package main

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"modernc.org/sqlite"
	sqlite3 "modernc.org/sqlite/lib"
)

//go:embed schema.sql
var schema string

type Database struct {
	conn *sql.DB
}

// Domain errors returned by Database methods. Handlers map these to HTTP
// status codes; any other error is treated as internal.
var (
	ErrNotFound         = errors.New("record not found")
	ErrPersonNameTaken  = errors.New("a person with this name already exists")
	ErrTypeNameTaken    = errors.New("a connection type with this name already exists")
	ErrConnectionExists = errors.New("a connection from this person to the other person already exists")
	ErrInvalidReference = errors.New("the referenced person or connection type does not exist")
	ErrTypeInUse        = errors.New("this connection type is still used by one or more connections")
)

func sqliteCode(err error) int {
	var sqliteErr *sqlite.Error
	if errors.As(err, &sqliteErr) {
		return sqliteErr.Code()
	}
	return 0
}

// isForeignKeyViolation reports foreign key failures. Deferred checks on
// INSERT report SQLITE_CONSTRAINT_FOREIGNKEY, while ON DELETE RESTRICT is
// enforced via an internal trigger and reports SQLITE_CONSTRAINT_TRIGGER.
func isForeignKeyViolation(err error) bool {
	switch sqliteCode(err) {
	case sqlite3.SQLITE_CONSTRAINT_FOREIGNKEY:
		return true
	case sqlite3.SQLITE_CONSTRAINT_TRIGGER:
		return strings.Contains(err.Error(), "FOREIGN KEY constraint failed")
	}
	return false
}

// translate replaces known SQLite errors with domain errors.
func translate(err error, onUnique, onForeignKey error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, sql.ErrNoRows):
		return ErrNotFound
	case onUnique != nil && sqliteCode(err) == sqlite3.SQLITE_CONSTRAINT_UNIQUE:
		return onUnique
	case onForeignKey != nil && isForeignKeyViolation(err):
		return onForeignKey
	}
	return err
}

// execAffecting runs a statement and returns ErrNotFound if no row was affected.
func (db *Database) execAffecting(query string, args ...any) error {
	result, err := db.conn.Exec(query, args...)
	if err != nil {
		return err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrNotFound
	}
	return nil
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
	_, err := db.conn.Exec(schema)
	if err != nil {
		return err
	}

	// Migrations for databases created by older versions.
	hasColor, err := db.columnExists("connection_types", "color")
	if err != nil {
		return err
	}
	if !hasColor {
		if _, err := db.conn.Exec("ALTER TABLE connection_types ADD COLUMN color TEXT DEFAULT '#666'"); err != nil {
			return fmt.Errorf("failed to add connection_types.color: %w", err)
		}
	}

	hasWeight, err := db.columnExists("connections", "weight")
	if err != nil {
		return err
	}
	if hasWeight {
		if _, err := db.conn.Exec("ALTER TABLE connections DROP COLUMN weight"); err != nil {
			return fmt.Errorf("failed to drop connections.weight: %w", err)
		}
	}
	return nil
}

func (db *Database) columnExists(table, column string) (bool, error) {
	var one int
	err := db.conn.QueryRow("SELECT 1 FROM pragma_table_info(?) WHERE name = ?", table, column).Scan(&one)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to inspect %s.%s: %w", table, column, err)
	}
	return true, nil
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
		return nil, translate(err, ErrTypeNameTaken, nil)
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
		return nil, translate(err, nil, nil)
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
	err := db.execAffecting(
		"UPDATE connection_types SET name = ?, description = ?, color = ? WHERE id = ?",
		name, description, color, id,
	)
	if err != nil {
		return nil, translate(err, ErrTypeNameTaken, nil)
	}

	return db.GetConnectionType(id)
}

func (db *Database) DeleteConnectionType(id int) error {
	err := db.execAffecting("DELETE FROM connection_types WHERE id = ?", id)
	return translate(err, nil, ErrTypeInUse)
}

// Person operations
func (db *Database) CreatePerson(name, description string) (*Person, error) {
	result, err := db.conn.Exec(
		"INSERT INTO people (name, description) VALUES (?, ?)",
		name, description,
	)
	if err != nil {
		return nil, translate(err, ErrPersonNameTaken, nil)
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
		return nil, translate(err, nil, nil)
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
	return db.execAffecting("DELETE FROM people WHERE id = ?", id)
}

// Connection operations

func (db *Database) CreateConnection(fromID, toID, typeID int, description string) (*Connection, error) {
	result, err := db.conn.Exec(
		"INSERT INTO connections (from_person_id, to_person_id, type_id, description) VALUES (?, ?, ?, ?)",
		fromID, toID, typeID, description,
	)
	if err != nil {
		return nil, translate(err, ErrConnectionExists, ErrInvalidReference)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return db.GetConnection(int(id))
}

const connectionColumns = "SELECT id, from_person_id, to_person_id, type_id, description, created_at, updated_at FROM connections"

func (db *Database) GetConnection(id int) (*Connection, error) {
	row := db.conn.QueryRow(connectionColumns+" WHERE id = ?", id)

	var c Connection
	err := row.Scan(&c.ID, &c.FromID, &c.ToID, &c.TypeID, &c.Description, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, translate(err, nil, nil)
	}

	return &c, nil
}

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
		if err := rows.Scan(&c.ID, &c.FromID, &c.ToID, &c.TypeID, &c.Description, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		connections = append(connections, c)
	}

	return connections, rows.Err()
}

func (db *Database) DeleteConnection(id int) error {
	return db.execAffecting("DELETE FROM connections WHERE id = ?", id)
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
