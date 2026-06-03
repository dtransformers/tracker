package health

import (
	"database/sql"
	"tracking/internal/database"
)


func GetApplicationHealth(connection *sql.DB) HealthStatus {

	dbStatus := database.Stats(connection)

	return HealthStatus{
		Database: dbStatus,
	}

}
