package v1

import (
	"context"
	"fmt"
	"time"
	"tracking/internal/ent"
	"tracking/internal/ent/visitorsession"
)

type Repository interface {
	Create(ctx context.Context, input CreateSessionInput) (*ent.VisitorSession, error)
	GetByUUID(ctx context.Context, uuid string) (*ent.VisitorSession, error)
	GetByVisitorID(ctx context.Context, visitorID int64) ([]*ent.VisitorSession, error)
	GetByDeviceID(ctx context.Context, deviceID int64) ([]*ent.VisitorSession, error)
	GetBySessionToken(ctx context.Context, token string) (*ent.VisitorSession, error)
	Update(ctx context.Context, uuid string, input UpdateSessionInput) (*ent.VisitorSession, error)
	Delete(ctx context.Context, uuid string) error
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

type CreateSessionInput struct {
	UUID                string
	VisitorID           int64
	DeviceID            int64
	SessionToken        string
	StartedAt           time.Time
	EntryPage           string
	AcquisitionSource   string
	AcquisitionMedium   string
	AcquisitionCampaign string
	ReferrerHost        string
}

type UpdateSessionInput struct {
	EndedAt        *time.Time
	LastActivityAt *time.Time
	ExitPage       *string
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
		SetReferrerHost(input.ReferrerHost)

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

	affected, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update visitor session: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	return r.GetByUUID(ctx, uuid)
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
