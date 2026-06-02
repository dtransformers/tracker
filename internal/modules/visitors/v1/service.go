package v1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"time"

	"github.com/google/uuid"

	"tracking/internal/ent"
	"tracking/internal/middlewares"
	consentsv1 "tracking/internal/modules/consents/v1"
	devicesv1 "tracking/internal/modules/devices/v1"
	sitesv1 "tracking/internal/modules/sites/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
)

type Service struct {
	repo         Repository
	sitesService *sitesv1.Service
	deviceRepo   devicesv1.Repository
	sessionSvc   *visitor_sessionsv1.Service
	consentsSvc  *consentsv1.Service
	logger       *slog.Logger
	timeout      time.Duration
}

func NewService(repo Repository, sitesService *sitesv1.Service, deviceRepo devicesv1.Repository, sessionSvc *visitor_sessionsv1.Service, consentsSvc *consentsv1.Service, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:         repo,
		sitesService: sitesService,
		deviceRepo:   deviceRepo,
		sessionSvc:   sessionSvc,
		consentsSvc:  consentsSvc,
		logger:       logger,
		timeout:      5 * time.Second,
	}
}

type CreateVisitorRequest struct {
	PublicKey string `json:"-"`

	VisitorUUID string `json:"visitor_uuid" swaggerignore:"true"`
	SessionUUID string `json:"session_uuid" swaggerignore:"true"`

	FingerprintHash string `json:"fingerprint_hash" binding:"required"`
	DeviceType      string `json:"device_type" swaggerignore:"true"`
	OS              string `json:"os" swaggerignore:"true"`
	Browser         string `json:"browser" swaggerignore:"true"`
	BrowserVersion  string `json:"browser_version" swaggerignore:"true"`
	Language        string `json:"language" swaggerignore:"true"`
	ScreenWidth     int    `json:"screen_width" swaggerignore:"true"`
	ScreenHeight    int    `json:"screen_height" swaggerignore:"true"`
	// Geo location information extracted from ip
	AddressIP string `json:"address_ip" swaggerignore:"true"`
	Country   string `json:"country" swaggerignore:"true"`
	Region    string `json:"region" swaggerignore:"true"`
	City      string `json:"city" swaggerignore:"true"`
	Timezone  string `json:"timezone" swaggerignore:"true"`
	ISP       string `json:"isp" swaggerignore:"true"`

	// Session Info
	SessionToken        string `json:"session_token" swaggerignore:"true"`
	EntryPage           string `json:"entry_page"`
	ExitPage            string `json:"exit_page"`
	AcquisitionSource   string `json:"acquisition_source"`
	AcquisitionMedium   string `json:"acquisition_medium"`
	AcquisitionCampaign string `json:"acquisition_campaign"`
	CampaignTerm        string `json:"campaign_term"`
	CampaignContent     string `json:"campaign_content"`
	CampaignID          string `json:"campaign_id"`
	ReferrerHost        string `json:"referrer_host"`
}

type VisitorResponse struct {
	UUID            string    `json:"uuid"`
	SessionID       string    `json:"session_id"`
	PrimaryUserID   *string   `json:"primary_user_id,omitempty"`
	FirstSeenAt     time.Time `json:"first_seen_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
	IsBot           bool      `json:"is_bot"`
	IsVerifiedHuman bool      `json:"is_verified_human"`
}

type FullVisitorSessionDetails struct {
	Visitor *ent.Visitor        `json:"visitor"`
	Device  *ent.Device         `json:"device"`
	Session *ent.VisitorSession `json:"session"`
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

func (s *Service) GetSessionByUUID(ctx context.Context, VisitorUUID string, SessionToken string) (*FullVisitorSessionDetails, error) {
	session, err := s.sessionSvc.GetSessionByTokenRaw(ctx, SessionToken)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found,Please Clear Cookies")
	}
	device, err := s.deviceRepo.GetByID(ctx, int(session.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found,Please Clear Cookies")
	}
	v, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("visitor not found, Please Clear Cookies")
	}

	return toFullSessionDetails(v, device, session), nil
}

func (s *Service) GetExistingSession(ctx context.Context, VisitorUUID string, SessionUUID string) (*VisitorResponse, error) {
	session, err := s.sessionSvc.GetBySessionTokenRaw(ctx, SessionUUID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("session not found,Please Clear Cookies")
	}
	device, err := s.getDeviceByID(ctx, int(session.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found,Please Clear Cookies")
	}
	v, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, errors.New("visitor not found, Please Clear Cookies")
	}
	response := toVisitorResponse(v)
	response.SessionID = session.SessionToken
	return response, nil

}

func (s *Service) createNewSessionForVisitor(ctx context.Context, VisitorUUID string, entryPage string, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	visitor, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if visitor == nil {
		return nil, errors.New("visitor not found")
	}

	visitor, err = s.updateVisitor(ctx, visitor)
	if err != nil {
		return nil, err
	}

	visitorLastActiveSession, err := s.sessionSvc.GetVisitorLastActiveSession(ctx, int64(visitor.ID))
	if err != nil {
		return nil, err
	}
	if visitorLastActiveSession == nil {
		visitorLastSession, err := s.sessionSvc.GetVisitorLastSession(ctx, int64(visitor.ID))
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errors.New("visitor does not exisits, Please Clear Coockies")
			}
			return nil, err
		}
		if visitorLastSession == nil {
			return nil, errors.New("visitor does not exisits, Please Clear Coockies")
		}
		visitorLastActiveSession, err = s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
			UUID:                uuid.New().String(),
			VisitorID:           int64(visitor.ID),
			DeviceID:            int64(visitorLastSession.DeviceID),
			SessionToken:        uuid.New().String(),
			StartedAt:           time.Now().UTC(),
			EntryPage:           entryPage,
			AcquisitionSource:   visitorLastSession.AcquisitionSource,
			AcquisitionMedium:   visitorLastSession.AcquisitionMedium,
			AcquisitionCampaign: visitorLastSession.AcquisitionCampaign,
			CampaignTerm:        visitorLastSession.CampaignTerm,
			CampaignContent:     visitorLastSession.CampaignContent,
			CampaignID:          visitorLastSession.CampaignID,
			ReferrerHost:        visitorLastSession.ReferrerHost,
		})
		if err != nil {
			return nil, err
		}
	}
	device, err := s.getDeviceByID(ctx, int(visitorLastActiveSession.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, errors.New("device not found")
	}

	newSession, err := s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           int64(visitor.ID),
		DeviceID:            int64(device.ID),
		SessionToken:        uuid.New().String(),
		StartedAt:           time.Now().UTC(),
		EntryPage:           entryPage,
		AcquisitionSource:   visitorLastActiveSession.AcquisitionSource,
		AcquisitionMedium:   visitorLastActiveSession.AcquisitionMedium,
		AcquisitionCampaign: visitorLastActiveSession.AcquisitionCampaign,
		CampaignTerm:        visitorLastActiveSession.CampaignTerm,
		CampaignContent:     visitorLastActiveSession.CampaignContent,
		CampaignID:          visitorLastActiveSession.CampaignID,
		ReferrerHost:        visitorLastActiveSession.ReferrerHost,
	})
	if err != nil {
		return nil, err
	}
	if newSession == nil {
		return nil, errors.New("failed to create new session")
	}

	_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
		UUID:      uuid.New().String(),
		VisitorID: int64(visitor.ID),
		SiteID:    int64(site.ID),
		SessionID: int64(newSession.ID),
	})
	if err != nil {
		s.logger.WarnContext(ctx, "failed to create site visitor link", "error", err)
	} else {
		err = s.consentsSvc.SetDefaultConsent(ctx, int64(visitor.ID), int64(site.ID))
		if err != nil {
			s.logger.WarnContext(ctx, "failed to set default visitor consent", "error", err)
		}
	}

	response := toVisitorResponse(visitor)
	response.SessionID = newSession.SessionToken
	return response, nil
}

func (s *Service) enrichRequestWithTracking(req *CreateVisitorRequest, tracking *middlewares.TrackingInfo) {
	if tracking == nil {
		return
	}

	if tracking.Cookies != nil {
		if tracking.Cookies["visitor_uuid"] != "" {
			req.VisitorUUID = tracking.Cookies["visitor_uuid"]
		}
		if tracking.Cookies["session_id"] != "" {
			req.SessionUUID = tracking.Cookies["session_id"]
		}
	}

	if req.DeviceType == "" {
		req.DeviceType = tracking.DeviceType
	}
	if req.OS == "" {
		req.OS = tracking.OS
	}
	if req.Browser == "" {
		req.Browser = tracking.Browser
	}
	if req.BrowserVersion == "" {
		req.BrowserVersion = tracking.BrowserVersion
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

func (s *Service) resolveVisitorAndDevice(ctx context.Context, req CreateVisitorRequest, isBot bool) (*ent.Visitor, *ent.Device, bool, error) {
	device, err := s.getDevice(ctx, req.FingerprintHash)
	if err != nil {
		return nil, nil, false, err
	}

	var v *ent.Visitor
	var isNewDevice bool

	if device == nil {
		isNewDevice = true
		v, err = s.createVisitor(ctx, req, isBot)
		if err != nil {
			return nil, nil, false, err
		}
		device, err = s.createDevice(ctx, req, int64(v.ID))
		if err != nil {
			return nil, nil, false, err
		}
	} else {
		device, err = s.updateDevice(ctx, device, req)
		if err != nil {
			return nil, nil, false, err
		}
		v, err = s.repo.GetByID(ctx, device.VisitorID)
		if err != nil {
			return nil, nil, false, err
		}
		v, err = s.updateVisitor(ctx, v)
		if err != nil {
			return nil, nil, false, err
		}
	}
	return v, device, isNewDevice, nil
}

func (s *Service) resolveSessionForCreate(ctx context.Context, req CreateVisitorRequest, v *ent.Visitor, device *ent.Device, isNewDevice bool) (int64, string, bool, error) {
	var sessionID int64
	var sessionToken string
	var isNewSession bool

	if isNewDevice || req.SessionToken != "" {
		if req.SessionToken == "" {
			req.SessionToken = uuid.New().String()
		}

		session, err := s.sessionSvc.GetBySessionTokenRaw(ctx, req.SessionToken)
		if err != nil {
			return 0, "", false, err
		}

		if session == nil {
			sess, err := s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
				UUID:                uuid.New().String(),
				VisitorID:           int64(v.ID),
				DeviceID:            int64(device.ID),
				SessionToken:        req.SessionToken,
				StartedAt:           time.Now().UTC(),
				EntryPage:           req.EntryPage,
				AcquisitionSource:   req.AcquisitionSource,
				AcquisitionMedium:   req.AcquisitionMedium,
				AcquisitionCampaign: req.AcquisitionCampaign,
				CampaignTerm:        req.CampaignTerm,
				CampaignContent:     req.CampaignContent,
				CampaignID:          req.CampaignID,
				ReferrerHost:        req.ReferrerHost,
			})
			if err != nil {
				return 0, "", false, err
			}
			sessionID = int64(sess.ID)
			sessionToken = sess.SessionToken
			isNewSession = true
		} else {
			now := time.Now()
			sess, err := s.sessionSvc.UpdateSession(ctx, session.UUID, visitor_sessionsv1.UpdateSessionInput{
				LastActivityAt: &now,
				ExitPage:       &req.ExitPage,
			})
			if err != nil {
				return 0, "", false, err
			}
			sessionID = int64(sess.ID)
			sessionToken = sess.SessionToken
		}
	}
	return sessionID, sessionToken, isNewSession, nil
}
func (s *Service) CreateSession(ctx context.Context, req CreateVisitorRequest, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var isBot bool
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		s.enrichRequestWithTracking(&req, tracking)
		isBot = tracking.Bot
		if req.VisitorUUID != "" && req.SessionUUID != "" {
			return s.GetExistingSession(ctx, req.VisitorUUID, req.SessionUUID)
		}
		if req.VisitorUUID != "" && req.SessionUUID == "" {
			return s.createNewSessionForVisitor(ctx, req.VisitorUUID, req.EntryPage, site)
		}
	}

	s.logger.InfoContext(ctx, "tracking visitor from ip", "address_ip", req.AddressIP)

	v, device, isNewDevice, err := s.resolveVisitorAndDevice(ctx, req, isBot)
	if err != nil {
		return nil, err
	}

	sessionID, sessionToken, isNewSession, err := s.resolveSessionForCreate(ctx, req, v, device, isNewDevice)
	if err != nil {
		return nil, err
	}

	if isNewSession && sessionID != 0 {
		_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
			UUID:      uuid.New().String(),
			VisitorID: int64(v.ID),
			SiteID:    int64(site.ID),
			SessionID: sessionID,
		})
		if err != nil {
			s.logger.WarnContext(ctx, "failed to create site visitor link", "error", err)
		} else {
			err = s.consentsSvc.SetDefaultConsent(ctx, int64(v.ID), int64(site.ID))
			if err != nil {
				s.logger.WarnContext(ctx, "failed to set default visitor consent", "error", err)
			}
		}
	}

	res := toVisitorResponse(v)
	if sessionToken != "" {
		res.SessionID = sessionToken
	}

	return res, nil
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

func (s *Service) getVistorByUUID(ctx context.Context, visitorUUID string) (*ent.Visitor, error) {

	v, err := s.repo.GetByUUID(ctx, visitorUUID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get visitor: %w", err)
	}
	return v, nil
}

func (s *Service) createVisitor(ctx context.Context, req CreateVisitorRequest, isBot bool) (*ent.Visitor, error) {

	v, err := s.repo.Create(ctx, CreateSessionInput{
		UUID:            uuid.New().String(),
		AnonymousID:     uuid.New().String(),
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

func (s *Service) getDeviceByID(ctx context.Context, id int) (*ent.Device, error) {
	if id == 0 {
		return nil, nil
	}
	device, err := s.deviceRepo.GetByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return device, nil
}

func (s *Service) createDevice(ctx context.Context, req CreateVisitorRequest, visitorID int64) (*ent.Device, error) {
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
		BrowserVersion:  req.BrowserVersion,
		Language:        req.Language,
		Timezone:        req.Timezone,
		ISP:             req.ISP,
		ScreenWidth:     req.ScreenWidth,
		ScreenHeight:    req.ScreenHeight,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}
	return d, nil
}

func (s *Service) updateDevice(ctx context.Context, device *ent.Device, req CreateVisitorRequest) (*ent.Device, error) {
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

func toVisitorResponse(v *ent.Visitor) *VisitorResponse {
	if v == nil {
		return nil
	}
	return &VisitorResponse{
		UUID:            v.UUID,
		PrimaryUserID:   v.PrimaryUserID,
		FirstSeenAt:     v.FirstSeenAt,
		LastSeenAt:      v.LastSeenAt,
		IsBot:           v.IsBot,
		IsVerifiedHuman: v.IsVerifiedHuman,
	}
}

func toFullSessionDetails(v *ent.Visitor, d *ent.Device, s *ent.VisitorSession) *FullVisitorSessionDetails {
	if v == nil || d == nil || s == nil {
		return nil
	}
	return &FullVisitorSessionDetails{Visitor: v, Device: d, Session: s}
}

func (s *Service) ResolveSession(ctx context.Context, sessionToken, fingerprintHash string, geo CreateVisitorRequest, site *ent.Site) (*ent.VisitorSession, error) {
	session, err := s.sessionSvc.GetBySessionTokenRaw(ctx, sessionToken)
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
				anonID = tracking.Cookies["visitor_uuid"]
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

		v, err = s.updateVisitor(ctx, v)
		if err != nil {
			return nil, err
		}
	}

	geo.FingerprintHash = fingerprintHash
	geo.SessionToken = sessionToken
	session, err = s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           int64(v.ID),
		DeviceID:            int64(device.ID),
		SessionToken:        geo.SessionToken,
		StartedAt:           time.Now().UTC(),
		EntryPage:           geo.EntryPage,
		AcquisitionSource:   geo.AcquisitionSource,
		AcquisitionMedium:   geo.AcquisitionMedium,
		AcquisitionCampaign: geo.AcquisitionCampaign,
		CampaignTerm:        geo.CampaignTerm,
		CampaignContent:     geo.CampaignContent,
		CampaignID:          geo.CampaignID,
		ReferrerHost:        geo.ReferrerHost,
	})
	if err != nil {
		return nil, err
	}

	_, err = s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
		UUID:      uuid.New().String(),
		VisitorID: int64(v.ID),
		SiteID:    int64(site.ID),
		SessionID: int64(session.ID),
	})
	if err != nil {
		s.logger.WarnContext(ctx, "failed to create site visitor link", "error", err)
	} else {
		err = s.consentsSvc.SetDefaultConsent(ctx, int64(v.ID), int64(site.ID))
		if err != nil {
			s.logger.WarnContext(ctx, "failed to set default visitor consent", "error", err)
		}
	}

	return session, nil
}

func (s *Service) CreateEvent(ctx context.Context, req CreateEventRequest, site *ent.Site) (*EventResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	geo := CreateVisitorRequest{
		Country:  req.Country,
		Region:   req.Region,
		City:     req.City,
		Timezone: req.Timezone,
	}
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		s.enrichRequestWithTracking(&geo, tracking)
	}

	session, err := s.ResolveSession(ctx, req.SessionToken, req.FingerprintHash, geo, site)
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

func (s *Service) CreatePageView(ctx context.Context, req CreatePageViewRequest, site *ent.Site) (*PageViewResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	geo := CreateVisitorRequest{
		Country:  req.Country,
		Region:   req.Region,
		City:     req.City,
		Timezone: req.Timezone,
	}
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		s.enrichRequestWithTracking(&geo, tracking)
	}

	session, err := s.ResolveSession(ctx, req.SessionToken, req.FingerprintHash, geo, site)
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
