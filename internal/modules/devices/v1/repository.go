package v1

import (
	"context"
	"fmt"
	"tracking/internal/ent"
	"tracking/internal/ent/device"
)

type Repository interface {
	Create(ctx context.Context, input CreateDeviceInput) (*ent.Device, error)
	GetByUUID(ctx context.Context, uuid string) (*ent.Device, error)
	GetByID(ctx context.Context, id int) (*ent.Device, error)
	GetByVisitorID(ctx context.Context, visitorID int64) ([]*ent.Device, error)
	GetByFingerprintHash(ctx context.Context, hash string) (*ent.Device, error)
	Update(ctx context.Context, uuid string, input UpdateDeviceInput) (*ent.Device, error)
	Delete(ctx context.Context, uuid string) error
}

type entRepository struct {
	client *ent.Client
}

func NewRepository(client *ent.Client) Repository {
	if client == nil {
		panic("ent client cannot be nil")
	}
	return &entRepository{client: client}
}

func (r *entRepository) Create(ctx context.Context, input CreateDeviceInput) (*ent.Device, error) {
	query := r.client.Device.
		Create().
		SetUUID(input.UUID).
		SetVisitorID(input.VisitorID).
		SetFingerprintHash(input.FingerprintHash).
		SetDeviceType(input.DeviceType).
		SetOs(input.OS).
		SetAddressIP(input.AddressIP).
		SetCountry(input.Country).
		SetRegion(input.Region).
		SetCity(input.City).
		SetBrowser(input.Browser).
		SetBrowserVersion(input.BrowserVersion).
		SetLanguage(input.Language).
		SetTimezone(input.Timezone).
		SetIsp(input.ISP).
		SetScreenWidth(input.ScreenWidth).
		SetScreenHeight(input.ScreenHeight)

	d, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return d, nil
}

func (r *entRepository) GetByUUID(ctx context.Context, uuid string) (*ent.Device, error) {
	d, err := r.client.Device.
		Query().
		Where(device.UUID(uuid)).
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get device by uuid: %w", err)
	}

	return d, nil
}

func (r *entRepository) GetByID(ctx context.Context, id int) (*ent.Device, error) {
	d, err := r.client.Device.
		Query().
		Where(device.ID(id)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get device by id: %w", err)
	}

	return d, nil
}

func (r *entRepository) GetByVisitorID(ctx context.Context, visitorID int64) ([]*ent.Device, error) {
	devices, err := r.client.Device.
		Query().
		Where(device.VisitorID(visitorID)).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get devices by visitor id: %w", err)
	}

	return devices, nil
}

func (r *entRepository) GetByFingerprintHash(ctx context.Context, hash string) (*ent.Device, error) {
	d, err := r.client.Device.
		Query().
		Where(device.FingerprintHash(hash)).
		First(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get device by fingerprint hash: %w", err)
	}

	return d, nil
}

func (r *entRepository) Update(ctx context.Context, uuid string, input UpdateDeviceInput) (*ent.Device, error) {
	query := r.client.Device.
		Update().
		Where(device.UUID(uuid))

	if input.LastSeenAt != nil {
		query.SetLastSeenAt(*input.LastSeenAt)
	}
	if input.AddressIP != nil {
		query.SetAddressIP(*input.AddressIP)
	}

	affected, err := query.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	if affected == 0 {
		return nil, &ent.NotFoundError{}
	}

	return r.GetByUUID(ctx, uuid)
}

func (r *entRepository) Delete(ctx context.Context, uuid string) error {
	affected, err := r.client.Device.
		Delete().
		Where(device.UUID(uuid)).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to delete device: %w", err)
	}

	if affected == 0 {
		return &ent.NotFoundError{}
	}

	return nil
}
