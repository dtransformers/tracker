package health

import (
	"database/sql"
	"tracking/internal/database"
)

type HealthStatus struct {
	Database database.PoolStats `json:"database"`
}

func GetApplicationHealth(connection *sql.DB) HealthStatus {

	dbStatus := database.Stats(connection)

	return HealthStatus{
		Database: dbStatus,
	}

}
