package v1

import "time"

type DeviceResponse struct {
	UUID            string    `json:"uuid"`
	VisitorID       int64     `json:"visitor_id"`
	FingerprintHash string    `json:"fingerprint_hash"`
	DeviceType      string    `json:"device_type"`
	OS              string    `json:"os"`
	AddressIP       string    `json:"address_ip"`
	Country         string    `json:"country"`
	Region          string    `json:"region"`
	City            string    `json:"city"`
	Browser         string    `json:"browser"`
	BrowserVersion  string    `json:"browser_version"`
	Language        string    `json:"language"`
	Timezone        string    `json:"timezone"`
	FirstSeenAt     time.Time `json:"first_seen_at"`
	LastSeenAt      time.Time `json:"last_seen_at"`
}

type CreateDeviceInput struct {
	UUID            string
	VisitorID       int64
	FingerprintHash string
	DeviceType      string
	OS              string
	AddressIP       string
	Country         string
	Region          string
	City            string
	Browser         string
	BrowserVersion  string
	Language        string
	Timezone        string
	ISP             string
	ScreenWidth     int
	ScreenHeight    int
}

type UpdateDeviceInput struct {
	LastSeenAt *time.Time
	AddressIP  *string
}
