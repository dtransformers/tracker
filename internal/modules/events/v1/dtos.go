package v1

type CreateEventRequest struct {
	TrackingEvent
}

type CreatePageViewRequest struct {
	SessionToken    string `json:"session_token"`
	FingerprintHash string `json:"fingerprint_hash"`
	URL             string `json:"url" binding:"required"`
	Title           string `json:"title"`
	Referrer        string `json:"referrer"`
}

type CreateEventInput struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
	Event     *TrackingEvent
}

type CreatePageViewInput struct {
	SessionID int64
	VisitorID int64
	DeviceID  int64
	URL       string
	Title     string
	Referrer  string
}
