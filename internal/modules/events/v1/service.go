package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	visitorsv1 "tracking/internal/modules/visitors/v1"
)

type Service struct {
	repo        Repository
	visitorsSvc *visitorsv1.Service
	logger      *slog.Logger
	timeout     time.Duration
}

func NewService(repo Repository, visitorsSvc *visitorsv1.Service, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:        repo,
		visitorsSvc: visitorsSvc,
		logger:      logger,
		timeout:     5 * time.Second,
	}
}

type CreateEventRequest struct {
	SessionToken    string                 `json:"session_token"`
	FingerprintHash string                 `json:"fingerprint_hash"`
	EventName       string                 `json:"event_name" binding:"required"`
	EventProperties map[string]interface{} `json:"event_properties"`
}

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest, site *ent.Site) (*ent.VisitorEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var anonID string
	var sessionID string
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		if tracking.Cookies != nil {
			anonID = tracking.Cookies["anonymous_id"]
			sessionID = tracking.Cookies["session_id"]
		}
	}
	if req.SessionToken == "" {
		req.SessionToken = sessionID
	}
	// geo is just empty since we don't extract from event payload usually, but we need to pass something
	geo := visitorsv1.CreateVisitorRequest{
		VisitorUUID:  anonID,
		SessionToken: req.SessionToken,
	}

	session, err := s.visitorsSvc.ResolveSession(ctx, req.SessionToken, req.FingerprintHash, geo, site)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve session: %w", err)
	}

	ev, err := s.repo.CreateEvent(ctx, CreateEventInput{
		SessionID:       int64(session.ID),
		VisitorID:       session.VisitorID,
		DeviceID:        session.DeviceID,
		EventName:       req.EventName,
		EventProperties: req.EventProperties,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return ev, nil
}

type CreatePageViewRequest struct {
	SessionToken    string `json:"session_token"`
	FingerprintHash string `json:"fingerprint_hash"`
	URL             string `json:"url" binding:"required"`
	Title           string `json:"title"`
	Referrer        string `json:"referrer"`
}

func (s *Service) CreatePageView(ctx context.Context, req CreatePageViewRequest, site *ent.Site) (*ent.PageView, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var anonID string
	var sessionID string
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		if tracking.Cookies != nil {
			anonID = tracking.Cookies["anonymous_id"]
			sessionID = tracking.Cookies["session_id"]
		}
	}
	if req.SessionToken == "" {
		req.SessionToken = sessionID
	}
	geo := visitorsv1.CreateVisitorRequest{
		VisitorUUID:  anonID,
		SessionToken: req.SessionToken,
		EntryPage:    req.URL,
	}

	session, err := s.visitorsSvc.ResolveSession(ctx, req.SessionToken, req.FingerprintHash, geo, site)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve session: %w", err)
	}

	pv, err := s.repo.CreatePageView(ctx, CreatePageViewInput{
		SessionID: int64(session.ID),
		VisitorID: session.VisitorID,
		DeviceID:  session.DeviceID,
		URL:       req.URL,
		Title:     req.Title,
		Referrer:  req.Referrer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page view: %w", err)
	}

	return pv, nil
}
