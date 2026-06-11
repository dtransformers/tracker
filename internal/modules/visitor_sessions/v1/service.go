package v1

import (
	"context"
	"fmt"
	"log/slog"
	"time"
	"tracking/internal/config"
	"tracking/internal/ent"
	"tracking/internal/middlewares"
	"tracking/pkg/constants"
	custom_errors "tracking/pkg/errors"
)

type Service struct {
	repo    Repository
	logger  *slog.Logger
	appCfg  config.AppConfig
	timeout time.Duration
}

func NewService(repo Repository, logger *slog.Logger, appCfg config.AppConfig) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:    repo,
		logger:  logger,
		appCfg:  appCfg,
		timeout: constants.DefaultServiceTimeout,
	}
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

func (s *Service) GetSessionByTokenRaw(ctx context.Context, SessionToken string) (*ent.VisitorSession, error) {
	session, err := s.repo.GetByToken(ctx, SessionToken)
	if err != nil {
		return nil, err
	}
	return session, nil
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

func (s *Service) GetBySessionTokenRaw(ctx context.Context, token string) (*ent.VisitorSession, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.GetBySessionToken(ctx, token)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}

	if !sess.IsActive {
		return nil, nil
	}

	if sess.EndedAt != nil && time.Now().UTC().After(*sess.EndedAt) {
		isActive := false
		sess, err = s.repo.Update(ctx, sess.UUID, UpdateSessionInput{IsActive: &isActive})
		if err != nil {
			s.logger.WarnContext(ctx, "failed to update expired session status", "error", err)
			return nil, err
		}
		return nil, nil
	}

	return sess, nil
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
		IsActive:            s.IsActive,
		LastActivityAt:      s.LastActivityAt,
		EntryPage:           s.EntryPage,
		ExitPage:            s.ExitPage,
		AcquisitionSource:   s.AcquisitionSource,
		AcquisitionMedium:   s.AcquisitionMedium,
		AcquisitionCampaign: s.AcquisitionCampaign,
		CampaignTerm:        s.CampaignTerm,
		CampaignContent:     s.CampaignContent,
		CampaignID:          s.CampaignID,
		ReferrerHost:        s.ReferrerHost,
		CreatedAt:           s.CreatedAt,
	}
}

func (s *Service) CreateSession(ctx context.Context, input CreateSessionInput) (*ent.VisitorSession, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.repo.DeactivateVisitorSessions(ctx, input.VisitorID)
	if err != nil {
		s.logger.WarnContext(ctx, "failed to deactivate previous sessions", "error", err)
	}

	if input.EndedAt == nil && s.appCfg.SessionExpiration > 0 {
		endedAt := time.Now().UTC().Add(time.Duration(s.appCfg.SessionExpiration) * time.Second)
		input.EndedAt = &endedAt
	}
	input.IsActive = true

	sess, err := s.repo.Create(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return sess, nil
}

func (s *Service) UpdateSession(ctx context.Context, uuid string, input UpdateSessionInput) (*ent.VisitorSession, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.Update(ctx, uuid, input)
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}
	return sess, nil
}

func (s *Service) GetVisitorLastActiveSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.GetVisitorLastActiveSession(ctx, visitorID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return sess, nil
}

func (s *Service) GetVisitorLastSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.repo.GetVisitorLastSession(ctx, visitorID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, err
	}
	return sess, nil
}

func (s *Service) ValidateSessionDetails(ctx context.Context, VisitorUUID string, SessionToken string) (*middlewares.SessionDetails, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	session, err := s.GetBySessionTokenRaw(ctx, SessionToken)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, custom_errors.NewUnauthorizedError(constants.Version, "session not found, please clear cookies")
		}
		return nil, err
	}
	if session == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "session not found, please clear cookies")
	}

	err = s.repo.VerifySessionOwnership(ctx, session.VisitorID, VisitorUUID, session.DeviceID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, custom_errors.NewUnauthorizedError(constants.Version, "session does not belong to this visitor")
		}
		return nil, err
	}

	return &middlewares.SessionDetails{
		SessionID: int64(session.ID),
		VisitorID: session.VisitorID,
		DeviceID:  session.DeviceID,
	}, nil
}

func (s *Service) KeepAliveSession(ctx context.Context, sessionToken string) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	sess, err := s.GetSessionByTokenRaw(ctx, sessionToken)
	if err != nil {
		return err
	}
	if sess == nil {
		return custom_errors.NewNotFoundError(constants.Version, "session not found")
	}

	now := time.Now().UTC()
	var endedAt *time.Time
	if s.appCfg.SessionExpiration > 0 {
		t := now.Add(time.Duration(s.appCfg.SessionExpiration) * time.Second)
		endedAt = &t
	}

	_, err = s.UpdateSession(ctx, sess.UUID, UpdateSessionInput{
		LastActivityAt: &now,
		EndedAt:        endedAt,
	})
	return err
}
