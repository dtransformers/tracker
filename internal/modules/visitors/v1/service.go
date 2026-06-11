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
	consentsv1 "tracking/internal/modules/consents/v1"
	devicesv1 "tracking/internal/modules/devices/v1"
	sitesv1 "tracking/internal/modules/sites/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
	"tracking/pkg/constants"
	custom_errors "tracking/pkg/errors"
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
		timeout:      constants.DefaultServiceTimeout,
	}
}

func (s *Service) GetSessionByUUID(ctx context.Context, VisitorUUID string, SessionUUID string) (*FullVisitorSessionDetails, error) {
	session, err := s.sessionSvc.GetBySessionUUIDRaw(ctx, SessionUUID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "session not found, please clear cookies")
	}
	device, err := s.deviceRepo.GetByID(ctx, int(session.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "device not found, please clear cookies")
	}
	v, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "visitor not found, please clear cookies")
	}
	if int(session.VisitorID) != v.ID {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "session does not belong to this visitor")
	}

	return toFullSessionDetails(v, device, session), nil
}

func (s *Service) VerifyVisitor(ctx context.Context, visitorUUID string, sessionUUID string) (*VerifyVisitorData, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	session, err := s.sessionSvc.GetBySessionUUIDRaw(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return &VerifyVisitorData{IsValid: false}, nil
	}

	v, err := s.getVistorByUUID(ctx, visitorUUID)
	if err != nil {
		return nil, err
	}
	if v == nil || int(session.VisitorID) != v.ID {
		return &VerifyVisitorData{IsValid: false}, nil
	}

	device, err := s.deviceRepo.GetByID(ctx, int(session.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return &VerifyVisitorData{IsValid: false}, nil
	}

	return &VerifyVisitorData{
		IsValid: true,
		Session: toVisitorSessionResponse(session),
		Device:  toVisitorDeviceResponse(device),
	}, nil
}

func (s *Service) GetExistingSession(ctx context.Context, VisitorUUID string, SessionUUID string) (*VisitorResponse, error) {
	session, err := s.sessionSvc.GetBySessionUUIDRaw(ctx, SessionUUID)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "session not found, please clear cookies")
	}
	device, err := s.getDeviceByID(ctx, int(session.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "device not found, please clear cookies")
	}
	v, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if v == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "visitor not found, please clear cookies")
	}
	response := toVisitorResponse(v)
	response.SessionToken = session.SessionToken
	return response, nil

}

func (s *Service) ensureVisitorLastActiveSession(ctx context.Context, visitorID int64, entryPage string) (*ent.VisitorSession, error) {
	visitorLastActiveSession, err := s.sessionSvc.GetVisitorLastActiveSession(ctx, visitorID)
	if err != nil {
		return nil, err
	}
	if visitorLastActiveSession != nil {
		return visitorLastActiveSession, nil
	}

	visitorLastSession, err := s.sessionSvc.GetVisitorLastSession(ctx, visitorID)
	visitorError := "visitor does not exist, please clear cookies"
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, custom_errors.NewNotFoundError(constants.Version, visitorError)
		}
		return nil, err
	}
	if visitorLastSession == nil {
		return nil, custom_errors.NewNotFoundError(constants.Version, visitorError)
	}

	return s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           visitorID,
		DeviceID:            int64(visitorLastSession.DeviceID),
		SessionUUID:         uuid.New().String(),
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
}

func (s *Service) createNewSessionForVisitor(ctx context.Context, VisitorUUID string, entryPage string, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	visitor, err := s.getVistorByUUID(ctx, VisitorUUID)
	if err != nil {
		return nil, err
	}
	if visitor == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "visitor not found")
	}

	visitor, err = s.updateVisitor(ctx, visitor)
	if err != nil {
		return nil, err
	}

	visitorLastActiveSession, err := s.ensureVisitorLastActiveSession(ctx, int64(visitor.ID), entryPage)
	if err != nil {
		return nil, err
	}
	device, err := s.getDeviceByID(ctx, int(visitorLastActiveSession.DeviceID))
	if err != nil {
		return nil, err
	}
	if device == nil {
		return nil, custom_errors.NewUnauthorizedError(constants.Version, "device not found")
	}

	newSession, err := s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           int64(visitor.ID),
		DeviceID:            int64(device.ID),
		SessionUUID:         uuid.New().String(),
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
		return nil, custom_errors.NewInternalError(constants.Version, "failed to create new session")
	}

	s.linkSiteVisitorAndConsent(ctx, int64(visitor.ID), int64(site.ID), int64(newSession.ID))

	response := toVisitorResponse(visitor)
	response.SessionToken = newSession.SessionToken
	return response, nil
}

func (s *Service) enrichRequestWithTracking(req *CreateVisitorRequest, tracking *middlewares.TrackingInfo) {
	if tracking == nil {
		return
	}

	if tracking.VisitorUUID != "" {
		req.VisitorUUID = tracking.VisitorUUID
	}
	if tracking.SessionUUID != "" {
		req.SessionUUID = tracking.SessionUUID
	}

	s.enrichDeviceTracking(req, tracking)

	req.AddressIP = tracking.OriginalIP

	if req.ReferrerHost == "" {
		if u, err := url.Parse(tracking.Referer); err == nil && u.Host != "" {
			req.ReferrerHost = u.Host
		} else {
			req.ReferrerHost = tracking.Referer
		}
	}

	s.enrichGeoTracking(req, tracking)
}

func (s *Service) enrichDeviceTracking(req *CreateVisitorRequest, tracking *middlewares.TrackingInfo) {
	if req.DeviceType == "" {
		if tracking.Mobile {
			req.DeviceType = "mobile"
		} else {
			req.DeviceType = "desktop"
		}
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
}

func (s *Service) enrichGeoTracking(req *CreateVisitorRequest, tracking *middlewares.TrackingInfo) {
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
		v, err = s.createVisitor(ctx, isBot)
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

func (s *Service) getOrCreateSessionForDeviceOrToken(ctx context.Context, req CreateVisitorRequest, v *ent.Visitor, device *ent.Device) (int64, string, bool, error) {
	if req.SessionUUID == "" {
		req.SessionUUID = uuid.New().String()
	}
	if req.SessionToken == "" {
		req.SessionToken = uuid.New().String()
	}

	session, err := s.sessionSvc.GetBySessionUUIDRaw(ctx, req.SessionUUID)
	if err != nil {
		return 0, "", false, err
	}
	fmt.Println(`Session Info From DB For The Visitor With ID: `, v.ID, `And UUID: `, req.SessionUUID, `For A IsNewSession:`, false, `For Visitor UUID : `, v.UUID, `And Device UUID : `, device.UUID)

	if session == nil {
		req.SessionToken = uuid.New().String()
		sess, err := s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
			UUID:                uuid.New().String(),
			VisitorID:           int64(v.ID),
			DeviceID:            int64(device.ID),
			SessionUUID:         req.SessionUUID,
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
		return int64(sess.ID), sess.UUID, true, nil
	}

	now := time.Now()
	sess, err := s.sessionSvc.UpdateSession(ctx, session.UUID, visitor_sessionsv1.UpdateSessionInput{
		LastActivityAt: &now,
		ExitPage:       &req.ExitPage,
	})
	if err != nil {
		return 0, "", false, err
	}
	return int64(sess.ID), sess.UUID, false, nil
}

func (s *Service) getOrCreateSessionForActiveVisitor(ctx context.Context, req CreateVisitorRequest, v *ent.Visitor, device *ent.Device) (int64, string, error) {
	sess, err := s.sessionSvc.GetVisitorLastActiveSession(ctx, int64(v.ID))
	if err != nil {
		return 0, "", err
	}
	if sess == nil {
		req.SessionToken = uuid.New().String()
		sess, err = s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
			UUID:                uuid.New().String(),
			VisitorID:           int64(v.ID),
			DeviceID:            int64(device.ID),
			SessionUUID:         req.SessionUUID,
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
			return 0, "", err
		}
		return int64(sess.ID), sess.UUID, nil
	}

	now := time.Now()
	sess, err = s.sessionSvc.UpdateSession(ctx, sess.UUID, visitor_sessionsv1.UpdateSessionInput{
		LastActivityAt: &now,
	})
	if err != nil {
		return 0, "", err
	}
	return int64(sess.ID), sess.UUID, nil
}

func (s *Service) resolveSessionForCreate(ctx context.Context, req CreateVisitorRequest, v *ent.Visitor, device *ent.Device, isNewDevice bool) (int64, string, bool, error) {
	if isNewDevice || req.SessionUUID != "" {
		return s.getOrCreateSessionForDeviceOrToken(ctx, req, v, device)
	}

	sessionID, sessionToken, err := s.getOrCreateSessionForActiveVisitor(ctx, req, v, device)
	if err != nil {
		return 0, "", false, err
	}
	return sessionID, sessionToken, false, nil
}
func (s *Service) CreateSession(ctx context.Context, req CreateVisitorRequest, site *ent.Site) (*VisitorResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var isBot bool
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		s.enrichRequestWithTracking(&req, tracking)
		isBot = tracking.Bot
		if tracking.VisitorUUID != "" && tracking.SessionUUID != "" {
			res, err := s.GetExistingSession(ctx, tracking.VisitorUUID, tracking.SessionUUID)
			if err == nil {
				return res, nil
			}
			return s.createNewSessionForVisitor(ctx, tracking.VisitorUUID, req.EntryPage, site)
		}
		if tracking.VisitorUUID != "" && tracking.SessionUUID == "" {
			return s.createNewSessionForVisitor(ctx, tracking.VisitorUUID, req.EntryPage, site)
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
	fmt.Println(`Generating New Session For The Visitor With ID: `, v.ID, `And Token: `, sessionToken, `For A IsNewSession:`, isNewSession, `For Visitor UUID : `, v.UUID, `And Device UUID : `, device.UUID)
	if isNewSession && sessionID != 0 {
		s.linkSiteVisitorAndConsent(ctx, int64(v.ID), int64(site.ID), sessionID)
	}

	res := toVisitorResponse(v)
	res.SessionToken = sessionToken

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

func (s *Service) createVisitor(ctx context.Context, isBot bool) (*ent.Visitor, error) {

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

func toVisitorSessionResponse(s *ent.VisitorSession) *visitor_sessionsv1.SessionResponse {
	if s == nil {
		return nil
	}
	return &visitor_sessionsv1.SessionResponse{
		UUID:                s.UUID,
		VisitorID:           s.VisitorID,
		DeviceID:            s.DeviceID,
		SessionToken:        s.SessionToken,
		StartedAt:           s.StartedAt,
		EndedAt:             s.EndedAt,
		IsActive:            s.IsActive,
		LastActivityAt:      s.LastActivityAt,
		EntryPage:           s.EntryPage,
		ExitPage:            s.ExitPage,
		AcquisitionSource:   s.AcquisitionSource,
		AcquisitionMedium:   s.AcquisitionMedium,
		AcquisitionCampaign: s.AcquisitionCampaign,
		CampaignTerm:        s.CampaignTerm,
		CampaignContent:     s.CampaignContent,
		CampaignID:          s.CampaignID,
		ReferrerHost:        s.ReferrerHost,
		CreatedAt:           s.CreatedAt,
	}
}

func toVisitorDeviceResponse(d *ent.Device) *devicesv1.DeviceResponse {
	if d == nil {
		return nil
	}
	return &devicesv1.DeviceResponse{
		UUID:            d.UUID,
		VisitorID:       d.VisitorID,
		FingerprintHash: d.FingerprintHash,
		DeviceType:      d.DeviceType,
		OS:              d.Os,
		AddressIP:       d.AddressIP,
		Country:         d.Country,
		Region:          d.Region,
		City:            d.City,
		Browser:         d.Browser,
		BrowserVersion:  d.BrowserVersion,
		Language:        d.Language,
		Timezone:        d.Timezone,
		FirstSeenAt:     d.FirstSeenAt,
		LastSeenAt:      d.LastSeenAt,
	}
}

const (
	errCreateSiteVisitorLink = "failed to create site visitor link"
	errSetDefaultConsent     = "failed to set default visitor consent"
)

func (s *Service) linkSiteVisitorAndConsent(ctx context.Context, visitorID, siteID, sessionID int64) {
	_, err := s.repo.CreateSiteVisitor(ctx, CreateSiteVisitorInput{
		UUID:      uuid.New().String(),
		VisitorID: visitorID,
		SiteID:    siteID,
		SessionID: sessionID,
	})
	if err != nil {
		s.logger.WarnContext(ctx, errCreateSiteVisitorLink, "error", err)
	} else {
		err = s.consentsSvc.SetDefaultConsent(ctx, visitorID, siteID)
		if err != nil {
			s.logger.WarnContext(ctx, errSetDefaultConsent, "error", err)
		}
	}
}

func (s *Service) ResolveSession(ctx context.Context, sessionUUID, fingerprintHash string, geo CreateVisitorRequest, site *ent.Site) (*ent.VisitorSession, error) {
	session, err := s.sessionSvc.GetBySessionUUIDRaw(ctx, sessionUUID)
	if err != nil {
		return nil, err
	}
	if session != nil {
		return session, nil
	}

	var isBot bool
	if tracking, ok := middlewares.GetTrackingInfo(ctx); ok && tracking != nil {
		isBot = tracking.Bot
	}

	geo.FingerprintHash = fingerprintHash
	geo.SessionUUID = sessionUUID

	v, device, _, err := s.resolveVisitorAndDevice(ctx, geo, isBot)
	if err != nil {
		return nil, err
	}

	geo.FingerprintHash = fingerprintHash
	geo.SessionToken = uuid.New().String()
	session, err = s.sessionSvc.CreateSession(ctx, visitor_sessionsv1.CreateSessionInput{
		UUID:                uuid.New().String(),
		VisitorID:           int64(v.ID),
		DeviceID:            int64(device.ID),
		SessionUUID:         geo.SessionUUID,
		SessionToken:        uuid.New().String(),
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

	s.linkSiteVisitorAndConsent(ctx, int64(v.ID), int64(site.ID), int64(session.ID))

	return session, nil
}
