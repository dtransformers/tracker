/**
* Config Module
* Load the app env from os or .env file
* Validate the app env stracture
* Return the app env
*
 */
package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
	"tracking/pkg/constants"

	"github.com/joho/godotenv"
)

type EnvProvider interface {
	LookupEnv(key string) (string, bool)
}

type OSEnv struct{}

func (OSEnv) LookupEnv(key string) (string, bool) { return os.LookupEnv(key) }

type LoadOptions struct {
	DotEnvPath         string
	RequireDotEnv      bool
	AllowEmptyOptional bool
}

type Config struct {
	App   AppConfig
	DB    DBConfig
	Redis RedisConfig
	TLS   TLSConfig
}
type RunMode string

const (
	Development RunMode = "development"
	Production  RunMode = "production"
	Testing     RunMode = "testing"
)

type TLSConfig struct {
	Enabled  bool
	CertFile string
	KeyFile  string
}
type AppConfig struct {
	Host    string
	Port    int
	RunMode RunMode
	IsDev   bool
}

type DBConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	Database        string
	MaxOpenConns    int           "default: 10"
	MaxIdleConns    int           "default: 5"
	ConnMaxLifetime time.Duration "default: 30m"
	ConnMaxIdleTime time.Duration "default: 10m"
	SSLMode         string
}

type RedisConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	Database int
}

func Load(opts LoadOptions) (Config, error) {
	return LoadFrom(OSEnv{}, opts)
}
func LoadFrom(env EnvProvider, opts LoadOptions) (Config, error) {
	if err := loadDotEnv(opts); err != nil {
		return Config{}, err
	}

	cfg := Config{
		App:   loadApp(env, opts),
		DB:    loadDB(env, opts),
		Redis: loadRedis(env, opts),
		TLS:   loadTLS(env, opts),
	}

	if err := Validate(cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func Validate(cfg Config) error {
	var errs []error

	if cfg.App.Port <= 0 || cfg.App.Port > 65535 {
		errs = append(errs, fmt.Errorf("APP_PORT out of range: %d", cfg.App.Port))
	}
	if strings.TrimSpace(cfg.App.Host) == "" {
		errs = append(errs, errors.New("APP_HOST is required"))
	}

	if strings.TrimSpace(cfg.DB.Host) == "" {
		errs = append(errs, errors.New("DB_HOST is required"))
	}
	if cfg.DB.Port <= 0 || cfg.DB.Port > 65535 {
		errs = append(errs, fmt.Errorf("DB_PORT out of range: %d", cfg.DB.Port))
	}
	if strings.TrimSpace(cfg.DB.Username) == "" {
		errs = append(errs, errors.New("DB_USERNAME is required"))
	}
	if strings.TrimSpace(cfg.DB.Database) == "" {
		errs = append(errs, errors.New("DB_DATABASE is required"))
	}
	if strings.TrimSpace(cfg.DB.Password) == "" {
		errs = append(errs, errors.New("DB_PASSWORD is required"))
	}

	if strings.TrimSpace(cfg.Redis.Host) == "" {
		errs = append(errs, errors.New("REDIS_HOST is required"))
	}
	if cfg.Redis.Port <= 0 || cfg.Redis.Port > 65535 {
		errs = append(errs, fmt.Errorf("REDIS_PORT out of range: %d", cfg.Redis.Port))
	}
	if cfg.Redis.Database < 0 {
		errs = append(errs, fmt.Errorf("REDIS_DATABASE must be >= 0, got: %d", cfg.Redis.Database))
	}

	return joinErrors(errs)
}

func loadApp(env EnvProvider, opts LoadOptions) AppConfig {
	port, _err := mustValidPort(env, "APP_PORT")
	if _err != nil {
		panic(_err)
	}
	runMode := mustString(env, "RUN_MODE")
	if runMode != "development" && runMode != "production" && runMode != "testing" {
		runMode = "development"
	}
	isDev := runMode == "development"
	return AppConfig{
		Host:    mustString(env, "APP_HOST"),
		Port:    port,
		RunMode: RunMode(runMode),
		IsDev:   isDev,
	}
}

func loadDB(env EnvProvider, opts LoadOptions) DBConfig {
	return DBConfig{
		Host:     mustString(env, "DB_HOST"),
		Port:     mustInt(env, "DB_PORT"),
		Username: mustString(env, "DB_USERNAME"),
		Password: mustString(env, "DB_PASSWORD"),
		Database: mustString(env, "DB_DATABASE"),
	}
}

func loadRedis(env EnvProvider, opts LoadOptions) RedisConfig {
	return RedisConfig{
		Host:     mustString(env, "REDIS_HOST"),
		Port:     mustInt(env, "REDIS_PORT"),
		Database: intOrDefault(env, "REDIS_DATABASE", 0),
		Username: stringOrDefault(env, "REDIS_USERNAME", "", opts.AllowEmptyOptional),
		Password: stringOrDefault(env, "REDIS_PASSWORD", "", opts.AllowEmptyOptional),
	}
}

func loadTLS(env EnvProvider, opts LoadOptions) TLSConfig {
	return TLSConfig{
		Enabled:  boolOrDefault(env, "TLS_ENABLED", false),
		CertFile: stringOrDefault(env, "TLS_CERT_FILE", "", opts.AllowEmptyOptional),
		KeyFile:  stringOrDefault(env, "TLS_KEY_FILE", "", opts.AllowEmptyOptional),
	}
}

func mustString(env EnvProvider, key string) string {
	v, ok := env.LookupEnv(key)
	if !ok {
		return ""
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return ""
	}
	return v
}

func mustValidPort(env EnvProvider, key string) (int, error) {
	var ErrPortMustBeDefined = errors.New("port must be defined at least")

	v, ok := env.LookupEnv(key)
	if !ok {
		return constants.DefaultConfig.Port, ErrPortMustBeDefined
	}
	v = strings.TrimSpace(v)
	if !ok {
		return constants.DefaultConfig.Port, ErrPortMustBeDefined
	}
	if v == "" {
		return constants.DefaultConfig.Port, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return constants.DefaultConfig.Port, nil
	}
	return n, nil
}

func mustInt(env EnvProvider, key string) int {
	v, ok := env.LookupEnv(key)
	if !ok {
		return 0
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return 0
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return n
}

func intOrDefault(env EnvProvider, key string, def int) int {
	v, ok := env.LookupEnv(key)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}

func stringOrDefault(env EnvProvider, key, def string, allowEmpty bool) string {
	v, ok := env.LookupEnv(key)
	if !ok {
		return def
	}
	v = strings.TrimSpace(v)
	if v == "" && !allowEmpty {
		return def
	}
	return v
}

func boolOrDefault(env EnvProvider, key string, def bool) bool {
	v, ok := env.LookupEnv(key)
	if !ok {
		return def
	}

	v = strings.TrimSpace(strings.ToLower(v))
	if v == "" {
		return def
	}

	switch v {
	case "1", "true", "yes", "y", "on":
		return true
	case "0", "false", "no", "n", "off":
		return false
	default:
		return def
	}
}

func loadDotEnv(opts LoadOptions) error {
	var err error
	if opts.DotEnvPath != "" {
		err = godotenv.Load(opts.DotEnvPath)
	} else {
		err = godotenv.Load()
	}

	if err != nil && opts.RequireDotEnv {
		return fmt.Errorf("failed to load .env: %w", err)
	}
	return nil
}

func joinErrors(errs []error) error {
	if len(errs) == 0 {
		return nil
	}
	return errors.Join(errs...)
}
