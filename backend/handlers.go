package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type API struct {
	db *Database
}

func NewAPI(db *Database) *API {
	return &API{db: db}
}

// ============ ConnectionType Handlers ============

// CreateConnectionTypeRequest represents the request body
type CreateConnectionTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type UpdateConnectionTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (api *API) CreateConnectionType(c *gin.Context) {
	var req CreateConnectionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#666"
	}

	ct, err := api.db.CreateConnectionType(req.Name, req.Description, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ct)
}

func (api *API) GetConnectionTypes(c *gin.Context) {
	types, err := api.db.GetAllConnectionTypes()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, types)
}

func (api *API) UpdateConnectionType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection type id"})
		return
	}

	var req UpdateConnectionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#666"
	}

	ct, err := api.db.UpdateConnectionType(id, req.Name, req.Description, req.Color)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ct)
}

func (api *API) DeleteConnectionType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection type id"})
		return
	}

	err = api.db.DeleteConnectionType(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection type deleted"})
}

// ============ Person Handlers ============

type CreatePersonRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

func (api *API) CreatePerson(c *gin.Context) {
	var req CreatePersonRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person, err := api.db.CreatePerson(req.Name, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (api *API) GetPeople(c *gin.Context) {
	people, err := api.db.GetAllPeople()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, people)
}

func (api *API) GetPerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	person, err := api.db.GetPerson(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "person not found"})
		return
	}

	c.JSON(http.StatusOK, person)
}

func (api *API) DeletePerson(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	err = api.db.DeletePerson(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person deleted"})
}

// ============ Connection Handlers ============

type CreateConnectionRequest struct {
	FromPersonID int     `json:"from_person_id" binding:"required"`
	ToPersonID   int     `json:"to_person_id" binding:"required"`
	TypeID       int     `json:"type_id" binding:"required"`
	Description  string  `json:"description"`
	Weight       float64 `json:"weight"`
}

func (api *API) CreateConnection(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := api.db.CreateConnection(req.FromPersonID, req.ToPersonID, req.TypeID, req.Description, req.Weight)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conn)
}

func (api *API) GetConnections(c *gin.Context) {
	connections, err := api.db.GetAllConnections()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (api *API) GetPersonConnections(c *gin.Context) {
	personID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid person id"})
		return
	}

	connections, err := api.db.GetConnectionsByPerson(personID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (api *API) DeleteConnection(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid connection id"})
		return
	}

	err = api.db.DeleteConnection(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection deleted"})
}

// ============ Graph Handler ============

func (api *API) GetGraph(c *gin.Context) {
	graph, err := api.db.GetGraph()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, graph)
}

// ============ Health Check ============

func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
