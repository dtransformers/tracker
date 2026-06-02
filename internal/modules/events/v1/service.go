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
	SessionToken    string `json:"session_token"`
	FingerprintHash string `json:"fingerprint_hash"`
	TrackingEvent
}

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest, site *ent.Site) (*ent.VisitorEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var visitorUUID string
	var sessionID string
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		if tracking.Cookies != nil {
			visitorUUID = tracking.Cookies["visitor_uuid"]
			sessionID = tracking.Cookies["session_id"]
		}
	}
	if sessionID == "" {
		return nil, fmt.Errorf("UnAuthoized Request")
	}
	details, err := s.visitorsSvc.GetSessionByUUID(ctx, visitorUUID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve session: %w", err)
	}

	ev, err := s.repo.CreateEvent(ctx, CreateEventInput{
		SessionID: int64(details.Session.ID),
		VisitorID: int64(details.Visitor.ID),
		DeviceID:  int64(details.Device.ID),
		Event:     &req.TrackingEvent,
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

	var visitorUUID string
	var sessionID string
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		if tracking.Cookies != nil {
			visitorUUID = tracking.Cookies["visitor_uuid"]
			sessionID = tracking.Cookies["session_id"]
		}
	}
	if sessionID == "" {
		return nil, fmt.Errorf("UnAuthoized Request")
	}
	details, err := s.visitorsSvc.GetSessionByUUID(ctx, visitorUUID, sessionID)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve session: %w", err)
	}

	pv, err := s.repo.CreatePageView(ctx, CreatePageViewInput{
		SessionID: int64(details.Session.ID),
		VisitorID: int64(details.Visitor.ID),
		DeviceID:  int64(details.Device.ID),
		URL:       req.URL,
		Title:     req.Title,
		Referrer:  req.Referrer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page view: %w", err)
	}

	return pv, nil
}
