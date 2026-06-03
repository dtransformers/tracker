package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	"tracking/pkg/constants"
	custom_errors "tracking/pkg/errors"
)

type Service struct {
	repo    Repository
	logger  *slog.Logger
	timeout time.Duration
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:    repo,
		logger:  logger,
		timeout: constants.DefaultServiceTimeout,
	}
}

func (s *Service) UpdateVisitorConsent(ctx context.Context, req UpdateConsentRequest, site *ent.Site) (*UpdateConsentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if !req.ConsentType.Valid() {
		return nil, custom_errors.NewBadRequestError(constants.Version, fmt.Sprintf("invalid consent type: %s", req.ConsentType))
	}
	if !req.Status.Valid() {
		return nil, custom_errors.NewBadRequestError(constants.Version, fmt.Sprintf("invalid consent status: %s", req.Status))
	}

	details, ok := middlewares.GetSessionDetails(ctx)
	if !ok || details == nil {
		return nil, fmt.Errorf("failed to retrieve session details from context")
	}

	consent, err := s.repo.SetVisitorConsent(ctx, details.VisitorID, int64(site.ID), string(req.ConsentType), string(req.Status))
	if err != nil {
		return nil, fmt.Errorf("failed to update consent: %w", err)
	}

	return &UpdateConsentResponse{
		Success: true,
		Message: "Consent updated successfully",
		Data:    consent,
	}, nil
}

func (s *Service) SetDefaultConsent(ctx context.Context, visitorID, siteID int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.repo.SetVisitorConsent(ctx, visitorID, siteID, string(ConsentTypeNone), string(ConsentStatusPending))
	return err
}
