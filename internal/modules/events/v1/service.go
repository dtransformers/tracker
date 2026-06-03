package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	"tracking/pkg/constants"
)

type Service struct {
	repo        Repository
	logger      *slog.Logger
	timeout     time.Duration
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:        repo,
		logger:      logger,
		timeout:     constants.DefaultServiceTimeout,
	}
}


func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest, site *ent.Site) (*ent.VisitorEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	details, ok := middlewares.GetSessionDetails(ctx)
	if !ok || details == nil {
		return nil, fmt.Errorf("failed to retrieve session details from context")
	}

	ev, err := s.repo.CreateEvent(ctx, CreateEventInput{
		SessionID: details.SessionID,
		VisitorID: details.VisitorID,
		DeviceID:  details.DeviceID,
		Event:     &req.TrackingEvent,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create event: %w", err)
	}

	return ev, nil
}


func (s *Service) CreatePageView(ctx context.Context, req CreatePageViewRequest, site *ent.Site) (*ent.PageView, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	details, ok := middlewares.GetSessionDetails(ctx)
	if !ok || details == nil {
		return nil, fmt.Errorf("failed to retrieve session details from context")
	}

	pv, err := s.repo.CreatePageView(ctx, CreatePageViewInput{
		SessionID: details.SessionID,
		VisitorID: details.VisitorID,
		DeviceID:  details.DeviceID,
		URL:       req.URL,
		Title:     req.Title,
		Referrer:  req.Referrer,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create page view: %w", err)
	}

	return pv, nil
}
