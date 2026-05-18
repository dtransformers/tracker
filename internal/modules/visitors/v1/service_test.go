package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tracking/internal/ent"
	sitesv1 "tracking/internal/modules/sites/v1"
	sitesMocks "tracking/internal/modules/sites/v1/mocks"
	visitorsv1 "tracking/internal/modules/visitors/v1"
	visitorsMocks "tracking/internal/modules/visitors/v1/mocks"
)

func TestTrackVisitor_NewVisitor(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, nil)

	pubKey := "test-pub-key"
	anonID := "anon-123"
	sessionID := int64(456)

	req := visitorsv1.TrackVisitorRequest{
		PublicKey:   pubKey,
		AnonymousID: anonID,
		SessionID:   sessionID,
	}

	mockSite := &ent.Site{
		ID:        1,
		PublicKey: pubKey,
		IsActive:  true,
	}

	mockVisitor := &ent.Visitor{
		ID:          10,
		UUID:        uuid.New().String(),
		AnonymousID: anonID,
		FirstSeenAt: time.Now(),
		LastSeenAt:  time.Now(),
	}

	mockSitesRepo.On("GetByPublicKey", mock.Anything, pubKey).Return(mockSite, nil)
	mockVisitorsRepo.On("GetByAnonymousID", mock.Anything, anonID).Return(nil, &ent.NotFoundError{})
	mockVisitorsRepo.On("Create", mock.Anything, mock.Anything).Return(mockVisitor, nil)
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil)

	res, err := service.TrackVisitor(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, anonID, res.AnonymousID)
	mockSitesRepo.AssertExpectations(t)
	mockVisitorsRepo.AssertExpectations(t)
}

func TestTrackVisitor_ExistingVisitor(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, nil)

	pubKey := "test-pub-key"
	anonID := "anon-123"
	sessionID := int64(456)

	req := visitorsv1.TrackVisitorRequest{
		PublicKey:   pubKey,
		AnonymousID: anonID,
		SessionID:   sessionID,
	}

	mockSite := &ent.Site{
		ID:        1,
		PublicKey: pubKey,
		IsActive:  true,
	}

	mockVisitor := &ent.Visitor{
		ID:          10,
		UUID:        uuid.New().String(),
		AnonymousID: anonID,
		FirstSeenAt: time.Now().Add(-1 * time.Hour),
		LastSeenAt:  time.Now().Add(-1 * time.Hour),
	}

	updatedVisitor := &ent.Visitor{
		ID:          10,
		UUID:        mockVisitor.UUID,
		AnonymousID: anonID,
		FirstSeenAt: mockVisitor.FirstSeenAt,
		LastSeenAt:  time.Now(),
	}

	mockSitesRepo.On("GetByPublicKey", mock.Anything, pubKey).Return(mockSite, nil)
	mockVisitorsRepo.On("GetByAnonymousID", mock.Anything, anonID).Return(mockVisitor, nil)
	mockVisitorsRepo.On("Update", mock.Anything, mockVisitor.UUID, mock.Anything).Return(updatedVisitor, nil)
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil)

	res, err := service.TrackVisitor(context.Background(), req)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, anonID, res.AnonymousID)
	mockSitesRepo.AssertExpectations(t)
	mockVisitorsRepo.AssertExpectations(t)
}

func TestTrackVisitor_InvalidSiteKey(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, nil)

	pubKey := "invalid-pub-key"
	req := visitorsv1.TrackVisitorRequest{
		PublicKey:   pubKey,
		AnonymousID: "anon-123",
		SessionID:   123,
	}

	mockSitesRepo.On("GetByPublicKey", mock.Anything, pubKey).Return(nil, &ent.NotFoundError{})

	res, err := service.TrackVisitor(context.Background(), req)

	assert.Error(t, err)
	assert.Nil(t, res)
	mockSitesRepo.AssertExpectations(t)
}
