package main

import (
	"errors"
	"log"
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

// respondError maps domain errors to HTTP responses. Unknown errors are
// logged and reported as a generic 500 so internal details are not leaked.
func respondError(c *gin.Context, err error, notFoundMsg string) {
	switch {
	case errors.Is(err, ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": notFoundMsg})
	case errors.Is(err, ErrPersonNameTaken),
		errors.Is(err, ErrTypeNameTaken),
		errors.Is(err, ErrConnectionExists),
		errors.Is(err, ErrTypeInUse):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, ErrInvalidReference):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		log.Printf("%s %s: %v", c.Request.Method, c.Request.URL.Path, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
	}
}

func parseID(c *gin.Context, invalidMsg string) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": invalidMsg})
		return 0, false
	}
	return id, true
}

// ============ ConnectionType Handlers ============

type ConnectionTypeRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

func (api *API) CreateConnectionType(c *gin.Context) {
	var req ConnectionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#666"
	}

	ct, err := api.db.CreateConnectionType(req.Name, req.Description, req.Color)
	if err != nil {
		respondError(c, err, "connection type not found")
		return
	}

	c.JSON(http.StatusCreated, ct)
}

func (api *API) GetConnectionTypes(c *gin.Context) {
	types, err := api.db.GetAllConnectionTypes()
	if err != nil {
		respondError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, types)
}

func (api *API) UpdateConnectionType(c *gin.Context) {
	id, ok := parseID(c, "invalid connection type id")
	if !ok {
		return
	}

	var req ConnectionTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Color == "" {
		req.Color = "#666"
	}

	ct, err := api.db.UpdateConnectionType(id, req.Name, req.Description, req.Color)
	if err != nil {
		respondError(c, err, "connection type not found")
		return
	}

	c.JSON(http.StatusOK, ct)
}

func (api *API) DeleteConnectionType(c *gin.Context) {
	id, ok := parseID(c, "invalid connection type id")
	if !ok {
		return
	}

	if err := api.db.DeleteConnectionType(id); err != nil {
		respondError(c, err, "connection type not found")
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
		respondError(c, err, "person not found")
		return
	}

	c.JSON(http.StatusCreated, person)
}

func (api *API) GetPeople(c *gin.Context) {
	people, err := api.db.GetAllPeople()
	if err != nil {
		respondError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, people)
}

func (api *API) GetPerson(c *gin.Context) {
	id, ok := parseID(c, "invalid person id")
	if !ok {
		return
	}

	person, err := api.db.GetPerson(id)
	if err != nil {
		respondError(c, err, "person not found")
		return
	}

	c.JSON(http.StatusOK, person)
}

func (api *API) DeletePerson(c *gin.Context) {
	id, ok := parseID(c, "invalid person id")
	if !ok {
		return
	}

	if err := api.db.DeletePerson(id); err != nil {
		respondError(c, err, "person not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "person deleted"})
}

// ============ Connection Handlers ============

type CreateConnectionRequest struct {
	FromPersonID int    `json:"from_person_id" binding:"required"`
	ToPersonID   int    `json:"to_person_id" binding:"required"`
	TypeID       int    `json:"type_id" binding:"required"`
	Description  string `json:"description"`
}

func (api *API) CreateConnection(c *gin.Context) {
	var req CreateConnectionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conn, err := api.db.CreateConnection(req.FromPersonID, req.ToPersonID, req.TypeID, req.Description)
	if err != nil {
		respondError(c, err, "connection not found")
		return
	}

	c.JSON(http.StatusCreated, conn)
}

func (api *API) GetConnections(c *gin.Context) {
	connections, err := api.db.GetAllConnections()
	if err != nil {
		respondError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (api *API) GetPersonConnections(c *gin.Context) {
	personID, ok := parseID(c, "invalid person id")
	if !ok {
		return
	}

	connections, err := api.db.GetConnectionsByPerson(personID)
	if err != nil {
		respondError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, connections)
}

func (api *API) DeleteConnection(c *gin.Context) {
	id, ok := parseID(c, "invalid connection id")
	if !ok {
		return
	}

	if err := api.db.DeleteConnection(id); err != nil {
		respondError(c, err, "connection not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "connection deleted"})
}

// ============ Graph Handler ============

func (api *API) GetGraph(c *gin.Context) {
	graph, err := api.db.GetGraph()
	if err != nil {
		respondError(c, err, "")
		return
	}

	c.JSON(http.StatusOK, graph)
}

// ============ Health Check ============

func (api *API) Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
