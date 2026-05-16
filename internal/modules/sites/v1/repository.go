package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"
	"tracking/internal/ent/site"
)

// Repository defines the interface for site data access.
type Repository interface {
	// Create persists a new site record.
	Create(ctx context.Context, input CreateSiteInput) (*ent.Site, error)
	// GetByID retrieves a site by its internal ID.
	GetByID(ctx context.Context, id int) (*ent.Site, error)
	// GetByUUID retrieves a site by its public UUID.
	GetByUUID(ctx context.Context, uuid string) (*ent.Site, error)
	// List returns all non-deleted sites.
	List(ctx context.Context) ([]*ent.Site, error)
	// Update modifies an existing site.
	Update(ctx context.Context, uuid string, input UpdateSiteInput) (*ent.Site, error)
	// Delete performs a soft delete on a site.
	Delete(ctx context.Context, uuid string) error
}

type siteRepository struct {
	client *ent.Client
}

// NewRepository creates a new site repository using the ent client.
func NewRepository(client *ent.Client) Repository {
	return &siteRepository{client: client}
}

type CreateSiteInput struct {
	UUID      string
	Name      string
	Domain    string
	PublicKey string
	SecretKey string
}

type UpdateSiteInput struct {
	Name     *string
	Domain   *string
	IsActive *bool
}

func (r *siteRepository) Create(ctx context.Context, input CreateSiteInput) (*ent.Site, error) {
	return r.client.Site.Create().
		SetUUID(input.UUID).
		SetName(input.Name).
		SetDomain(input.Domain).
		SetPublicKey(input.PublicKey).
		SetSecretKey(input.SecretKey).
		Save(ctx)
}

func (r *siteRepository) GetByID(ctx context.Context, id int) (*ent.Site, error) {
	s, err := r.client.Site.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("site not found: %w", err)
		}
		return nil, err
	}
	return s, nil
}

func (r *siteRepository) GetByUUID(ctx context.Context, uuidStr string) (*ent.Site, error) {
	s, err := r.client.Site.Query().
		Where(site.UUIDEQ(uuidStr)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("site not found: %w", err)
		}
		return nil, err
	}
	return s, nil
}

func (r *siteRepository) List(ctx context.Context) ([]*ent.Site, error) {
	return r.client.Site.Query().
		Where(site.IsDeletedEQ(false)).
		All(ctx)
}

func (r *siteRepository) Update(ctx context.Context, uuidStr string, input UpdateSiteInput) (*ent.Site, error) {
	updater := r.client.Site.Update().
		Where(site.UUIDEQ(uuidStr))

	if input.Name != nil {
		updater.SetName(*input.Name)
	}
	if input.Domain != nil {
		updater.SetDomain(*input.Domain)
	}
	if input.IsActive != nil {
		updater.SetIsActive(*input.IsActive)
	}

	n, err := updater.Save(ctx)
	if err != nil {
		return nil, err
	}
	if n == 0 {
		return nil, fmt.Errorf("site not found")
	}

	return r.GetByUUID(ctx, uuidStr)
}

func (r *siteRepository) Delete(ctx context.Context, uuidStr string) error {
	// Soft delete
	err := r.client.Site.Update().
		Where(site.UUIDEQ(uuidStr)).
		SetIsDeleted(true).
		Exec(ctx)
	return err
}
