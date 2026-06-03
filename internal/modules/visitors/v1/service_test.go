package v1_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"tracking/internal/config"
	"tracking/internal/ent"
	consentsv1 "tracking/internal/modules/consents/v1"
	consentsMocks "tracking/internal/modules/consents/v1/mocks"
	devicesv1 "tracking/internal/modules/devices/v1"
	devicesMocks "tracking/internal/modules/devices/v1/mocks"
	sitesv1 "tracking/internal/modules/sites/v1"
	sitesMocks "tracking/internal/modules/sites/v1/mocks"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
	visitor_sessionsMocks "tracking/internal/modules/visitor_sessions/v1/mocks"
	visitorsv1 "tracking/internal/modules/visitors/v1"
	visitorsMocks "tracking/internal/modules/visitors/v1/mocks"
)

func TestTrackVisitor_NewVisitor(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	pubKey := "test-pub-key"
	anonID := "anon-123"

	req := visitorsv1.CreateVisitorRequest{
		PublicKey:           pubKey,
		SessionToken:        "session-token",
		DeviceType:          "desktop",
		OS:                  "linux",
		AddressIP:           "[IP_ADDRESS]",
		Country:             "US",
		Region:              "CA",
		City:                "San Francisco",
		Browser:             "chrome",
		Language:            "en-US",
		Timezone:            "America/New_York",
		EntryPage:           "https://example.com",
		ExitPage:            "https://example.com",
		AcquisitionSource:   "google",
		AcquisitionMedium:   "cpc",
		AcquisitionCampaign: "campaign1",
		ReferrerHost:        "google.com",
		FingerprintHash:     "fingerprint-hash",
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

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            uuid.New().String(),
		FingerprintHash: "fingerprint-hash",
	}

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         uuid.New().String(),
		SessionToken: "session-token",
	}

	mockVisitorsRepo.On("GetByAnonymousID", mock.Anything, anonID).Return(nil, &ent.NotFoundError{}).Maybe()
	mockVisitorsRepo.On("Create", mock.Anything, mock.Anything).Return(mockVisitor, nil)
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil)

	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, "fingerprint-hash").Return(nil, &ent.NotFoundError{})
	mockDevicesRepo.On("Create", mock.Anything, mock.Anything).Return(mockDevice, nil)

	mockSessionsRepo.On("GetBySessionToken", mock.Anything, "session-token").Return(nil, &ent.NotFoundError{})
	mockSessionsRepo.On("DeactivateVisitorSessions", mock.Anything, mock.Anything).Return(nil).Maybe()
	mockSessionsRepo.On("Create", mock.Anything, mock.Anything).Return(mockSession, nil)
	mockConsentsRepo.On("SetVisitorConsent", mock.Anything, mock.Anything, mock.Anything, "none", "pending").Return(&ent.VisitorConsent{}, nil)

	res, err := service.CreateSession(context.Background(), req, mockSite)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockVisitorsRepo.AssertExpectations(t)
	mockDevicesRepo.AssertExpectations(t)
	mockSessionsRepo.AssertExpectations(t)
	mockConsentsRepo.AssertExpectations(t)
}

func TestTrackVisitor_ExistingVisitor(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	pubKey := "test-pub-key"
	anonID := "anon-123"
	req := visitorsv1.CreateVisitorRequest{
		PublicKey:           pubKey,
		SessionToken:        "session-token",
		DeviceType:          "desktop",
		OS:                  "linux",
		AddressIP:           "[IP_ADDRESS]",
		Country:             "US",
		Region:              "CA",
		City:                "San Francisco",
		Browser:             "chrome",
		Language:            "en-US",
		Timezone:            "America/New_York",
		EntryPage:           "https://example.com",
		ExitPage:            "https://example.com",
		AcquisitionSource:   "google",
		AcquisitionMedium:   "cpc",
		AcquisitionCampaign: "campaign1",
		ReferrerHost:        "google.com",
		FingerprintHash:     "fingerprint-hash",
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

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            "device-uuid",
		FingerprintHash: "fingerprint-hash",
		VisitorID:       10,
	}

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         "session-uuid",
		SessionToken: "session-token",
	}

	mockVisitorsRepo.On("GetByID", mock.Anything, int64(10)).Return(mockVisitor, nil)
	mockVisitorsRepo.On("Update", mock.Anything, mockVisitor.UUID, mock.Anything).Return(updatedVisitor, nil)
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil).Maybe()

	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, "fingerprint-hash").Return(mockDevice, nil)
	mockDevicesRepo.On("Update", mock.Anything, "device-uuid", mock.Anything).Return(mockDevice, nil)

	mockSessionsRepo.On("GetBySessionToken", mock.Anything, "session-token").Return(mockSession, nil)
	mockSessionsRepo.On("DeactivateVisitorSessions", mock.Anything, mock.Anything).Return(nil).Maybe()
	mockSessionsRepo.On("Update", mock.Anything, "session-uuid", mock.Anything).Return(mockSession, nil)

	res, err := service.CreateSession(context.Background(), req, mockSite)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockVisitorsRepo.AssertExpectations(t)
	mockDevicesRepo.AssertExpectations(t)
	mockSessionsRepo.AssertExpectations(t)
	mockConsentsRepo.AssertExpectations(t)
}

func TestTrackVisitor_GeoIP_SkipLocal(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	pubKey := "test-pub-key"
	anonID := "anon-123"

	req := visitorsv1.CreateVisitorRequest{
		PublicKey:       pubKey,
		AddressIP:       "127.0.0.1",
		FingerprintHash: "fingerprint-hash",
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

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            uuid.New().String(),
		FingerprintHash: "fingerprint-hash",
	}

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         uuid.New().String(),
		SessionToken: "session-token",
	}

	mockVisitorsRepo.On("GetByAnonymousID", mock.Anything, anonID).Return(nil, &ent.NotFoundError{}).Maybe()
	mockVisitorsRepo.On("Create", mock.Anything, mock.Anything).Return(mockVisitor, nil)
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil)

	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, "fingerprint-hash").Return(nil, &ent.NotFoundError{})
	mockDevicesRepo.On("Create", mock.Anything, mock.MatchedBy(func(in devicesv1.CreateDeviceInput) bool {
		return in.Country == "" && in.Region == "" && in.City == "" && in.Timezone == "" && in.AddressIP == "127.0.0.1"
	})).Return(mockDevice, nil)

	mockSessionsRepo.On("GetBySessionToken", mock.Anything, mock.AnythingOfType("string")).Return(nil, &ent.NotFoundError{})
	mockSessionsRepo.On("DeactivateVisitorSessions", mock.Anything, mock.Anything).Return(nil).Maybe()
	mockSessionsRepo.On("Create", mock.Anything, mock.Anything).Return(mockSession, nil)
	mockConsentsRepo.On("SetVisitorConsent", mock.Anything, mock.Anything, mock.Anything, "none", "pending").Return(&ent.VisitorConsent{}, nil)

	res, err := service.CreateSession(context.Background(), req, mockSite)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	mockVisitorsRepo.AssertExpectations(t)
	mockDevicesRepo.AssertExpectations(t)
	mockSessionsRepo.AssertExpectations(t)
	mockConsentsRepo.AssertExpectations(t)
}
