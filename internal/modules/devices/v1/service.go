package v1

import (
	"context"
	"log/slog"
	"time"
	"tracking/internal/ent"
)

type Service struct {
	repo    Repository
	logger  *slog.Logger
	timeout time.Duration
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	if logger == nil {
		logger = slog.Default()
	}
	return &Service{
		repo:    repo,
		logger:  logger,
		timeout: 5 * time.Second,
	}
}

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

func (s *Service) GetByUUID(ctx context.Context, uuid string) (*DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	d, err := s.repo.GetByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}
	return toDeviceResponse(d), nil
}

func (s *Service) GetByVisitorID(ctx context.Context, visitorID int64) ([]*DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	devices, err := s.repo.GetByVisitorID(ctx, visitorID)
	if err != nil {
		return nil, err
	}

	res := make([]*DeviceResponse, len(devices))
	for i, d := range devices {
		res[i] = toDeviceResponse(d)
	}
	return res, nil
}

func (s *Service) GetByFingerprintHash(ctx context.Context, hash string) (*DeviceResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	d, err := s.repo.GetByFingerprintHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return toDeviceResponse(d), nil
}

func toDeviceResponse(d *ent.Device) *DeviceResponse {
	if d == nil {
		return nil
	}
	return &DeviceResponse{
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
		Language:        d.Language,
		Timezone:        d.Timezone,
		FirstSeenAt:     d.FirstSeenAt,
		LastSeenAt:      d.LastSeenAt,
	}
}
