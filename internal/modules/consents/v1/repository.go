package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"
	"tracking/internal/ent/visitorconsent"

	"github.com/google/uuid"
)

type Repository interface {
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

func (r *entRepository) SetVisitorConsent(ctx context.Context, visitorID, siteID int64, consentType, status string) (*ent.VisitorConsent, error) {
	vc, err := r.client.VisitorConsent.Query().
		Where(
			visitorconsent.VisitorID(visitorID),
			visitorconsent.SiteID(siteID),
			visitorconsent.ConsentType(consentType),
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
