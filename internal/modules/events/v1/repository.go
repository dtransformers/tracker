package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"

	"github.com/google/uuid"
)

type Repository interface {
	CreateEvent(ctx context.Context, input CreateEventInput) (*ent.VisitorEvent, error)
	CreatePageView(ctx context.Context, input CreatePageViewInput) (*ent.PageView, error)
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

type CreateEventInput struct {
	SessionID       int64
	VisitorID       int64
	DeviceID        int64
	EventName       string
	EventProperties map[string]interface{}
}

type CreatePageViewInput struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
	URL       string
	Title     string
	Referrer  string
}

func (r *entRepository) CreateEvent(ctx context.Context, input CreateEventInput) (*ent.VisitorEvent, error) {
	ev, err := r.client.VisitorEvent.
		Create().
		SetUUID(uuid.New().String()).
		SetSessionID(input.SessionID).
		SetVisitorID(input.VisitorID).
		SetDeviceID(input.DeviceID).
		SetEventName(input.EventName).
		SetEventProperties(input.EventProperties).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create visitor event: %w", err)
	}
	return ev, nil
}

func (r *entRepository) CreatePageView(ctx context.Context, input CreatePageViewInput) (*ent.PageView, error) {
	pv, err := r.client.PageView.
		Create().
		SetUUID(uuid.New().String()).
		SetSessionID(input.SessionID).
		SetVisitorID(input.VisitorID).
		SetDeviceID(input.DeviceID).
		SetURL(input.URL).
		SetTitle(input.Title).
		SetReferrer(input.Referrer).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create page view: %w", err)
	}
	return pv, nil
}
