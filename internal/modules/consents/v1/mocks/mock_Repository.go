package mocks

import (
	"context"
	"tracking/internal/ent"

	"github.com/stretchr/testify/mock"
)

type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) SetVisitorConsent(ctx context.Context, visitorID, siteID int64, consentType, status string) (*ent.VisitorConsent, error) {
	args := m.Called(ctx, visitorID, siteID, consentType, status)
	if rf, ok := args.Get(0).(func(context.Context, int64, int64, string, string) (*ent.VisitorConsent, error)); ok {
		return rf(ctx, visitorID, siteID, consentType, status)
	}
	var r0 *ent.VisitorConsent
	if args.Get(0) != nil {
		r0 = args.Get(0).(*ent.VisitorConsent)
	}
	return r0, args.Error(1)
}
