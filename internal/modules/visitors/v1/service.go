package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"tracking/internal/ent"
	devicesv1 "tracking/internal/modules/devices/v1"
	sitesv1 "tracking/internal/modules/sites/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"

	"github.com/google/uuid"
)

type Service struct {
	repo         Repository
	sitesService *sitesv1.Service
	deviceRepo   devicesv1.Repository
	sessionRepo  visitor_sessionsv1.Repository
	logger       *slog.Logger
	timeout      time.Duration
}

func NewService(repo Repository, sitesService *sitesv1.Service, deviceRepo devicesv1.Repository, sessionRepo visitor_sessionsv1.Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:         repo,
		sitesService: sitesService,
		deviceRepo:   deviceRepo,
		sessionRepo:  sessionRepo,
		logger:       logger,
		timeout:      5 * time.Second,
	}
}

type CreateVisitorRequest struct {
	// AnonymousID is extracted from cookie by the handler
	AnonymousID string `json:"-"`

	// PublicKey is extracted from header on the middleware level
	PublicKey string `json:"-"`

	// Device Info Extracted by the tracking SDK
	FingerprintHash string `json:"fingerprint_hash" binding:"required"`
	DeviceType      string `json:"device_type"`
	OS              string `json:"os"`
	Browser         string `json:"browser"`
	Language        string `json:"language"`

	// Geo location information extracted from ip
	AddressIP string `json:"address_ip"`
	Country   string `json:"country"`
	Region    string `json:"region"`
	City      string `json:"city"`
	Timezone  string `json:"timezone"`

	// Session Info
	SessionToken        string `json:"session_token"`
	EntryPage           string `json:"entry_page"`
	ExitPage            string `json:"exit_page"`
	AcquisitionSource   string `json:"acquisition_source"`
	AcquisitionMedium   string `json:"acquisition_medium"`
	AcquisitionCampaign string `json:"acquisition_campaign"`
	ReferrerHost        string `json:"referrer_host"`
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

/*
CreateVisitor godoc
*/
func (s *Service) CreateVisitor(ctx context.Context, req CreateVisitorRequest, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	s.logger.InfoContext(ctx, "tracking visitor", "anonymous_id", req.AnonymousID, "site_id", site.ID)

	v, err := s.handleVisitor(ctx, req)
	if err != nil {
		return nil, err
	}

	device, err := s.handleDevice(ctx, req, int64(v.ID))
	if err != nil {
		return nil, err
	}

	sessionID, err := s.handleSession(ctx, req, int64(v.ID), int64(device.ID))
	if err != nil {
		return nil, err
	}

	if sessionID != 0 {
		_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
			UUID:      uuid.New().String(),
			VisitorID: int64(v.ID),
			SiteID:    int64(site.ID),
			SessionID: sessionID,
		})
		if err != nil {
			s.logger.WarnContext(ctx, "failed to create site visitor link", "error", err)
		}
	}

	return toVisitorResponse(v), nil
}

func (s *Service) handleVisitor(ctx context.Context, req CreateVisitorRequest) (*ent.Visitor, error) {
	v, err := s.repo.GetByAnonymousID(ctx, req.AnonymousID)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, fmt.Errorf("failed to get visitor: %w", err)
		}

		v, err = s.repo.Create(ctx, CreateVisitorInput{
			UUID:        uuid.New().String(),
			AnonymousID: req.AnonymousID,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create visitor: %w", err)
		}
	} else {
		now := time.Now()
		v, err = s.repo.Update(ctx, v.UUID, UpdateVisitorInput{
			LastSeenAt: &now,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to update visitor: %w", err)
		}
	}
	return v, nil
}

func (s *Service) handleDevice(ctx context.Context, req CreateVisitorRequest, visitorID int64) (*ent.Device, error) {
	device, err := s.deviceRepo.GetByFingerprintHash(ctx, req.FingerprintHash)
	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to check device: %w", err)
	}

	if device == nil {
		d, err := s.deviceRepo.Create(ctx, devicesv1.CreateDeviceInput{
			UUID:            uuid.New().String(),
			VisitorID:       visitorID,
			FingerprintHash: req.FingerprintHash,
			DeviceType:      req.DeviceType,
			OS:              req.OS,
			AddressIP:       req.AddressIP,
			Country:         req.Country,
			Region:          req.Region,
			City:            req.City,
			Browser:         req.Browser,
			Language:        req.Language,
			Timezone:        req.Timezone,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create device: %w", err)
		}
		return d, nil
	}

	now := time.Now()
	d, err := s.deviceRepo.Update(ctx, device.UUID, devicesv1.UpdateDeviceInput{
		LastSeenAt: &now,
		AddressIP:  &req.AddressIP,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return d, nil
}

func (s *Service) handleSession(ctx context.Context, req CreateVisitorRequest, visitorID int64, deviceID int64) (int64, error) {
	if req.SessionToken == "" {
		return 0, nil
	}

	session, err := s.sessionRepo.GetBySessionToken(ctx, req.SessionToken)
	if err != nil && !ent.IsNotFound(err) {
		return 0, fmt.Errorf("failed to get session: %w", err)
	}

	if session == nil {
		sess, err := s.sessionRepo.Create(ctx, visitor_sessionsv1.CreateSessionInput{
			UUID:                uuid.New().String(),
			VisitorID:           visitorID,
			DeviceID:            deviceID,
			SessionToken:        req.SessionToken,
			StartedAt:           time.Now(),
			EntryPage:           req.EntryPage,
			AcquisitionSource:   req.AcquisitionSource,
			AcquisitionMedium:   req.AcquisitionMedium,
			AcquisitionCampaign: req.AcquisitionCampaign,
			ReferrerHost:        req.ReferrerHost,
		})
		if err != nil {
			return 0, fmt.Errorf("failed to create session: %w", err)
		}
		return int64(sess.ID), nil
	}

	now := time.Now()
	sess, err := s.sessionRepo.Update(ctx, session.UUID, visitor_sessionsv1.UpdateSessionInput{
		LastActivityAt: &now,
		ExitPage:       &req.ExitPage,
	})
	if err != nil {
		return 0, fmt.Errorf("failed to update session: %w", err)
	}
	return int64(sess.ID), nil
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
