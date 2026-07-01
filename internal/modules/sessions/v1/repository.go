package v1

import (
	"context"
	"fmt"
	"time"
	"tracking/internal/ent"
	"tracking/internal/ent/visitor"
	"tracking/internal/ent/visitorsession"
)

type Repository interface {
	Create(ctx context.Context, input CreateSessionInput) (*ent.VisitorSession, error)
	GetByUUID(ctx context.Context, uuid string) (*ent.VisitorSession, error)
	GetByToken(ctx context.Context, uuid string) (*ent.VisitorSession, error)


	GetByVisitorID(ctx context.Context, visitorID int64) ([]*ent.VisitorSession, error)
	GetVisitorLastActiveSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error)
	GetVisitorLastSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error)
	GetByDeviceID(ctx context.Context, deviceID int64) ([]*ent.VisitorSession, error)
	GetBySessionToken(ctx context.Context, token string) (*ent.VisitorSession, error)
	Update(ctx context.Context, uuid string, input UpdateSessionInput) (*ent.VisitorSession, error)
	DeactivateVisitorSessions(ctx context.Context, visitorID int64) error
	Delete(ctx context.Context, uuid string) error
	VerifySessionOwnership(ctx context.Context, visitorID int64, visitorUUID string, deviceID int64) error
}

type entRepository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	if client == nil {
		panic("ent client cannot be nil")
	}
	return &entRepository{client: client}
}

func (r *entRepository) Create(ctx context.Context, input CreateSessionInput) (*ent.VisitorSession, error) {
	query := r.client.VisitorSession.
		Create().
		SetUUID(input.UUID).
		SetVisitorID(input.VisitorID).
		SetDeviceID(input.DeviceID).
		SetSessionToken(input.SessionToken).
		SetStartedAt(input.StartedAt).
		SetEntryPage(input.EntryPage).
		SetAcquisitionSource(input.AcquisitionSource).
		SetAcquisitionMedium(input.AcquisitionMedium).
		SetAcquisitionCampaign(input.AcquisitionCampaign).
		SetCampaignTerm(input.CampaignTerm).
		SetCampaignContent(input.CampaignContent).
		SetCampaignID(input.CampaignID).
		SetReferrerHost(input.ReferrerHost).
		SetIsActive(input.IsActive)

	if input.EndedAt != nil {
		query.SetEndedAt(*input.EndedAt)
	}

	s, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create visitor session: %w", err)
	}

	return s, nil
}

func (r *entRepository) GetByUUID(ctx context.Context, uuid string) (*ent.VisitorSession, error) {
	s, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.UUID(uuid)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor session by uuid: %w", err)
	}

	return s, nil
}

func (r *entRepository) GetByToken(ctx context.Context, token string) (*ent.VisitorSession, error) {
	s, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.SessionToken(token)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor session by token: %w", err)
	}

	return s, nil
}

func (r *entRepository) GetByVisitorID(ctx context.Context, visitorID int64) ([]*ent.VisitorSession, error) {
	sessions, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.VisitorID(visitorID)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor sessions by visitor id: %w", err)
	}

	return sessions, nil
}

func (r *entRepository) GetVisitorLastActiveSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error) {
	sessions, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.VisitorID(visitorID)).
		Where(visitorsession.EndedAtGT(time.Now())).
		Where(visitorsession.IsActive(true)).
		Order(ent.Desc(visitorsession.FieldStartedAt)).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor session by visitor id: %w", err)
	}

	return sessions, nil
}

func (r *entRepository) GetVisitorLastSession(ctx context.Context, visitorID int64) (*ent.VisitorSession, error) {
	sessions, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.VisitorID(visitorID)).
		Order(ent.Desc(visitorsession.FieldStartedAt)).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor session by visitor id: %w", err)
	}

	return sessions, nil
}

func (r *entRepository) GetByDeviceID(ctx context.Context, deviceID int64) ([]*ent.VisitorSession, error) {
	sessions, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.DeviceID(deviceID)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor sessions by device id: %w", err)
	}

	return sessions, nil
}

func (r *entRepository) GetBySessionToken(ctx context.Context, token string) (*ent.VisitorSession, error) {
	s, err := r.client.VisitorSession.
		Query().
		Where(visitorsession.SessionToken(token)).
		Where(visitorsession.IsActive(true)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor session by token: %w", err)
	}

	return s, nil
}

func (r *entRepository) Update(ctx context.Context, uuid string, input UpdateSessionInput) (*ent.VisitorSession, error) {
	query := r.client.VisitorSession.
		Update().
		Where(visitorsession.UUID(uuid))

	if input.EndedAt != nil {
		query.SetEndedAt(*input.EndedAt)
	}
	if input.LastActivityAt != nil {
		query.SetLastActivityAt(*input.LastActivityAt)
	}
	if input.ExitPage != nil {
		query.SetExitPage(*input.ExitPage)
	}
	if input.IsActive != nil {
		query.SetIsActive(*input.IsActive)
	}

	affected, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update visitor session: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	return r.GetByUUID(ctx, uuid)
}

func (r *entRepository) DeactivateVisitorSessions(ctx context.Context, visitorID int64) error {
	_, err := r.client.VisitorSession.
		Update().
		Where(
			visitorsession.VisitorID(visitorID),
			visitorsession.IsActive(true),
		).
		SetIsActive(false).
		Save(ctx)

	if err != nil {
		return fmt.Errorf("failed to deactivate visitor sessions: %w", err)
	}

	return nil
}

func (r *entRepository) Delete(ctx context.Context, uuid string) error {
	affected, err := r.client.VisitorSession.
		Delete().
		Where(visitorsession.UUID(uuid)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete visitor session: %w", err)
	}

	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}

func (r *entRepository) VerifySessionOwnership(ctx context.Context, visitorID int64, visitorUUID string, deviceID int64) error {
	visitorExists, err := r.client.Visitor.
		Query().
		Where(
			visitor.ID(int(visitorID)),
			visitor.UUID(visitorUUID),
		).
		Exist(ctx)
	if err != nil {
		return fmt.Errorf("failed to query visitor: %w", err)
	}
	if !visitorExists {
		return &ent.NotFoundError{}
	}

	return nil
}
