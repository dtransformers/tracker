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
		TenantID:  10,
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
	assert.Equal(t, input.TenantID, site.TenantID)
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
		TenantID:  1,
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
	assert.Equal(t, created.TenantID, retrieved.TenantID)
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
		TenantID:  1,
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

func TestRepository_ListByTenant(t *testing.T) {
	client := enttest.Open(t, "sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()

	repo := v1.NewRepository(client)
	ctx := context.Background()

	tenant1 := int64(100)
	tenant2 := int64(200)

	// Create sites for tenant 1
	repo.Create(ctx, v1.CreateSiteInput{UUID: "s1", Name: "Site 1", Domain: "s1.com", TenantID: tenant1, PublicKey: "p1", SecretKey: "s1"})
	repo.Create(ctx, v1.CreateSiteInput{UUID: "s2", Name: "Site 2", Domain: "s2.com", TenantID: tenant1, PublicKey: "p2", SecretKey: "s2"})

	// Create site for tenant 2
	repo.Create(ctx, v1.CreateSiteInput{UUID: "s3", Name: "Site 3", Domain: "s3.com", TenantID: tenant2, PublicKey: "p3", SecretKey: "s3"})

	// List for tenant 1
	sites1, err := repo.ListByTenant(ctx, tenant1)
	require.NoError(t, err)
	assert.Len(t, sites1, 2)

	// List for tenant 2
	sites2, err := repo.ListByTenant(ctx, tenant2)
	require.NoError(t, err)
	assert.Len(t, sites2, 1)
	assert.Equal(t, "s3", sites2[0].UUID)
}
