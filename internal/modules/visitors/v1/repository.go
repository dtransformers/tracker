package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"
	"tracking/internal/ent/visitor"
	"tracking/internal/ent/visitorconsent"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, input CreateSessionInput) (*ent.Visitor, error)
	GetByID(ctx context.Context, id int64) (*ent.Visitor, error)
	GetByUUID(ctx context.Context, uuid string) (*ent.Visitor, error)
	GetByAnonymousID(ctx context.Context, anonymousID string) (*ent.Visitor, error)
	Update(ctx context.Context, uuid string, input UpdateVisitorInput) (*ent.Visitor, error)
	Delete(ctx context.Context, uuid string) error

	CreateSiteVisitor(ctx context.Context, input CreateSiteVisitorInput) (*ent.SiteVisitor, error)
	SetVisitorConsent(ctx context.Context, visitorID, siteID int64, consentType, status string) (*ent.VisitorConsent, error)
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

func (r *entRepository) Create(ctx context.Context, input CreateSessionInput) (*ent.Visitor, error) {
	query := r.client.Visitor.
		Create().
		SetUUID(input.UUID).
		SetAnonymousID(input.AnonymousID).
		SetIsBot(input.IsBot).
		SetIsVerifiedHuman(input.IsVerifiedHuman)

	if input.PrimaryUserID != nil {
		query.SetPrimaryUserID(*input.PrimaryUserID)
	}

	v, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create visitor: %w", err)
	}

	return v, nil
}

func (r *entRepository) GetByID(ctx context.Context, id int64) (*ent.Visitor, error) {
	v, err := r.client.Visitor.Get(ctx, int(id))
	if err != nil {
		return nil, fmt.Errorf("failed to get visitor by id: %w", err)
	}
	return v, nil
}

func (r *entRepository) GetByUUID(ctx context.Context, uuid string) (*ent.Visitor, error) {
	v, err := r.client.Visitor.
		Query().
		Where(visitor.UUID(uuid)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor by uuid: %w", err)
	}

	return v, nil
}

func (r *entRepository) GetByAnonymousID(ctx context.Context, anonymousID string) (*ent.Visitor, error) {
	v, err := r.client.Visitor.
		Query().
		Where(visitor.AnonymousID(anonymousID)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get visitor by anonymous id: %w", err)
	}

	return v, nil
}

func (r *entRepository) Update(ctx context.Context, uuid string, input UpdateVisitorInput) (*ent.Visitor, error) {
	query := r.client.Visitor.
		Update().
		Where(visitor.UUID(uuid))

	if input.PrimaryUserID != nil {
		query.SetPrimaryUserID(*input.PrimaryUserID)
	}
	if input.LastSeenAt != nil {
		query.SetLastSeenAt(*input.LastSeenAt)
	}
	if input.IsBot != nil {
		query.SetIsBot(*input.IsBot)
	}
	if input.IsVerifiedHuman != nil {
		query.SetIsVerifiedHuman(*input.IsVerifiedHuman)
	}

	affected, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update visitor: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	return r.GetByUUID(ctx, uuid)
}

func (r *entRepository) Delete(ctx context.Context, uuid string) error {
	affected, err := r.client.Visitor.
		Delete().
		Where(visitor.UUID(uuid)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete visitor: %w", err)
	}

	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}

func (r *entRepository) CreateSiteVisitor(ctx context.Context, input CreateSiteVisitorInput) (*ent.SiteVisitor, error) {
	sv, err := r.client.SiteVisitor.
		Create().
		SetUUID(input.UUID).
		SetVisitorID(input.VisitorID).
		SetSessionID(input.SessionID).
		SetSiteID(input.SiteID).
		Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create site visitor: %w", err)
	}

	return sv, nil
}

func (r *entRepository) SetVisitorConsent(ctx context.Context, visitorID, siteID int64, consentType, status string) (*ent.VisitorConsent, error) {
	vc, err := r.client.VisitorConsent.Query().
		Where(
			visitorconsent.VisitorID(visitorID),
			visitorconsent.SiteID(siteID),
		).First(ctx)

	if err != nil && !ent.IsNotFound(err) {
		return nil, fmt.Errorf("failed to query visitor consent: %w", err)
	}

	if vc != nil {
		return vc.Update().
			SetConsentType(consentType).
			SetStatus(status).
			Save(ctx)
	}

	return r.client.VisitorConsent.Create().
		SetUUID(uuid.New().String()).
		SetVisitorID(visitorID).
		SetSiteID(siteID).
		SetConsentType(consentType).
		SetStatus(status).
		Save(ctx)
}
