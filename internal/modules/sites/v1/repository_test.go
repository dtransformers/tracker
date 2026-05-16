package v1_test

import (
	"context"
	"testing"

	"tracking/internal/ent/enttest"
	v1 "tracking/internal/modules/sites/v1"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRepository_Create(t *testing.T) {
	// Setup in-memory database for testing
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := v1.NewRepository(client)
	ctx := context.Background()

	input := v1.CreateSiteInput{
		UUID:      "test-uuid-123",
		Name:      "Test Site",
		Domain:    "test.example.com",
		PublicKey: "public-key-123",
		SecretKey: "secret-key-456",
	}

	site, err := repo.Create(ctx, input)

	require.NoError(t, err)
	assert.Equal(t, input.UUID, site.UUID)
	assert.Equal(t, input.Name, site.Name)
	assert.Equal(t, input.Domain, site.Domain)
	assert.Equal(t, input.PublicKey, site.PublicKey)
	assert.Equal(t, input.SecretKey, site.SecretKey)
	assert.True(t, site.IsActive)
}

func TestRepository_GetByUUID(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := v1.NewRepository(client)
	ctx := context.Background()

	// Create a site first
	input := v1.CreateSiteInput{
		UUID:      "test-uuid-456",
		Name:      "Test Site",
		Domain:    "test.example.com",
		PublicKey: "public-key-123",
		SecretKey: "secret-key-456",
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	// Retrieve it
	retrieved, err := repo.GetByUUID(ctx, created.UUID)

	require.NoError(t, err)
	assert.Equal(t, created.UUID, retrieved.UUID)
	assert.Equal(t, created.Name, retrieved.Name)
}

func TestRepository_Update(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := v1.NewRepository(client)
	ctx := context.Background()

	// Create a site
	input := v1.CreateSiteInput{
		UUID:      "test-uuid-789",
		Name:      "Original Name",
		Domain:    "original.example.com",
		PublicKey: "public-key-123",
		SecretKey: "secret-key-456",
	}
	created, err := repo.Create(ctx, input)
	require.NoError(t, err)

	// Update it
	newName := "Updated Name"
	updateInput := v1.UpdateSiteInput{
		Name: &newName,
	}

	updated, err := repo.Update(ctx, created.UUID, updateInput)

	require.NoError(t, err)
	assert.Equal(t, newName, updated.Name)
	assert.Equal(t, created.Domain, updated.Domain) // Should remain unchanged
}
