package v1

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"tracking/internal/ent"

	"github.com/google/uuid"
)

// Service handles the business logic for sites.
type Service struct {
	repo Repository
}

// NewService creates a new Site service with the given repository.
func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

// CreateRequest defines the payload for creating a new site.
type CreateRequest struct {
	// Name of the site
	Name   string `json:"name" binding:"required"`
	// Domain of the site (must be unique)
	Domain string `json:"domain" binding:"required"`
}

// UpdateRequest defines the payload for updating an existing site.
type UpdateRequest struct {
	// Optional new name for the site
	Name     *string `json:"name"`
	// Optional new domain for the site
	Domain   *string `json:"domain"`
	// Optional active status
	IsActive *bool   `json:"is_active"`
}

func (s *Service) CreateSite(ctx context.Context, req CreateRequest) (*ent.Site, error) {
	u := uuid.New().String()
	pubKey := generateRandomKey(16)
	secKey := generateRandomKey(32)

	return s.repo.Create(ctx, CreateSiteInput{
		UUID:      u,
		Name:      req.Name,
		Domain:    req.Domain,
		PublicKey: pubKey,
		SecretKey: secKey,
	})
}

func (s *Service) GetSite(ctx context.Context, uuid string) (*ent.Site, error) {
	return s.repo.GetByUUID(ctx, uuid)
}

func (s *Service) ListSites(ctx context.Context) ([]*ent.Site, error) {
	return s.repo.List(ctx)
}

func (s *Service) UpdateSite(ctx context.Context, uuid string, req UpdateRequest) (*ent.Site, error) {
	return s.repo.Update(ctx, uuid, UpdateSiteInput{
		Name:     req.Name,
		Domain:   req.Domain,
		IsActive: req.IsActive,
	})
}

func (s *Service) DeleteSite(ctx context.Context, uuid string) error {
	return s.repo.Delete(ctx, uuid)
}

func generateRandomKey(length int) string {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return ""
	}
	return hex.EncodeToString(b)
}
