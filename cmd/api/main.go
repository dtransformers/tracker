// @title Tracking API
// @version 1.0
// @description This is a tracking service API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
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
	"tracking/internal/modules/sites"
	"tracking/internal/server"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "tracking/docs"
)

type application struct {
	cfg config.Config
	db  *database.Database
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("App failed to start: %v", err)
	}
}

func run() error {
	cfg, err := config.Load(config.LoadOptions{
		RequireDotEnv: false,
	})
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	db, err := initDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}

	app := &application{
		cfg: cfg,
		db:  db,
	}

	router := app.setupRouter()

	fmt.Printf("starting server on %s:%d (tls=%v) \n", cfg.App.Host, cfg.App.Port, cfg.TLS.Enabled)
	fmt.Printf("Check Server Health On http://%s:%d/health \n", cfg.App.Host, cfg.App.Port)

	return server.Run(context.Background(), router, app.serverSettings())
}

func (app *application) setupRouter() *gin.Engine {
	if app.cfg.App.IsDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global Middlewares
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	router.Use(gin.Recovery())

	// Register Routes
	app.registerRoutes(router)

	// Setup Visualizer in development mode
	if app.cfg.App.IsDev {
		app.setupVisualizer(router)
	}

	return router
}

func (app *application) registerRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		applicationHealth := health.GetApplicationHealth(app.db.SQL)
		c.JSON(http.StatusOK, applicationHealth)
	})

	sites.RegisterV1(router, app.db.Client)

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

func (app *application) setupVisualizer(router *gin.Engine) {
	host := app.cfg.App.Host
	port := strconv.Itoa(app.cfg.App.Port)

	fmt.Printf("Starting debugging tools server: http://%s:%s/debug/statsviz/\n", host, port)

	mux := http.NewServeMux()
	statsviz.Register(mux)

	router.Any("/debug/statsviz/*filepath", gin.WrapH(mux))
}

func (app *application) serverSettings() server.Settings {
	return server.Settings{
		Host:              app.cfg.App.Host,
		Port:              app.cfg.App.Port,
		ReadTimeout:       15 * time.Second,
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ShutdownTimeout:   10 * time.Second,
		TLS: server.TLSSettings{
			Enabled:  app.cfg.TLS.Enabled,
			CertFile: app.cfg.TLS.CertFile,
			KeyFile:  app.cfg.TLS.KeyFile,
		},
	}
}

func initDatabase(cfg config.Config) (*database.Database, error) {
	const (
		MaxOpenConnections = 10
		MaxIdleConns       = 5
		SSLMode            = "disable"
		ConnMaxLifetime    = 30 * time.Minute
		ConnMaxIdleTime    = 10 * time.Minute
	)

	enableDebug := cfg.App.RunMode == config.Development

	fmt.Printf("connecting to Database on %s:%d\n", cfg.DB.Host, cfg.DB.Port)

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
		EnableDebug:     enableDebug,
	})
	if err != nil {
		return nil, err
	}

	return db, nil
}
