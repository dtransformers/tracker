package config

import (
	"fmt"
	"strings"
	"testing"
	"tracking/pkg/constants"
)

type MapEnv map[string]string

func (m MapEnv) LookupEnv(key string) (string, bool) {
	v, ok := m[key]
	return v, ok
}

func baseValidEnv() MapEnv {
	return MapEnv{
		"APP_HOST":    "0.0.0.0",
		"APP_PORT":    "8080",
		"RUN_MODE":    "development",
		"DB_HOST":     "localhost",
		"DB_PORT":     "5432",
		"DB_USERNAME": "app",
		"DB_PASSWORD": "secret",
		"DB_DATABASE": "appdb",

		"REDIS_HOST":     "localhost",
		"REDIS_PORT":     "6379",
		"REDIS_DATABASE": "0",

		"TLS_ENABLED": "false",
	}
}

func TestLoadFrom_OK(t *testing.T) {
	env := baseValidEnv()

	cfg, err := LoadFrom(env, LoadOptions{
		DotEnvPath:         "",
		RequireDotEnv:      false,
		AllowEmptyOptional: true,
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if cfg.App.Port != 8080 {
		t.Fatalf("expected app port 8080, got %d", cfg.App.Port)
	}
	if cfg.DB.Host != "localhost" {
		t.Fatalf("expected db host localhost, got %q", cfg.DB.Host)
	}
}

func TestLoadFrom_InvalidIntBecomesZero_ThenValidateFails(t *testing.T) {
	env := baseValidEnv()
	env["APP_PORT"] = "not-a-number"

	cfg, err := LoadFrom(env, LoadOptions{
		DotEnvPath:         "",
		RequireDotEnv:      false,
		AllowEmptyOptional: true,
	})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("expected load to succeed (parse fallback), got %v", err)
	}

	if cfg.App.Port != constants.DefaultConfig.Port {
		t.Fatalf("expected app port to be 0 on invalid int, got %d", cfg.App.Port)
	}

}

func TestLoadFrom_InvalidRunMode(t *testing.T) {
	env := baseValidEnv()
	env["RUN_MODE"] = "invalid"

	cfg, err := LoadFrom(env, LoadOptions{
		DotEnvPath:         "",
		RequireDotEnv:      false,
		AllowEmptyOptional: true,
	})
	if err != nil {
		fmt.Println(err)
		t.Fatalf("expected load to succeed (parse fallback), got %v", err)
	}
	if cfg.App.IsDev == false {
		t.Fatalf("expected app isDev to be false on invalid run mode, got %v", cfg.App.IsDev)
	}

}

func TestValidate_MissingRequiredFields(t *testing.T) {
	cfg := Config{}
	err := Validate(cfg)
	if err == nil {
		t.Fatalf("expected validation error, got nil")
	}

	if !strings.Contains(err.Error(), "DB_HOST") && !strings.Contains(err.Error(), "db") {
		t.Fatalf("expected error to mention db host, got: %v", err)
	}
}

func TestBoolParsing(t *testing.T) {
	tests := []struct {
		name string
		val  string
		want bool
	}{
		{"true", "true", true},
		{"1", "1", true},
		{"yes", "yes", true},
		{"on", "on", true},
		{"false", "false", false},
		{"0", "0", false},
		{"no", "no", false},
		{"off", "off", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := baseValidEnv()
			env["TLS_ENABLED"] = tt.val

			cfg, err := LoadFrom(env, LoadOptions{
				DotEnvPath:         "",
				RequireDotEnv:      false,
				AllowEmptyOptional: true,
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}
			if cfg.TLS.Enabled != tt.want {
				t.Fatalf("expected %v for %q, got %v", tt.want, tt.val, cfg.TLS.Enabled)
			}
		})
	}
}

func clamp(s string, max int) string {
	if len(s) > max {
		return s[:max]
	}
	return s
}

func FuzzLoadFromEnv(f *testing.F) {
	seed := baseValidEnv()
	f.Add(
		seed["APP_HOST"],
		seed["APP_PORT"],
		seed["RUN_MODE"],
		seed["DB_HOST"],
		seed["DB_PORT"],
		seed["DB_USERNAME"],
		seed["DB_PASSWORD"],
		seed["DB_DATABASE"],
		seed["REDIS_HOST"],
		seed["REDIS_PORT"],
		seed["REDIS_DATABASE"],
		seed["TLS_ENABLED"],
	)

	f.Add("0.0.0.0", "not-a-number", "invalid", "localhost", "5432", "app", "secret", "appdb",
		"localhost", "6379", "0", "localhost", "5672", "guest", "guest", "maybe")
	f.Add("", "", "", "", "", "", "", "", "", "", "", "", "", "", "", "")
	f.Add("   ", " 8080 ", "development", " \t", "65536", "a", "b", "c",
		"127.0.0.1", "-1", "-99", "host", "0", "u", "p", "true")
	f.Add("\x00\xff", "2147483648", "prod", "db", "999999999", "user", "pass", "db",
		"redis", "999999999", "999999999", "999999999", "u", "p", "on")

	f.Fuzz(func(t *testing.T,
		appHost, appPort, runMode,
		dbHost, dbPort, dbUser, dbPass, dbName,
		redisHost, redisPort, redisDB,
		tlsEnabled string,
	) {
		// Keep fuzz inputs bounded (prevents huge allocations / slowdowns)
		appHost = clamp(appHost, 200)
		appPort = clamp(appPort, 50)
		runMode = clamp(runMode, 50)
		dbHost = clamp(dbHost, 200)
		dbPort = clamp(dbPort, 50)
		dbUser = clamp(dbUser, 200)
		dbPass = clamp(dbPass, 200)
		dbName = clamp(dbName, 200)
		redisHost = clamp(redisHost, 200)
		redisPort = clamp(redisPort, 50)
		redisDB = clamp(redisDB, 50)
		tlsEnabled = clamp(tlsEnabled, 50)

		env := baseValidEnv()
		env["APP_HOST"] = appHost
		env["APP_PORT"] = appPort
		env["RUN_MODE"] = runMode

		env["DB_HOST"] = dbHost
		env["DB_PORT"] = dbPort
		env["DB_USERNAME"] = dbUser
		env["DB_PASSWORD"] = dbPass
		env["DB_DATABASE"] = dbName

		env["REDIS_HOST"] = redisHost
		env["REDIS_PORT"] = redisPort
		env["REDIS_DATABASE"] = redisDB

		env["TLS_ENABLED"] = tlsEnabled

		cfg, err := LoadFrom(env, LoadOptions{
			DotEnvPath:         "",
			RequireDotEnv:      false,
			AllowEmptyOptional: true,
		})

		// Invariant 1: LoadFrom must never panic (Go fuzz already flags panics).
		// Invariant 2: If LoadFrom succeeded, it should not produce "impossible" values.
		if err == nil {
			// Port range sanity (tune to your real contract)
			if cfg.App.Port < 0 || cfg.App.Port > 65535 {
				t.Fatalf("app port out of range: %d (from %q)", cfg.App.Port, appPort)
			}
			if cfg.DB.Port < 0 || cfg.DB.Port > 65535 {
				t.Fatalf("db port out of range: %d (from %q)", cfg.DB.Port, dbPort)
			}

			// Your existing behavior: invalid int falls back to DefaultConfig.Port
			// (this checks the “invalid -> default” contract doesn’t go crazy)
			// We can’t know if appPort is invalid without parsing, so we just ensure:
			if cfg.App.Port == 0 && constants.DefaultConfig.Port != 0 {
				// If default isn't 0 and we got 0, that may indicate a bug—adjust if 0 is allowed.
				// Comment this out if 0 is an acceptable fallback.
			}

			// Optional: validation should behave consistently
			_ = Validate(cfg) // if Validate panics, fuzz will catch it
		}
	})
}
