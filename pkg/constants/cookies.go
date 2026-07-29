package constants

const (
	// CookieVisitorUUID uniquely identifies a visitor across multiple sessions.
	// It is a long-lived cookie used to recognize returning users and attribute their activities consistently.
	CookieVisitorUUID = "visitor_uuid"

	// CookieSessionToken is used to track the visitor's current active session.
	// It helps in grouping interactions that occur within a continuous timeframe.
	CookieSessionToken = "session_token"

	// CookieDeviceUUID is used to track the visitor's device across multiple sessions.
	// It helps in recognizing returning users and attributing their activities consistently.
	CookieDeviceUUID = "device_uuid"
)
