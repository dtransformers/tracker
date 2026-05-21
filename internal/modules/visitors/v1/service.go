package v1

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	devicesv1 "tracking/internal/modules/devices/v1"
	sitesv1 "tracking/internal/modules/sites/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
)

type Service struct {
	repo         Repository
	sitesService *sitesv1.Service
	deviceRepo   devicesv1.Repository
	sessionRepo  visitor_sessionsv1.Repository
	logger       *slog.Logger
	timeout      time.Duration
}

func NewService(repo Repository, sitesService *sitesv1.Service, deviceRepo devicesv1.Repository, sessionRepo visitor_sessionsv1.Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:         repo,
		sitesService: sitesService,
		deviceRepo:   deviceRepo,
		sessionRepo:  sessionRepo,
		logger:       logger,
		timeout:      5 * time.Second,
	}
}

type CreateSessionRequest struct {
	PublicKey string `json:"-"`

	AnonymousID string `json:"anonymous_id" swaggerignore:"true"`

	// Device Info Extracted by the tracking SDK
	FingerprintHash string `json:"fingerprint_hash" binding:"required"`
	DeviceType      string `json:"device_type" swaggerignore:"true"`
	OS              string `json:"os" swaggerignore:"true"`
	Browser         string `json:"browser" swaggerignore:"true"`
	Language        string `json:"language" swaggerignore:"true"`

	// Geo location information extracted from ip
	AddressIP string `json:"address_ip" swaggerignore:"true"`
	Country   string `json:"country" swaggerignore:"true"`
	Region    string `json:"region" swaggerignore:"true"`
	City      string `json:"city" swaggerignore:"true"`
	Timezone  string `json:"timezone" swaggerignore:"true"`

	// Session Info
	SessionToken        string `json:"session_token" swaggerignore:"true"`
	EntryPage           string `json:"entry_page"`
	ExitPage            string `json:"exit_page"`
	AcquisitionSource   string `json:"acquisition_source"`
	AcquisitionMedium   string `json:"acquisition_medium"`
	AcquisitionCampaign string `json:"acquisition_campaign"`
	ReferrerHost        string `json:"referrer_host"`
}

type VisitorResponse struct {
	UUID            string    `json:"uuid"`
	AnonymousID     string    `json:"anonymous_id"`
	PrimaryUserID   *string   `json:"primary_user_id,omitempty"`
	FirstSeenAt     time.Time `json:"first_seen_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	IsBot           bool      `json:"is_bot"`
	IsVerifiedHuman bool      `json:"is_verified_human"`
}

type CreateEventRequest struct {
	SessionToken    string                 `json:"session_token" binding:"required"`
	FingerprintHash string                 `json:"fingerprint_hash" binding:"required"`
	EventName       string                 `json:"event_name" binding:"required"`
	EventProperties map[string]interface{} `json:"event_properties"`

	// Geolocation fields if provided
	Country  string `json:"country" swaggerignore:"true"`
	Region   string `json:"region" swaggerignore:"true"`
	City     string `json:"city" swaggerignore:"true"`
	Timezone string `json:"timezone" swaggerignore:"true"`
}

type EventResponse struct {
	UUID            string                 `json:"uuid"`
	SessionUUID     string                 `json:"session_uuid"`
	VisitorUUID     string                 `json:"visitor_uuid"`
	DeviceUUID      string                 `json:"device_uuid"`
	EventName       string                 `json:"event_name"`
	EventProperties map[string]interface{} `json:"event_properties,omitempty"`
	CreatedAt       time.Time              `json:"created_at"`
}

type CreatePageViewRequest struct {
	SessionToken    string `json:"session_token" binding:"required"`
	FingerprintHash string `json:"fingerprint_hash" binding:"required"`
	URL             string `json:"url" binding:"required"`
	Title           string `json:"title"`
	Referrer        string `json:"referrer"`

	// Geolocation fields if provided
	Country  string `json:"country" swaggerignore:"true"`
	Region   string `json:"region" swaggerignore:"true"`
	City     string `json:"city" swaggerignore:"true"`
	Timezone string `json:"timezone" swaggerignore:"true"`
}

type PageViewResponse struct {
	UUID        string    `json:"uuid"`
	SessionUUID string    `json:"session_uuid"`
	VisitorUUID string    `json:"visitor_uuid"`
	DeviceUUID  string    `json:"device_uuid"`
	URL         string    `json:"url"`
	Title       string    `json:"title,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

/*
CreateSession godoc
responsible on creating the visitor session  since the sdk didn't find any current session, the session creation implecitly requires a device if not found, means also the visitor entity is not found since the device and the entity are conencted, also we need to store this event as a visit to the website on the site visitor entity
input: the device info, the request info, the site info : all those info are on the request context
response: the visitor session
*/
func (s *Service) CreateSession(ctx context.Context, req CreateSessionRequest, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var originalIP string
	var isBot bool
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		originalIP = tracking.OriginalIP
		if tracking.Cookies != nil && tracking.Cookies["anonymous_id"] != "" {
			req.AnonymousID = tracking.Cookies["anonymous_id"]
		}
		req.DeviceType = tracking.DeviceType
		if req.OS == "" {
			req.OS = tracking.OS
		}
		if req.Browser == "" {
			req.Browser = tracking.Browser
		}
		if req.DeviceType == "" {
			req.DeviceType = tracking.DeviceType
		}
		if req.Language == "" {
			req.Language = tracking.AcceptLanguage
		}
		req.AddressIP = tracking.OriginalIP

		if req.ReferrerHost == "" {
			if u, err := url.Parse(tracking.Referer); err == nil && u.Host != "" {
				req.ReferrerHost = u.Host
			} else {
				req.ReferrerHost = tracking.Referer
			}
		}
		isBot = tracking.Bot

		if req.Country == "" {
			req.Country = tracking.Country
		}
		if req.Region == "" {
			req.Region = tracking.Region
		}
		if req.City == "" {
			req.City = tracking.City
		}
		if req.Timezone == "" {
			req.Timezone = tracking.Timezone
		}
	}

	s.logger.InfoContext(ctx, "tracking visitor from ip", "address_ip", originalIP)

	device, err := s.getDevice(ctx, req.FingerprintHash)
	if err != nil {
		return nil, err
	}

	var v *ent.Visitor
	var isNewDevice bool

	if device == nil {
		isNewDevice = true
		// With the creation of every device, a new visitor must be created.
		// If req.AnonymousID already exists, generate a new one to avoid unique constraint violations.
		if req.AnonymousID != "" {
			existingV, err := s.getVisitor(ctx, req.AnonymousID)
			if err != nil {
				return nil, err
			}
			if existingV != nil {
				req.AnonymousID = uuid.New().String()
			}
		}
		v, err = s.createVisitor(ctx, req, isBot)
		if err != nil {
			return nil, err
		}
		device, err = s.createDevice(ctx, req, int64(v.ID))
		if err != nil {
			return nil, err
		}
	} else {
		device, err = s.updateDevice(ctx, device, req)
		if err != nil {
			return nil, err
		}
		v, err = s.repo.GetByID(ctx, device.VisitorID)
		if err != nil {
			return nil, err
		}
		v, err = s.updateVisitor(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	var sessionID int64
	if isNewDevice || req.SessionToken != "" {
		if req.SessionToken == "" {
			req.SessionToken = uuid.New().String()
		}

		session, err := s.getSession(ctx, req.SessionToken)
		if err != nil {
			return nil, err
		}

		if session == nil {
			sess, err := s.createSession(ctx, req, int64(v.ID), int64(device.ID))
			if err != nil {
				return nil, err
			}
			sessionID = int64(sess.ID)
		} else {
			sess, err := s.updateSession(ctx, session, req)
			if err != nil {
				return nil, err
			}
			sessionID = int64(sess.ID)
		}
	}

	if sessionID != 0 {
		_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
			UUID:      uuid.New().String(),
			VisitorID: int64(v.ID),
			SiteID:    int64(site.ID),
			SessionID: sessionID,
		})
		if err != nil {
			s.logger.WarnContext(ctx, "failed to create site visitor link", "error", err)
		}
	}

	return toVisitorResponse(v), nil
}

func (s *Service) getVisitor(ctx context.Context, anonymousID string) (*ent.Visitor, error) {
	if anonymousID == "" {
		return nil, nil
	}
	v, err := s.repo.GetByAnonymousID(ctx, anonymousID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get visitor: %w", err)
	}
	return v, nil
}

func (s *Service) createVisitor(ctx context.Context, req CreateSessionRequest, isBot bool) (*ent.Visitor, error) {
	anonID := req.AnonymousID
	if anonID == "" {
		anonID = uuid.New().String()
	}
	v, err := s.repo.Create(ctx, CreateSessionInput{
		UUID:            uuid.New().String(),
		AnonymousID:     anonID,
		IsBot:           isBot,
		IsVerifiedHuman: !isBot,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create visitor: %w", err)
	}
	return v, nil
}

func (s *Service) updateVisitor(ctx context.Context, v *ent.Visitor) (*ent.Visitor, error) {
	now := time.Now()
	updated, err := s.repo.Update(ctx, v.UUID, UpdateVisitorInput{
		LastSeenAt: &now,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update visitor: %w", err)
	}
	return updated, nil
}

func (s *Service) getDevice(ctx context.Context, fingerprintHash string) (*ent.Device, error) {
	if fingerprintHash == "" {
		return nil, nil
	}
	device, err := s.deviceRepo.GetByFingerprintHash(ctx, fingerprintHash)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return device, nil
}

func (s *Service) createDevice(ctx context.Context, req CreateSessionRequest, visitorID int64) (*ent.Device, error) {
	d, err := s.deviceRepo.Create(ctx, devicesv1.CreateDeviceInput{
		UUID:            uuid.New().String(),
		VisitorID:       visitorID,
		FingerprintHash: req.FingerprintHash,
		DeviceType:      req.DeviceType,
		OS:              req.OS,
		AddressIP:       req.AddressIP,
		Country:         req.Country,
		Region:          req.Region,
		City:            req.City,
		Browser:         req.Browser,
		Language:        req.Language,
		Timezone:        req.Timezone,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}
	return d, nil
}

func (s *Service) updateDevice(ctx context.Context, device *ent.Device, req CreateSessionRequest) (*ent.Device, error) {
	now := time.Now()
	d, err := s.deviceRepo.Update(ctx, device.UUID, devicesv1.UpdateDeviceInput{
		LastSeenAt: &now,
		AddressIP:  &req.AddressIP,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}
	return d, nil
}

func (s *Service) getSession(ctx context.Context, sessionToken string) (*ent.VisitorSession, error) {
	if sessionToken == "" {
		return nil, nil
	}
	session, err := s.sessionRepo.GetBySessionToken(ctx, sessionToken)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}
	return session, nil
}

func (s *Service) createSession(ctx context.Context, req CreateSessionRequest, visitorID int64, deviceID int64) (*ent.VisitorSession, error) {
	sess, err := s.sessionRepo.Create(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           visitorID,
		DeviceID:            deviceID,
		SessionToken:        req.SessionToken,
		StartedAt:           time.Now(),
		EntryPage:           req.EntryPage,
		AcquisitionSource:   req.AcquisitionSource,
		AcquisitionMedium:   req.AcquisitionMedium,
		AcquisitionCampaign: req.AcquisitionCampaign,
		ReferrerHost:        req.ReferrerHost,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %w", err)
	}
	return sess, nil
}

func (s *Service) updateSession(ctx context.Context, session *ent.VisitorSession, req CreateSessionRequest) (*ent.VisitorSession, error) {
	now := time.Now()
	sess, err := s.sessionRepo.Update(ctx, session.UUID, visitor_sessionsv1.UpdateSessionInput{
		LastActivityAt: &now,
		ExitPage:       &req.ExitPage,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update session: %w", err)
	}
	return sess, nil
}

func toVisitorResponse(v *ent.Visitor) *VisitorResponse {
	if v == nil {
		return nil
	}
	return &VisitorResponse{
		UUID:            v.UUID,
		AnonymousID:     v.AnonymousID,
		PrimaryUserID:   v.PrimaryUserID,
		FirstSeenAt:     v.FirstSeenAt,
		LastSeenAt:      v.LastSeenAt,
		IsBot:           v.IsBot,
		IsVerifiedHuman: v.IsVerifiedHuman,
	}
}

func (s *Service) resolveSession(ctx context.Context, sessionToken, fingerprintHash string, geo CreateSessionRequest) (*ent.VisitorSession, error) {
	session, err := s.getSession(ctx, sessionToken)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	device, err := s.getDevice(ctx, fingerprintHash)
	if err != nil {
		return nil, err
	}

	var v *ent.Visitor
	if device == nil {
		var anonID string
		if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
			if tracking.Cookies != nil {
				anonID = tracking.Cookies["anonymous_id"]
			}
		}
		if anonID == "" {
			anonID = uuid.New().String()
		}
		existingV, err := s.getVisitor(ctx, anonID)
		if err != nil {
			return nil, err
		}
		if existingV != nil {
			anonID = uuid.New().String()
		}

		geo.AnonymousID = anonID
		geo.FingerprintHash = fingerprintHash
		geo.SessionToken = sessionToken

		var isBot bool
		if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
			isBot = tracking.Bot
		}

		v, err = s.createVisitor(ctx, geo, isBot)
		if err != nil {
			return nil, err
		}

		device, err = s.createDevice(ctx, geo, int64(v.ID))
		if err != nil {
			return nil, err
		}
	} else {
		v, err = s.repo.GetByID(ctx, device.VisitorID)
		if err != nil {
			return nil, err
		}
	}

	geo.FingerprintHash = fingerprintHash
	geo.SessionToken = sessionToken
	session, err = s.createSession(ctx, geo, int64(v.ID), int64(device.ID))
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest) (*EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var geo CreateSessionRequest
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		geo.AddressIP = tracking.OriginalIP
		geo.DeviceType = tracking.DeviceType
		geo.OS = tracking.OS
		geo.Browser = tracking.Browser
		geo.Language = tracking.AcceptLanguage
		geo.Country = tracking.Country
		geo.Region = tracking.Region
		geo.City = tracking.City
		geo.Timezone = tracking.Timezone
	}
	if req.Country != "" {
		geo.Country = req.Country
	}
	if req.Region != "" {
		geo.Region = req.Region
	}
	if req.City != "" {
		geo.City = req.City
	}
	if req.Timezone != "" {
		geo.Timezone = req.Timezone
	}

	session, err := s.resolveSession(ctx, req.SessionToken, req.FingerprintHash, geo)
	if err != nil {
		return nil, err
	}

	v, err := s.repo.GetByID(ctx, session.VisitorID)
	if err != nil {
		return nil, err
	}

	device, err := s.getDevice(ctx, req.FingerprintHash)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device not found")
	}

	ev, err := s.repo.CreateEvent(ctx, CreateEventInput{
		SessionID:       int64(session.ID),
		VisitorID:       int64(v.ID),
		DeviceID:        int64(device.ID),
		EventName:       req.EventName,
		EventProperties: req.EventProperties,
	})
	if err != nil {
		return nil, err
	}

	return &EventResponse{
		UUID:            ev.UUID,
		SessionUUID:     session.UUID,
		VisitorUUID:     v.UUID,
		DeviceUUID:      device.UUID,
		EventName:       ev.EventName,
		EventProperties: ev.EventProperties,
		CreatedAt:       ev.CreatedAt,
	}, nil
}

func (s *Service) CreatePageView(ctx context.Context, req CreatePageViewRequest) (*PageViewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var geo CreateSessionRequest
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		geo.AddressIP = tracking.OriginalIP
		geo.DeviceType = tracking.DeviceType
		geo.OS = tracking.OS
		geo.Browser = tracking.Browser
		geo.Language = tracking.AcceptLanguage
		geo.Country = tracking.Country
		geo.Region = tracking.Region
		geo.City = tracking.City
		geo.Timezone = tracking.Timezone
	}
	if req.Country != "" {
		geo.Country = req.Country
	}
	if req.Region != "" {
		geo.Region = req.Region
	}
	if req.City != "" {
		geo.City = req.City
	}
	if req.Timezone != "" {
		geo.Timezone = req.Timezone
	}

	session, err := s.resolveSession(ctx, req.SessionToken, req.FingerprintHash, geo)
	if err != nil {
		return nil, err
	}

	v, err := s.repo.GetByID(ctx, session.VisitorID)
	if err != nil {
		return nil, err
	}

	device, err := s.getDevice(ctx, req.FingerprintHash)
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, fmt.Errorf("device not found")
	}

	pv, err := s.repo.CreatePageView(ctx, CreatePageViewInput{
		SessionID: int64(session.ID),
		VisitorID: int64(v.ID),
		DeviceID:  int64(device.ID),
		URL:       req.URL,
		Title:     req.Title,
		Referrer:  req.Referrer,
	})
	if err != nil {
		return nil, err
	}

	return &PageViewResponse{
		UUID:        pv.UUID,
		SessionUUID: session.UUID,
		VisitorUUID: v.UUID,
		DeviceUUID:  device.UUID,
		URL:         pv.URL,
		Title:       pv.Title,
		Referrer:    pv.Referrer,
		CreatedAt:   pv.CreatedAt,
	}, nil
}
