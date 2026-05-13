package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
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

func HealthHandler(connection *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		applicationHealth := GetApplicationHealth(connection)

		w.Header().Set("Content-Type", "application/json")

		if err := json.NewEncoder(w).Encode(applicationHealth); err != nil {
			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)
			return
		}
	}
}
