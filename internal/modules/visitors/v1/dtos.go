package v1

import (
	"time"

	"tracking/internal/ent"
	devicesv1 "tracking/internal/modules/devices/v1"
	visitor_sessionsv1 "tracking/internal/modules/visitor_sessions/v1"
)

type CreateVisitorRequest struct {
	PublicKey string `json:"-"`

	VisitorUUID string `json:"visitor_uuid" swaggerignore:"true"`

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
	SessionToken    string    `json:"session_token"`
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

type VerifyVisitorData struct {
	IsValid bool                                        `json:"is_valid"`
	Session *visitor_sessionsv1.SessionResponse         `json:"session"`
	Device  *devicesv1.DeviceResponse                   `json:"device"`
}

type VerifyVisitorAPIResponse struct {
	Success bool               `json:"success"`
	Data    *VerifyVisitorData `json:"data"`
	Version string             `json:"version"`
	Error   *string            `json:"error"`
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
	SessionToken    string                 `json:"session_token"`
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
	SessionToken string   `json:"session_token"`
	VisitorUUID string    `json:"visitor_uuid"`
	DeviceUUID  string    `json:"device_uuid"`
	URL         string    `json:"url"`
	Title       string    `json:"title,omitempty"`
	Referrer    string    `json:"referrer,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

type CreateSessionInput struct {
	UUID            string
	AnonymousID     string
	PrimaryUserID   *string
	IsBot           bool
	IsVerifiedHuman bool
}

type UpdateVisitorInput struct {
	PrimaryUserID   *string
	LastSeenAt      *time.Time
	IsBot           *bool
	IsVerifiedHuman *bool
}

type CreateSiteVisitorInput struct {
	UUID      string
	VisitorID int64
	SessionID int64
	SiteID    int64
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
