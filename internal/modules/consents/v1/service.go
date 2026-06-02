package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
)

type Service struct {
	repo       Repository
	sessionSvc *visitor_sessionsv1.Service
	logger     *slog.Logger
	timeout    time.Duration
}

func NewService(repo Repository, sessionSvc *visitor_sessionsv1.Service, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:       repo,
		sessionSvc: sessionSvc,
		logger:     logger,
		timeout:    5 * time.Second,
	}
}

type UpdateConsentRequest struct {
	ConsentType string `json:"consent_type" binding:"required"`
	Status      string `json:"status" binding:"required"`
}

type UpdateConsentResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *Service) UpdateVisitorConsent(ctx context.Context, req UpdateConsentRequest, site *ent.Site) (*UpdateConsentResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var sessionToken string
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		if tracking.Cookies != nil {
			sessionToken = tracking.Cookies["session_id"]
		}
	}

	if sessionToken == "" {
		return nil, fmt.Errorf("missing session token in cookies")
	}

	session, err := s.sessionSvc.GetBySessionToken(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, fmt.Errorf("session not found")
	}

	_, err = s.repo.SetVisitorConsent(ctx, session.VisitorID, int64(site.ID), req.ConsentType, req.Status)
	if err != nil {
		return nil, fmt.Errorf("failed to update consent: %w", err)
	}

	return &UpdateConsentResponse{
		Success: true,
		Message: "Consent updated successfully",
	}, nil
}

func (s *Service) SetDefaultConsent(ctx context.Context, visitorID, siteID int64) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	_, err := s.repo.SetVisitorConsent(ctx, visitorID, siteID, "none", "pending")
	return err
}
