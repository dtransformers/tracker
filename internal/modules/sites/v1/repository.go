package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"
	"tracking/internal/ent/site"
)

type Repository interface {
	Create(ctx context.Context, input CreateSiteInput) (*ent.Site, error)

	GetByUUID(ctx context.Context, uuid string) (*ent.Site, error)

	GetByDomain(ctx context.Context, domain string) (*ent.Site, error)
	GetByPublicKey(ctx context.Context, publicKey string) (*ent.Site, error)
	List(ctx context.Context) ([]*ent.Site, error)

	Update(ctx context.Context, uuid string, input UpdateSiteInput) (*ent.Site, error)

	Delete(ctx context.Context, uuid string) error

	RegenerateKeys(ctx context.Context, uuid, publicKey, secretKey string) (*ent.Site, error)

	ListActive(ctx context.Context) ([]*ent.Site, error)
	ListByTenant(ctx context.Context, tenantID int64) ([]*ent.Site, error)
	Count(ctx context.Context) (int, error)
}

type entRepository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	if client == nil {
		panic("ent client cannot be nil")
	}
	return &entRepository{client: client}
}

type CreateSiteInput struct {
	UUID      string
	Name      string
	Domain    string
	TenantID  int64
	PublicKey string
	SecretKey string
}

type UpdateSiteInput struct {
	Name     *string
	Domain   *string
	TenantID *int64
	IsActive *bool
}

func (r *entRepository) Create(ctx context.Context, input CreateSiteInput) (*ent.Site, error) {
	site, err := r.client.Site.
		Create().
		SetUUID(input.UUID).
		SetName(input.Name).
		SetDomain(input.Domain).
		SetTenantID(input.TenantID).
		SetPublicKey(input.PublicKey).
		SetSecretKey(input.SecretKey).
		SetIsActive(true). // Default to active
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	return site, nil
}

func (r *entRepository) GetByUUID(ctx context.Context, uuid string) (*ent.Site, error) {
	site, err := r.client.Site.
		Query().
		Where(site.UUID(uuid)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get site by uuid: %w", err)
	}

	return site, nil
}

// GetByDomain retrieves a site by its domain.
func (r *entRepository) GetByDomain(ctx context.Context, domain string) (*ent.Site, error) {
	site, err := r.client.Site.
		Query().
		Where(site.Domain(domain)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get site by domain: %w", err)
	}

	return site, nil
}

// GetByPublicKey retrieves a site by its public key.
func (r *entRepository) GetByPublicKey(ctx context.Context, publicKey string) (*ent.Site, error) {
	site, err := r.client.Site.
		Query().
		Where(site.PublicKey(publicKey)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to get site by public key: %w", err)
	}

	return site, nil
}

// List retrieves all sites ordered by creation date.
func (r *entRepository) List(ctx context.Context) ([]*ent.Site, error) {
	sites, err := r.client.Site.
		Query().
		Order(ent.Desc(site.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	return sites, nil
}

// ListActive retrieves all active sites.
func (r *entRepository) ListActive(ctx context.Context) ([]*ent.Site, error) {
	sites, err := r.client.Site.
		Query().
		Where(site.IsActive(true)).
		Order(ent.Desc(site.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list active sites: %w", err)
	}

	return sites, nil
}

// ListByTenant retrieves all sites for a specific tenant.
func (r *entRepository) ListByTenant(ctx context.Context, tenantID int64) ([]*ent.Site, error) {
	sites, err := r.client.Site.
		Query().
		Where(site.TenantID(tenantID)).
		Order(ent.Desc(site.FieldCreatedAt)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to list sites by tenant: %w", err)
	}

	return sites, nil
}

// Update updates a site with the provided fields.
func (r *entRepository) Update(ctx context.Context, uuid string, input UpdateSiteInput) (*ent.Site, error) {
	// Start building the update query
	updateQuery := r.client.Site.
		Update().
		Where(site.UUID(uuid))

	// Apply updates only for non-nil fields
	if input.Name != nil {
		updateQuery = updateQuery.SetName(*input.Name)
	}

	if input.Domain != nil {
		updateQuery = updateQuery.SetDomain(*input.Domain)
	}

	if input.TenantID != nil {
		updateQuery = updateQuery.SetTenantID(*input.TenantID)
	}

	if input.IsActive != nil {
		updateQuery = updateQuery.SetIsActive(*input.IsActive)
	}

	// Execute the update
	affected, err := updateQuery.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update site: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	// Retrieve and return the updated site
	return r.GetByUUID(ctx, uuid)
}

// Delete deletes a site by its UUID.
func (r *entRepository) Delete(ctx context.Context, uuid string) error {
	affected, err := r.client.Site.
		Delete().
		Where(site.UUID(uuid)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete site: %w", err)
	}

	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}

// RegenerateKeys updates the security keys for a site.
func (r *entRepository) RegenerateKeys(ctx context.Context, uuid, publicKey, secretKey string) (*ent.Site, error) {
	affected, err := r.client.Site.
		Update().
		Where(site.UUID(uuid)).
		SetPublicKey(publicKey).
		SetSecretKey(secretKey).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to regenerate keys: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	// Retrieve and return the updated site
	return r.GetByUUID(ctx, uuid)
}

// Count returns the total number of sites.
func (r *entRepository) Count(ctx context.Context) (int, error) {
	count, err := r.client.Site.
		Query().
		Count(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to count sites: %w", err)
	}

	return count, nil
}

// Transaction support methods (optional but useful)

// WithTx returns a new repository that uses the given transaction.
func (r *entRepository) WithTx(tx *ent.Tx) Repository {
	return &entRepository{client: tx.Client()}
}

// BeginTx starts a new transaction.
func (r *entRepository) BeginTx(ctx context.Context) (*ent.Tx, error) {
	return r.client.Tx(ctx)
}
