package v1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"tracking/internal/ent"
	"tracking/internal/ent/enttest"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTest(t *testing.T) (*gin.Engine, *ent.Client) {
	gin.SetMode(gin.TestMode)
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")

	router := gin.New()
	group := router.Group("/api/v1")
	Register(group, client)

	return router, client
}

func TestCreateSite(t *testing.T) {
	router, client := setupTest(t)
	defer client.Close()

	t.Run("successful creation", func(t *testing.T) {
		reqBody := CreateRequest{
			Name:   gofakeit.AppName(),
			Domain: gofakeit.DomainName(),
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/v1/sites", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var res map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &res)

		assert.Equal(t, reqBody.Name, res["name"])
		assert.Equal(t, reqBody.Domain, res["domain"])
		assert.NotEmpty(t, res["uuid"])
		assert.NotEmpty(t, res["public_key"])
	})
}

func FuzzCreateSite(f *testing.F) {
	// Add seed corpus
	f.Add("Valid Site", "valid.com")
	f.Add("", "empty-name.com")
	f.Add("No Domain", "")
	f.Add("Long Name "+fmt.Sprint(make([]byte, 1000)), "long.com")

	f.Fuzz(func(t *testing.T, name string, domain string) {
		// Fresh memory DB for each fuzz iteration
		client := enttest.Open(t, "sqlite3", "file:fuzz?mode=memory&cache=shared&_fk=1")
		defer client.Close()

		repo := NewRepository(client)
		svc := NewService(repo)
		handler := NewHandler(svc)

		router := gin.New()
		router.POST("/sites", handler.Create)

		reqBody := CreateRequest{
			Name:   name,
			Domain: domain,
		}
		jsonBody, _ := json.Marshal(reqBody)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/sites", bytes.NewBuffer(jsonBody))
		req.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(w, req)

		// We just want to make sure it doesn't panic
		// and returns a reasonable status code
		if w.Code != http.StatusCreated && w.Code != http.StatusBadRequest && w.Code != http.StatusInternalServerError {
			t.Errorf("Unexpected status code: %d", w.Code)
		}
	})
}

// Helper for faking input in tests
func fakeCreateRequest() CreateRequest {
	return CreateRequest{
		Name:   gofakeit.AppName(),
		Domain: gofakeit.DomainName(),
	}
}
