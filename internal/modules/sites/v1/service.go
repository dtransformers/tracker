package v1

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"
	"time"
	"tracking/internal/ent"

	"github.com/google/uuid"
)

// Configuration constants for key generation
const (
	PublicKeyBytes  = 32 // Generates 64 hex characters
	SecretKeyBytes  = 64 // Generates 128 hex characters
	DefaultTimeout  = 5 * time.Second
	MaxNameLength   = 255
	MinNameLength   = 1
	MaxDomainLength = 253 // RFC 1035
)

var (
	ErrInvalidDomain = errors.New("invalid domain format")
	ErrInvalidName   = errors.New("name must be between 1 and 255 characters")
	ErrEmptyUpdate   = errors.New("at least one field must be updated")
	ErrDomainExists  = errors.New("domain already exists")
	ErrKeyGeneration = errors.New("failed to generate security key")
	ErrSiteNotFound  = errors.New("site not found")
)

const errInvalidUUIDFormat = "invalid UUID format: %w"

// Service handles the business logic for sites.
type Service struct {
	repo         Repository
	logger       *slog.Logger
	pubKeyLength int
	secKeyLength int
	timeout      time.Duration
}

// ServiceConfig holds configuration options for the Service.
type ServiceConfig struct {
	PublicKeyLength int
	SecretKeyLength int
	Timeout         time.Duration
}

// NewService creates a new Site service with the given repository and logger.
func NewService(repo Repository, logger *slog.Logger) *Service {
	return NewServiceWithConfig(repo, logger, ServiceConfig{
		PublicKeyLength: PublicKeyBytes,
		SecretKeyLength: SecretKeyBytes,
		Timeout:         DefaultTimeout,
	})
}

// NewServiceWithConfig creates a new Site service with custom configuration.
func NewServiceWithConfig(repo Repository, logger *slog.Logger, cfg ServiceConfig) *Service {
	if logger == nil {
		logger = slog.Default()
	}

	return &Service{
		repo:         repo,
		logger:       logger,
		pubKeyLength: cfg.PublicKeyLength,
		secKeyLength: cfg.SecretKeyLength,
		timeout:      cfg.Timeout,
	}
}

// CreateSite creates a new site with generated security keys.
func (s *Service) CreateSite(ctx context.Context, req CreateRequest) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "creating site",
		"name", req.Name,
		"domain", req.Domain,
		"tenant_id", req.TenantID,
	)

	if err := s.validateCreateRequest(req); err != nil {
		s.logger.WarnContext(ctx, "validation failed",
			"error", err,
			"domain", req.Domain,
		)
		return nil, err
	}

	if err := s.checkDomainUniqueness(ctx, req.Domain, ""); err != nil {
		return nil, err
	}

	pubKey, err := generateSecureKey(s.pubKeyLength)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate public key", "error", err)
		return nil, fmt.Errorf("%w: public key", ErrKeyGeneration)
	}

	secKey, err := generateSecureKey(s.secKeyLength)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate secret key", "error", err)
		return nil, fmt.Errorf("%w: secret key", ErrKeyGeneration)
	}

	site, err := s.repo.Create(ctx, CreateSiteInput{
		UUID:      uuid.New().String(),
		Name:      strings.TrimSpace(req.Name),
		Domain:    normalizeDomain(req.Domain),
		TenantID:  req.TenantID,
		PublicKey: pubKey,
		SecretKey: secKey,
	})
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to create site",
			"error", err,
			"domain", req.Domain,
		)
		return nil, fmt.Errorf("failed to create site: %w", err)
	}

	s.logger.InfoContext(ctx, "site created successfully",
		"uuid", site.UUID,
		"domain", site.Domain,
	)

	return toResponse(site), nil
}

func (s *Service) GetSite(ctx context.Context, siteUUID string) (*FullSiteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if _, err := uuid.Parse(siteUUID); err != nil {
		return nil, fmt.Errorf(errInvalidUUIDFormat, err)
	}

	site, err := s.repo.GetByUUID(ctx, siteUUID)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrSiteNotFound
		}
		s.logger.ErrorContext(ctx, "failed to get site",
			"error", err,
			"uuid", siteUUID,
		)
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	return toFullResponse(site), nil
}

func (s *Service) VerifySiteKey(ctx context.Context, publicKey string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	site, err := s.repo.GetByPublicKey(ctx, publicKey)
	if err != nil {
		if ent.IsNotFound(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to verify site key: %w", err)
	}

	return site.IsActive, nil
}

func (s *Service) GetSiteKey(ctx context.Context, publicKey string) (*ent.Site, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	site, err := s.repo.GetByPublicKey(ctx, publicKey)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrSiteNotFound
		}
		return nil, fmt.Errorf("failed to get site by key: %w", err)
	}

	return site, nil
}

func (s *Service) ListSites(ctx context.Context, tenantID *int64) (*ListResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var sites []*ent.Site
	var err error

	if tenantID != nil {
		s.logger.InfoContext(ctx, "listing sites for tenant", "tenant_id", *tenantID)
		sites, err = s.repo.ListByTenant(ctx, *tenantID)
	} else {
		s.logger.InfoContext(ctx, "listing all sites")
		sites, err = s.repo.List(ctx)
	}

	if err != nil {
		s.logger.ErrorContext(ctx, "failed to list sites", "error", err)
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	responses := make([]*Response, len(sites))
	for i, site := range sites {
		responses[i] = toResponse(site)
	}

	return &ListResponse{
		Sites: responses,
		Total: len(responses),
	}, nil
}

func (s *Service) UpdateSite(ctx context.Context, siteUUID string, req UpdateRequest) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "updating site", "uuid", siteUUID)

	if _, err := uuid.Parse(siteUUID); err != nil {
		return nil, fmt.Errorf(errInvalidUUIDFormat, err)
	}

	if err := s.validateUpdateRequest(req); err != nil {
		s.logger.WarnContext(ctx, "validation failed",
			"error", err,
			"uuid", siteUUID,
		)
		return nil, err
	}

	if req.Domain != nil {
		if err := s.checkDomainUniqueness(ctx, *req.Domain, siteUUID); err != nil {
			return nil, err
		}
	}

	input := UpdateSiteInput{
		Name:     req.Name,
		TenantID: req.TenantID,
		IsActive: req.IsActive,
	}
	if req.Domain != nil {
		normalized := normalizeDomain(*req.Domain)
		input.Domain = &normalized
	}

	site, err := s.repo.Update(ctx, siteUUID, input)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrSiteNotFound
		}
		s.logger.ErrorContext(ctx, "failed to update site",
			"error", err,
			"uuid", siteUUID,
		)
		return nil, fmt.Errorf("failed to update site: %w", err)
	}

	s.logger.InfoContext(ctx, "site updated successfully",
		"uuid", site.UUID,
		"domain", site.Domain,
	)

	return toResponse(site), nil
}

func (s *Service) DeleteSite(ctx context.Context, siteUUID string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "deleting site", "uuid", siteUUID)

	if _, err := uuid.Parse(siteUUID); err != nil {
		return fmt.Errorf(errInvalidUUIDFormat, err)
	}

	err := s.repo.Delete(ctx, siteUUID)
	if err != nil {
		if ent.IsNotFound(err) {
			return ErrSiteNotFound
		}
		s.logger.ErrorContext(ctx, "failed to delete site",
			"error", err,
			"uuid", siteUUID,
		)
		return fmt.Errorf("failed to delete site: %w", err)
	}

	s.logger.InfoContext(ctx, "site deleted successfully", "uuid", siteUUID)

	return nil
}

func (s *Service) RegenerateKeys(ctx context.Context, siteUUID string) (*Response, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "regenerating keys", "uuid", siteUUID)

	if _, err := uuid.Parse(siteUUID); err != nil {
		return nil, fmt.Errorf(errInvalidUUIDFormat, err)
	}

	pubKey, err := generateSecureKey(s.pubKeyLength)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate public key", "error", err)
		return nil, fmt.Errorf("%w: public key", ErrKeyGeneration)
	}

	secKey, err := generateSecureKey(s.secKeyLength)
	if err != nil {
		s.logger.ErrorContext(ctx, "failed to generate secret key", "error", err)
		return nil, fmt.Errorf("%w: secret key", ErrKeyGeneration)
	}

	site, err := s.repo.RegenerateKeys(ctx, siteUUID, pubKey, secKey)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, ErrSiteNotFound
		}
		s.logger.ErrorContext(ctx, "failed to regenerate keys",
			"error", err,
			"uuid", siteUUID,
		)
		return nil, fmt.Errorf("failed to regenerate keys: %w", err)
	}

	s.logger.InfoContext(ctx, "keys regenerated successfully", "uuid", siteUUID)

	return toResponse(site), nil
}

func (s *Service) validateCreateRequest(req CreateRequest) error {
	if err := validateName(req.Name); err != nil {
		return err
	}

	if err := validateDomain(req.Domain); err != nil {
		return err
	}

	return nil
}

func (s *Service) validateUpdateRequest(req UpdateRequest) error {
	if req.Name == nil && req.Domain == nil && req.TenantID == nil && req.IsActive == nil {
		return ErrEmptyUpdate
	}

	if req.Name != nil {
		if err := validateName(*req.Name); err != nil {
			return err
		}
	}

	if req.Domain != nil {
		if err := validateDomain(*req.Domain); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) checkDomainUniqueness(ctx context.Context, domain, excludeUUID string) error {
	normalized := normalizeDomain(domain)

	existing, err := s.repo.GetByDomain(ctx, normalized)
	if err != nil && !ent.IsNotFound(err) {
		s.logger.ErrorContext(ctx, "failed to check domain uniqueness",
			"error", err,
			"domain", normalized,
		)
		return fmt.Errorf("failed to check domain uniqueness: %w", err)
	}

	if existing != nil && existing.UUID != excludeUUID {
		s.logger.WarnContext(ctx, "domain already exists",
			"domain", normalized,
			"existing_uuid", existing.UUID,
		)
		return ErrDomainExists
	}

	return nil
}

func validateName(name string) error {
	trimmed := strings.TrimSpace(name)
	if len(trimmed) < MinNameLength || len(trimmed) > MaxNameLength {
		return fmt.Errorf("%w: must be between %d and %d characters",
			ErrInvalidName, MinNameLength, MaxNameLength)
	}
	return nil
}

func validateDomain(domain string) error {
	trimmed := strings.TrimSpace(domain)
	if len(trimmed) == 0 {
		return fmt.Errorf("%w: cannot be empty", ErrInvalidDomain)
	}

	if len(trimmed) > MaxDomainLength {
		return fmt.Errorf("%w: exceeds maximum length of %d characters",
			ErrInvalidDomain, MaxDomainLength)
	}

	normalized := normalizeDomain(trimmed)

	testURL := "http://" + normalized
	u, err := url.Parse(testURL)
	if err != nil || u.Host != normalized {
		return fmt.Errorf("%w: invalid format", ErrInvalidDomain)
	}

	if strings.Contains(normalized, " ") {
		return fmt.Errorf("%w: cannot contain spaces", ErrInvalidDomain)
	}

	if strings.HasPrefix(normalized, ".") || strings.HasSuffix(normalized, ".") {
		return fmt.Errorf("%w: cannot start or end with a dot", ErrInvalidDomain)
	}

	return nil
}

func normalizeDomain(domain string) string {
	normalized := strings.TrimSpace(domain)
	normalized = strings.ToLower(normalized)
	normalized = strings.TrimPrefix(normalized, "http://")
	normalized = strings.TrimPrefix(normalized, "https://")
	normalized = strings.TrimSuffix(normalized, "/")
	return normalized
}

func generateSecureKey(length int) (string, error) {
	if length <= 0 {
		return "", fmt.Errorf("length must be greater than 0")
	}

	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}

func toResponse(site *ent.Site) *Response {
	if site == nil {
		return nil
	}

	return &Response{
		UUID:      site.UUID,
		Name:      site.Name,
		Domain:    site.Domain,
		TenantID:  site.TenantID,
		PublicKey: site.PublicKey,
		IsActive:  site.IsActive,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}
}

func toFullResponse(site *ent.Site) *FullSiteResponse {
	if site == nil {
		return nil
	}

	return &FullSiteResponse{
		UUID:      site.UUID,
		Name:      site.Name,
		Domain:    site.Domain,
		TenantID:  site.TenantID,
		PublicKey: site.PublicKey,
		SecretKey: site.SecretKey,
		IsActive:  site.IsActive,
		CreatedAt: site.CreatedAt,
		UpdatedAt: site.UpdatedAt,
	}
}
func (c ServiceConfig) Validate() error {
	if c.PublicKeyLength <= 0 {
		return fmt.Errorf("public key length must be positive")
	}
	if c.SecretKeyLength <= 0 {
		return fmt.Errorf("secret key length must be positive")
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}
	return nil
}
