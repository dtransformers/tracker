package v1

import "time"

type CreateRequest struct {
	// Name of the site (required, 1-255 characters)
	Name string `json:"name" binding:"required"`
	// Domain of the site (required, must be unique and valid)
	Domain string `json:"domain" binding:"required"`
	// Tenant ID of the site
	TenantID int64 `json:"tenant_id" binding:"required"`
}

type UpdateRequest struct {
	// Optional new name for the site
	Name *string `json:"name"`
	// Optional new domain for the site
	Domain *string `json:"domain"`
	// Optional new tenant ID
	TenantID *int64 `json:"tenant_id"`
	// Optional active status
	IsActive *bool `json:"is_active"`
}

type Response struct {
	UUID      string `json:"uuid"`
	Name      string `json:"name"`
	Domain    string `json:"domain"`
	TenantID  int64  `json:"tenant_id"`
	PublicKey string `json:"public_key"`

	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type FullSiteResponse struct {
	UUID      string    `json:"uuid"`
	Name      string    `json:"name"`
	Domain    string    `json:"domain"`
	TenantID  int64     `json:"tenant_id"`
	PublicKey string    `json:"public_key"`
	SecretKey string    `json:"secret_key"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ListResponse struct {
	Sites []*Response `json:"sites"`
	Total int         `json:"total"`
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
