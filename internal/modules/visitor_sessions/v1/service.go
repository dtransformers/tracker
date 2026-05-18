package v1

import (
	"context"
	"log/slog"
	"time"
	"tracking/internal/ent"
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
		timeout: 5 * time.Second,
	}
}

type SessionResponse struct {
	UUID                string     `json:"uuid"`
	VisitorID           int64      `json:"visitor_id"`
	DeviceID            int64      `json:"device_id"`
	SessionToken        string     `json:"session_token"`
	StartedAt           time.Time  `json:"started_at"`
	EndedAt             *time.Time `json:"ended_at,omitempty"`
	LastActivityAt      time.Time  `json:"last_activity_at"`
	EntryPage           string     `json:"entry_page"`
	ExitPage            string     `json:"exit_page"`
	AcquisitionSource   string     `json:"acquisition_source"`
	AcquisitionMedium   string     `json:"acquisition_medium"`
	AcquisitionCampaign string     `json:"acquisition_campaign"`
	ReferrerHost        string     `json:"referrer_host"`
	CreatedAt           time.Time  `json:"created_at"`
}

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*SessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return toSessionResponse(sess), nil
}

func (s *Service) GetByVisitorID(ctx context.Context, visitorID int64) ([]*SessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sessions, err := s.repo.GetByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, err
	}

	res := make([]*SessionResponse, len(sessions))
	for i, sess := range sessions {
		res[i] = toSessionResponse(sess)
	}
	return res, nil
}

func (s *Service) GetByDeviceID(ctx context.Context, deviceID int64) ([]*SessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sessions, err := s.repo.GetByDeviceID(ctx, deviceID)
	if err != nil {
		return nil, err
	}

	res := make([]*SessionResponse, len(sessions))
	for i, sess := range sessions {
		res[i] = toSessionResponse(sess)
	}
	return res, nil
}

func (s *Service) GetBySessionToken(ctx context.Context, token string) (*SessionResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.GetBySessionToken(ctx, token)
	if err != nil {
		return nil, err
	}
	return toSessionResponse(sess), nil
}

func toSessionResponse(s *ent.VisitorSession) *SessionResponse {
	if s == nil {
		return nil
	}
	var endedAt *time.Time
	// Note: s.EndedAt might be a time.Time, check if it's optional in schema. 
	// The schema says field.Time("ended_at").Optional().Nillable(). So it's a pointer.
	endedAt = s.EndedAt

	return &SessionResponse{
		UUID:                s.UUID,
		VisitorID:           s.VisitorID,
		DeviceID:            s.DeviceID,
		SessionToken:        s.SessionToken,
		StartedAt:           s.StartedAt,
		EndedAt:             endedAt,
		LastActivityAt:      s.LastActivityAt,
		EntryPage:           s.EntryPage,
		ExitPage:            s.ExitPage,
		AcquisitionSource:   s.AcquisitionSource,
		AcquisitionMedium:   s.AcquisitionMedium,
		AcquisitionCampaign: s.AcquisitionCampaign,
		ReferrerHost:        s.ReferrerHost,
		CreatedAt:           s.CreatedAt,
	}
}
