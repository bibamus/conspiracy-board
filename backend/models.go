package main

import "time"

type ConnectionType struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Color       string    `json:"color"`
	CreatedAt   time.Time `json:"created_at"`
}

type Person struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Connection struct {
	ID          int       `json:"id"`
	FromID      int       `json:"from_person_id"`
	ToID        int       `json:"to_person_id"`
	TypeID      int       `json:"type_id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GraphNode struct {
	Person      Person       `json:"person"`
	Connections []Connection `json:"connections"`
}

type Graph struct {
	Nodes []*GraphNode     `json:"nodes"`
	Types []ConnectionType `json:"types"`
}
