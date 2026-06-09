package constants

// Cookies naming and meanings
const (
	// CookieSessionID is used to track the visitor's current active session.
	// It helps in grouping interactions that occur within a continuous timeframe.
	CookieSessionID = "session_id"

	// CookieVisitorUUID uniquely identifies a visitor across multiple sessions.
	// It is a long-lived cookie used to recognize returning users and attribute their activities consistently.
	CookieVisitorUUID = "visitor_uuid"
)
