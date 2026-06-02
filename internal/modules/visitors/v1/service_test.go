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
	mockSessionsRepo.On("Create", mock.Anything, mock.Anything).Return(mockSession, nil)
	mockConsentsRepo.On("SetVisitorConsent", mock.Anything, mock.Anything, mock.Anything, "none", "pending").Return(&ent.VisitorConsent{}, nil)

	res, err := service.CreateSession(context.Background(), req, mockSite)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, anonID, res.AnonymousID)
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
	mockSessionsRepo.On("Update", mock.Anything, "session-uuid", mock.Anything).Return(mockSession, nil)

	res, err := service.CreateSession(context.Background(), req, mockSite)

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, anonID, res.AnonymousID)
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

func TestCreateEvent_ExistingSession(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	sessionToken := "sess-123"
	fingerprint := "fp-123"

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         "session-uuid",
		VisitorID:    10,
		DeviceID:     20,
		SessionToken: sessionToken,
	}

	mockVisitor := &ent.Visitor{
		ID:   10,
		UUID: "visitor-uuid",
	}

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            "device-uuid",
		FingerprintHash: fingerprint,
	}

	mockEvent := &ent.VisitorEvent{
		UUID:            "event-uuid",
		EventName:       "click",
		EventProperties: map[string]interface{}{"button": "blue"},
		CreatedAt:       time.Now(),
	}

	mockSessionsRepo.On("GetBySessionToken", mock.Anything, sessionToken).Return(mockSession, nil)
	mockVisitorsRepo.On("GetByID", mock.Anything, int64(10)).Return(mockVisitor, nil)
	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, fingerprint).Return(mockDevice, nil)
	mockVisitorsRepo.On("CreateEvent", mock.Anything, mock.Anything).Return(mockEvent, nil)

	req := visitorsv1.CreateEventRequest{
		SessionToken:    sessionToken,
		FingerprintHash: fingerprint,
		EventName:       "click",
		EventProperties: map[string]interface{}{"button": "blue"},
	}

	res, err := service.CreateEvent(context.Background(), req, &ent.Site{})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "event-uuid", res.UUID)
	assert.Equal(t, "session-uuid", res.SessionUUID)
	assert.Equal(t, "visitor-uuid", res.VisitorUUID)
	assert.Equal(t, "device-uuid", res.DeviceUUID)
	assert.Equal(t, "click", res.EventName)
}

func TestCreatePageView_ExistingSession(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	sessionToken := "sess-123"
	fingerprint := "fp-123"

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         "session-uuid",
		VisitorID:    10,
		DeviceID:     20,
		SessionToken: sessionToken,
	}

	mockVisitor := &ent.Visitor{
		ID:   10,
		UUID: "visitor-uuid",
	}

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            "device-uuid",
		FingerprintHash: fingerprint,
	}

	mockPV := &ent.PageView{
		UUID:      "pv-uuid",
		URL:       "https://example.com/home",
		Title:     "Home Page",
		Referrer:  "https://google.com",
		CreatedAt: time.Now(),
	}

	mockSessionsRepo.On("GetBySessionToken", mock.Anything, sessionToken).Return(mockSession, nil)
	mockVisitorsRepo.On("GetByID", mock.Anything, int64(10)).Return(mockVisitor, nil)
	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, fingerprint).Return(mockDevice, nil)
	mockVisitorsRepo.On("CreatePageView", mock.Anything, mock.Anything).Return(mockPV, nil)

	req := visitorsv1.CreatePageViewRequest{
		SessionToken:    sessionToken,
		FingerprintHash: fingerprint,
		URL:             "https://example.com/home",
		Title:           "Home Page",
		Referrer:        "https://google.com",
	}

	res, err := service.CreatePageView(context.Background(), req, &ent.Site{})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "pv-uuid", res.UUID)
	assert.Equal(t, "session-uuid", res.SessionUUID)
	assert.Equal(t, "visitor-uuid", res.VisitorUUID)
	assert.Equal(t, "device-uuid", res.DeviceUUID)
	assert.Equal(t, "https://example.com/home", res.URL)
}

func TestCreateEvent_NewSession(t *testing.T) {
	mockSitesRepo := new(sitesMocks.MockRepository)
	sitesService := sitesv1.NewService(mockSitesRepo, nil)

	mockVisitorsRepo := new(visitorsMocks.MockRepository)
	mockDevicesRepo := new(devicesMocks.MockRepository)
	mockSessionsRepo := new(visitor_sessionsMocks.MockRepository)
	mockConsentsRepo := new(consentsMocks.MockRepository)
	sessionSvc := visitor_sessionsv1.NewService(mockSessionsRepo, nil, config.AppConfig{})
	consentsSvc := consentsv1.NewService(mockConsentsRepo, sessionSvc, nil)
	service := visitorsv1.NewService(mockVisitorsRepo, sitesService, mockDevicesRepo, sessionSvc, consentsSvc, nil)

	sessionToken := "new-sess-123"
	fingerprint := "new-fp-123"

	mockVisitor := &ent.Visitor{
		ID:   10,
		UUID: "new-visitor-uuid",
	}

	mockDevice := &ent.Device{
		ID:              20,
		UUID:            "new-device-uuid",
		FingerprintHash: fingerprint,
	}

	mockSession := &ent.VisitorSession{
		ID:           30,
		UUID:         "new-session-uuid",
		VisitorID:    10,
		DeviceID:     20,
		SessionToken: sessionToken,
	}

	mockEvent := &ent.VisitorEvent{
		UUID:            "new-event-uuid",
		EventName:       "page_load",
		EventProperties: nil,
		CreatedAt:       time.Now(),
	}

	// 1. getSession -> not found
	mockSessionsRepo.On("GetBySessionToken", mock.Anything, sessionToken).Return(nil, &ent.NotFoundError{}).Once()
	// 2. getDevice -> not found
	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, fingerprint).Return(nil, &ent.NotFoundError{}).Once()
	// 3. createVisitor
	mockVisitorsRepo.On("GetByAnonymousID", mock.Anything, mock.AnythingOfType("string")).Return(nil, &ent.NotFoundError{}).Once()
	mockVisitorsRepo.On("Create", mock.Anything, mock.Anything).Return(mockVisitor, nil).Once()
	// 4. createDevice
	mockDevicesRepo.On("Create", mock.Anything, mock.Anything).Return(mockDevice, nil).Once()
	// 5. createSession
	mockSessionsRepo.On("GetBySessionToken", mock.Anything, sessionToken).Return(nil, &ent.NotFoundError{}).Once()
	mockSessionsRepo.On("Create", mock.Anything, mock.Anything).Return(mockSession, nil).Once()
	mockVisitorsRepo.On("CreateSiteVisitor", mock.Anything, mock.Anything).Return(&ent.SiteVisitor{}, nil).Once()
	mockConsentsRepo.On("SetVisitorConsent", mock.Anything, mock.Anything, mock.Anything, "none", "pending").Return(&ent.VisitorConsent{}, nil).Once()

	// Later when generating response
	mockVisitorsRepo.On("GetByID", mock.Anything, int64(10)).Return(mockVisitor, nil).Once()
	mockDevicesRepo.On("GetByFingerprintHash", mock.Anything, fingerprint).Return(mockDevice, nil).Once()
	mockVisitorsRepo.On("CreateEvent", mock.Anything, mock.Anything).Return(mockEvent, nil).Once()

	req := visitorsv1.CreateEventRequest{
		SessionToken:    sessionToken,
		FingerprintHash: fingerprint,
		EventName:       "page_load",
	}

	res, err := service.CreateEvent(context.Background(), req, &ent.Site{})

	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, "new-event-uuid", res.UUID)
	assert.Equal(t, "new-session-uuid", res.SessionUUID)
	assert.Equal(t, "new-visitor-uuid", res.VisitorUUID)
	assert.Equal(t, "new-device-uuid", res.DeviceUUID)
}
