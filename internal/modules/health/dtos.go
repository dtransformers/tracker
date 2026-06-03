package health

import (
	"tracking/internal/database"
)

type HealthStatus struct {
	Database database.PoolStats `json:"database"`
}
