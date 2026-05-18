package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"tracking/internal/ent"
	sitesv1 "tracking/internal/modules/sites/v1"
	"github.com/google/uuid"
)

type Service struct {
	repo         Repository
	sitesService *sitesv1.Service
	logger       *slog.Logger
	timeout      time.Duration
}

func NewService(repo Repository, sitesService *sitesv1.Service, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:         repo,
		sitesService: sitesService,
		logger:       logger,
		timeout:      5 * time.Second,
	}
}

type TrackVisitorRequest struct {
	PublicKey   string `json:"public_key" binding:"required"`
	AnonymousID string `json:"anonymous_id" binding:"required"`
	SessionID   int64  `json:"session_id" binding:"required"`
}

type VisitorResponse struct {
	UUID            string    `json:"uuid"`
	AnonymousID     string    `json:"anonymous_id"`
	PrimaryUserID   *string   `json:"primary_user_id,omitempty"`
	FirstSeenAt     time.Time `json:"first_seen_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	IsBot           bool      `json:"is_bot"`
	IsVerifiedHuman bool      `json:"is_verified_human"`
}

func (s *Service) TrackVisitor(ctx context.Context, req TrackVisitorRequest) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "tracking visitor", "anonymous_id", req.AnonymousID, "public_key", req.PublicKey)

	// 1. Verify site key
	ok, err := s.sitesService.VerifySiteKey(ctx, req.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to verify site key: %w", err)
	}
	if !ok {
		return nil, fmt.Errorf("invalid or inactive site key")
	}

	// 2. Get site to get site_id
	site, err := s.sitesService.GetSiteKey(ctx, req.PublicKey)
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	// 3. Check if visitor exists by anonymous ID
	v, err := s.repo.GetByAnonymousID(ctx, req.AnonymousID)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get visitor: %w", err)
		}

		// Visitor not found, create new one
		v, err = s.repo.Create(ctx, CreateVisitorInput{
			UUID:        uuid.New().String(),
			AnonymousID: req.AnonymousID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create visitor: %w", err)
		}
	} else {
		// Update last seen
		now := time.Now()
		v, err = s.repo.Update(ctx, v.UUID, UpdateVisitorInput{
			LastSeenAt: &now,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update visitor: %w", err)
		}
	}

	// 4. Create SiteVisitor link
	_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
		UUID:      uuid.New().String(),
		VisitorID: int64(v.ID),
		SiteID:    int64(site.ID),
		SessionID: req.SessionID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create site visitor link: %w", err)
	}

	return toVisitorResponse(v), nil
}

func toVisitorResponse(v *ent.Visitor) *VisitorResponse {
	if v == nil {
		return nil
	}
	return &VisitorResponse{
		UUID:            v.UUID,
		AnonymousID:     v.AnonymousID,
		PrimaryUserID:   v.PrimaryUserID,
		FirstSeenAt:     v.FirstSeenAt,
		LastSeenAt:      v.LastSeenAt,
		IsBot:           v.IsBot,
		IsVerifiedHuman: v.IsVerifiedHuman,
	}
}
