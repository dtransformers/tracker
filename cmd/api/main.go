// @title Tracking API
// @version 1.0
// @description This is a tracking service API.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:3000
// @BasePath /
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"tracking/docs"
	"tracking/internal/config"
	"tracking/internal/database"
	"tracking/internal/modules/health"
	"tracking/internal/modules/sites"
	"tracking/internal/server"

	"github.com/arl/statsviz"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type application struct {
	cfg    config.Config
	db     *database.Database
	router *gin.Engine
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("Application failed to start: %v", err)
	}
}

func run() error {
	ctx := context.Background()
	cfg, err := loadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	db, err := initDatabase(cfg)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer db.Client.Close()

	app := newApplication(cfg, db)

	return app.start(ctx)
}

func newApplication(cfg config.Config, db *database.Database) *application {
	app := &application{
		cfg: cfg,
		db:  db,
	}
	app.router = app.setupRouter()
	return app
}

// start starts the HTTP server
func (app *application) start(ctx context.Context) error {
	fmt.Printf("Starting server on %s:%d (tls=%v)\n", app.cfg.App.Host, app.cfg.App.Port, app.cfg.TLS.Enabled)
	fmt.Printf("Health check: http://%s:%d/health\n", app.cfg.App.Host, app.cfg.App.Port)

	if app.cfg.App.IsDev {
		fmt.Printf("Swagger docs: http://%s:%d/swagger/index.html\n", app.cfg.App.Host, app.cfg.App.Port)
	}

	return server.Run(ctx, app.router, app.buildServerSettings())
}

func (app *application) setupRouter() *gin.Engine {
	if app.cfg.App.IsDev {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	app.applyMiddlewares(router)

	app.registerRoutes(router)

	if app.cfg.App.IsDev {
		app.setupDevTools(router)
	}

	return router
}

func (app *application) applyMiddlewares(router *gin.Engine) {
	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{"/health"},
	}))
	router.Use(gin.Recovery())
}

func (app *application) registerRoutes(router *gin.Engine) {
	app.registerHealthRoutes(router)

	sites.RegisterV1(router, app.db.Client)
	app.registerSwaggerRoutes(router)
}

func (app *application) registerHealthRoutes(router *gin.Engine) {
	router.GET("/health", func(c *gin.Context) {
		applicationHealth := health.GetApplicationHealth(app.db.SQL)
		c.JSON(http.StatusOK, applicationHealth)
	})
}

func (app *application) registerSwaggerRoutes(router *gin.Engine) {
	docs.SwaggerInfo.Host = fmt.Sprintf("%s:%d", app.cfg.App.Host, app.cfg.App.Port)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/swagger", func(c *gin.Context) {
		c.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
}

func (app *application) setupDevTools(router *gin.Engine) {
	host := app.cfg.App.Host
	port := strconv.Itoa(app.cfg.App.Port)

	fmt.Printf("Debug tools: http://%s:%s/debug/statsviz/\n", host, port)

	mux := http.NewServeMux()
	if err := statsviz.Register(mux); err != nil {
		log.Printf("Warning: failed to register statsviz: %v", err)
		return
	}

	router.Any("/debug/statsviz/*filepath", gin.WrapH(mux))
}

func (app *application) buildServerSettings() server.Settings {
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

func loadConfig() (config.Config, error) {
	return config.Load(config.LoadOptions{
		RequireDotEnv: false,
	})
}

func initDatabase(cfg config.Config) (*database.Database, error) {
	dbConfig := buildDatabaseConfig(cfg)

	fmt.Printf("Connecting to database on %s:%d\n", cfg.DB.Host, cfg.DB.Port)

	db, err := database.Connect(dbConfig)
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	if err := db.SQL.Ping(); err != nil {
		db.Client.Close()
		return nil, fmt.Errorf("database ping failed: %w", err)
	}

	fmt.Println("Database connection established")
	return db, nil
}

func buildDatabaseConfig(cfg config.Config) database.Config {
	return database.Config{
		Host:            cfg.DB.Host,
		Port:            cfg.DB.Port,
		User:            cfg.DB.Username,
		Password:        cfg.DB.Password,
		Database:        cfg.DB.Database,
		SSLMode:         getSSLMode(cfg),
		MaxOpenConns:    getMaxOpenConnections(cfg),
		MaxIdleConns:    getMaxIdleConnections(cfg),
		ConnMaxLifetime: 30 * time.Minute,
		ConnMaxIdleTime: 10 * time.Minute,
		EnableDebug:     cfg.App.RunMode == config.Development,
	}
}

func getSSLMode(cfg config.Config) string {
	if cfg.DB.SSLMode != "" {
		return cfg.DB.SSLMode
	}
	return "disable"
}

func getMaxOpenConnections(cfg config.Config) int {
	if cfg.DB.MaxOpenConns > 0 {
		return cfg.DB.MaxOpenConns
	}
	return 10
}

func getMaxIdleConnections(cfg config.Config) int {
	if cfg.DB.MaxIdleConns > 0 {
		return cfg.DB.MaxIdleConns
	}
	return 5
}
