package visitor_sessions

import (
	"log/slog"
	"tracking/internal/ent"
	v1 "tracking/internal/modules/visitor_sessions/v1"
)

func RegisterV1(client *ent.Client, logger *slog.Logger) *v1.Service {
	if logger == nil {
		logger = slog.Default()
	}

	repo := v1.NewRepository(client)
	return v1.NewService(repo, logger)
}
