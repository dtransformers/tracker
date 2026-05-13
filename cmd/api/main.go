package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
	"tracking/internal/config"
	"tracking/internal/database"
	"tracking/internal/modules/health"
	"tracking/internal/server"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
)

func setupVisualizer(router *gin.Engine, host string, port string) {
	fmt.Printf(
		"Starting debugging tools server: http://%s:%s/debug/statsviz/\n",
		host,
		port,
	)

	mux := http.NewServeMux()

	statsviz.Register(mux)

	router.Any("/debug/statsviz/*filepath", gin.WrapH(mux))
}

func setupServer(
	cfg config.Config,
	databaseConnection *database.Database,
) {
	if cfg.App.IsDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{
			"/health",
		},
	}))

	router.Use(gin.Recovery())

	setupRoutes(router, databaseConnection)

	if cfg.App.IsDev {
		setupVisualizer(router, cfg.App.Host, strconv.Itoa(cfg.App.Port))
	}
	s := server.Settings{
		Host:              cfg.App.Host,
		Port:              cfg.App.Port,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ShutdownTimeout:   10 * time.Second,

		TLS: server.TLSSettings{
			Enabled:  cfg.TLS.Enabled,
			CertFile: cfg.TLS.CertFile,
			KeyFile:  cfg.TLS.KeyFile,
		},
	}

	if err := server.Run(
		context.Background(),
		router,
		s,
	); err != nil {
		log.Fatal(err)
	}
}

func setupRoutes(
	router *gin.Engine,
	databaseConnection *database.Database,
) {
	router.GET("/health", func(c *gin.Context) {
		applicationHealth := health.GetApplicationHealth(
			databaseConnection.SQL,
		)

		c.JSON(http.StatusOK, applicationHealth)
	})
}

func setupDatabaseConnection(cfg config.Config) (*database.Database, error) {
	const MaxOpenConnections = 10
	const SSLMode = "disable"
	const MaxIdleConns = 5
	const ConnMaxLifetime = 30 * time.Minute
	const ConnMaxIdleTime = 10 * time.Minute
	var EnableDebug = cfg.App.RunMode == config.Development
	fmt.Printf("connecting to Database on %s:%d \n", cfg.DB.Host, cfg.DB.Port)
	db, err := database.Connect(database.Config{
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		User:            cfg.DB.Username,
		Password:        cfg.DB.Password,
		Database:        cfg.DB.Database,
		SSLMode:         SSLMode,
		MaxOpenConns:    MaxOpenConnections,
		MaxIdleConns:    MaxIdleConns,
		ConnMaxLifetime: ConnMaxLifetime,
		ConnMaxIdleTime: ConnMaxIdleTime,
		EnableDebug:     EnableDebug,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	if err := database.Migrate(db.Gorm); err != nil {
		return nil, fmt.Errorf("failed to migrate database: %w", err)
	}
	return db, nil
}
func main() {
	loadOptions := config.LoadOptions{
		RequireDotEnv: false,
	}
	envirmentConfig, err := config.Load(loadOptions)
	if err != nil {
		panic(`App Failed to load config ` + err.Error())
	}
	fmt.Printf(" starting server on %s:%d (tls=%v) \n", envirmentConfig.App.Host, envirmentConfig.App.Port, envirmentConfig.TLS.Enabled)

	databaseConnection, err := setupDatabaseConnection(envirmentConfig)
	if err != nil {
		panic(`App Failed to load database ` + err.Error())
	}
	setupServer(envirmentConfig, databaseConnection)
}
