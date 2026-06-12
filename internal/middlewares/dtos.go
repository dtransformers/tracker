package middlewares

import "time"

type HMACConfig struct {
	GetSecretKey func(publicKey string) (string, error)

	MaxClockSkew time.Duration

	NonceStore NonceStore

	RequireContentInSignature bool
}

type NonceStore interface {
	CheckAndStore(publicKey, nonce string, expiry time.Time) bool
}

const SessionDetailsKey = "session_details"

type SessionDetails struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
}

type TrackingInfo struct {
	OriginalIP     string            `json:"original_ip"`
	ForwardedIPs   []string          `json:"forwarded_ips,omitempty"`
	UserAgent      string            `json:"user_agent"`
	Browser        string            `json:"browser"`
	BrowserVersion string            `json:"browser_version"`
	OS             string            `json:"os"`
	Platform       string            `json:"platform"`
	Mobile         bool              `json:"mobile"`
	Bot            bool              `json:"bot"`
	VisitorUUID    string            `json:"visitor_uuid,omitempty"`
	SessionToken   string            `json:"session_token,omitempty"`
	Cookies        map[string]string `json:"cookies,omitempty"`
	Referer        string            `json:"referer,omitempty"`
	AcceptLanguage string            `json:"accept_language,omitempty"`
	Host           string            `json:"host"`
	Timestamp      time.Time         `json:"timestamp"`
	Country        string            `json:"country,omitempty"`
	Region         string            `json:"region,omitempty"`
	City           string            `json:"city,omitempty"`
	Timezone       string            `json:"timezone,omitempty"`
}

const TrackingInfoKey = "tracking_info"
