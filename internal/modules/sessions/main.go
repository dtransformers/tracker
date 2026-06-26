package session

import (
	"log/slog"
	"tracking/internal/config"
	"tracking/internal/ent"
	v1 "tracking/internal/modules/sessions/v1"
)

func RegisterV1(client *ent.Client, logger *slog.Logger, appCfg config.AppConfig) *v1.Service {
	if logger == nil {
		logger = slog.Default()
	}

	repo := v1.NewRepository(client)
	return v1.NewService(repo, logger, appCfg)
}
