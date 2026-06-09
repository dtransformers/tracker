package v1

import "time"

type SessionResponse struct {
	UUID                string     `json:"uuid"`
	VisitorID           int64      `json:"visitor_id"`
	DeviceID            int64      `json:"device_id"`
	SessionUUID         string     `json:"session_uuid"`
	SessionToken        string     `json:"session_token"`
	StartedAt           time.Time  `json:"started_at"`
	EndedAt             *time.Time `json:"ended_at,omitempty"`
	IsActive            bool       `json:"is_active"`
	LastActivityAt      time.Time  `json:"last_activity_at"`
	EntryPage           string     `json:"entry_page"`
	ExitPage            string     `json:"exit_page"`
	AcquisitionSource   string     `json:"acquisition_source"`
	AcquisitionMedium   string     `json:"acquisition_medium"`
	AcquisitionCampaign string     `json:"acquisition_campaign"`
	CampaignTerm        string     `json:"campaign_term"`
	CampaignContent     string     `json:"campaign_content"`
	CampaignID          string     `json:"campaign_id"`
	ReferrerHost        string     `json:"referrer_host"`
	CreatedAt           time.Time  `json:"created_at"`
}

type CreateSessionInput struct {
	UUID                string
	VisitorID           int64
	DeviceID            int64
	SessionUUID         string
	SessionToken        string
	StartedAt           time.Time
	EntryPage           string
	AcquisitionSource   string
	AcquisitionMedium   string
	AcquisitionCampaign string
	CampaignTerm        string
	CampaignContent     string
	CampaignID          string
	ReferrerHost        string
	EndedAt             *time.Time
	IsActive            bool
}

type UpdateSessionInput struct {
	EndedAt        *time.Time
	LastActivityAt *time.Time
	ExitPage       *string
	IsActive       *bool
}
