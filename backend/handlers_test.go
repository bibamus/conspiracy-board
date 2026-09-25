package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func newTestRouter(t *testing.T) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	db, err := NewDatabase(filepath.Join(t.TempDir(), "api.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return newRouter(NewAPI(db))
}

func doRequest(t *testing.T, r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestEmptyCollectionsAreJSONArrays(t *testing.T) {
	r := newTestRouter(t)

	for _, path := range []string{"/api/people", "/api/connection-types", "/api/connections", "/api/people/1/connections"} {
		rec := doRequest(t, r, http.MethodGet, path, "")
		if rec.Code != http.StatusOK || strings.TrimSpace(rec.Body.String()) != "[]" {
			t.Errorf("GET %s: got %d %q, want 200 []", path, rec.Code, rec.Body.String())
		}
	}

	if rec := doRequest(t, r, http.MethodPost, "/api/people", `{"name":"Alice"}`); rec.Code != http.StatusCreated {
		t.Fatalf("create person: %d %s", rec.Code, rec.Body.String())
	}
	rec := doRequest(t, r, http.MethodGet, "/api/graph", "")
	want := `"connections":[]`
	if !strings.Contains(rec.Body.String(), want) || !strings.Contains(rec.Body.String(), `"types":[]`) {
		t.Errorf("GET /api/graph: got %s, want empty arrays", rec.Body.String())
	}
}

// Guards the SQLite error-code mapping in database.go against driver changes.
func TestConstraintErrorsMapToStatusCodes(t *testing.T) {
	r := newTestRouter(t)
	steps := []struct {
		method, path, body string
		want               int
	}{
		{http.MethodPost, "/api/people", `{"name":"A"}`, http.StatusCreated},
		{http.MethodPost, "/api/people", `{"name":"B"}`, http.StatusCreated},
		{http.MethodPost, "/api/people", `{"name":"A"}`, http.StatusConflict},
		{http.MethodPost, "/api/connection-types", `{"name":"knows"}`, http.StatusCreated},
		{http.MethodPost, "/api/connection-types", `{"name":"knows"}`, http.StatusConflict},
		{http.MethodPost, "/api/connections", `{"from_person_id":1,"to_person_id":2,"type_id":1}`, http.StatusCreated},
		{http.MethodPost, "/api/connections", `{"from_person_id":1,"to_person_id":2,"type_id":1}`, http.StatusConflict},
		{http.MethodPost, "/api/connections", `{"from_person_id":1,"to_person_id":99,"type_id":1}`, http.StatusBadRequest},
		{http.MethodDelete, "/api/connection-types/1", ``, http.StatusConflict},
		{http.MethodGet, "/api/people/99", ``, http.StatusNotFound},
		{http.MethodGet, "/api/graph", ``, http.StatusOK},
	}
	for _, s := range steps {
		if rec := doRequest(t, r, s.method, s.path, s.body); rec.Code != s.want {
			t.Fatalf("%s %s %s: got %d %s, want %d", s.method, s.path, s.body, rec.Code, rec.Body.String(), s.want)
		}
	}
}

func TestValidationErrorsAreUserFacing(t *testing.T) {
	r := newTestRouter(t)
	cases := []struct {
		path, body, want string
	}{
		{"/api/people", `{}`, "name is required"},
		{"/api/connection-types", `{"description":"x"}`, "name is required"},
		{"/api/connections", `{"from_person_id":1}`, "to_person_id is required; type_id is required"},
		{"/api/people", `{"name":`, "request body must be valid JSON"},
		{"/api/people", ``, "request body must be valid JSON"},
		{"/api/connections", `{"from_person_id":"1","to_person_id":2,"type_id":3}`, "from_person_id must be of type int"},
	}

	for _, tc := range cases {
		rec := doRequest(t, r, http.MethodPost, tc.path, tc.body)
		var resp struct{ Error string }
		json.Unmarshal(rec.Body.Bytes(), &resp)
		if rec.Code != http.StatusBadRequest || resp.Error != tc.want {
			t.Errorf("POST %s %s: got %d %q, want 400 %q", tc.path, tc.body, rec.Code, resp.Error, tc.want)
		}
	}
}
