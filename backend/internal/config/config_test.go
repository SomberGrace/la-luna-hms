package config

import (
	"testing"
	"time"
)

func TestLoadConfigUsesLocalDatabaseDefaults(t *testing.T) {
	t.Setenv("DB_MAX_OPEN_CONNS", "")
	t.Setenv("DB_MAX_IDLE_CONNS", "")
	t.Setenv("DB_CONN_MAX_LIFETIME", "")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "")
	t.Setenv("RUN_STARTUP_MIGRATIONS", "")
	t.Setenv("SEED_DATABASE", "")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DBMaxOpenConns != 10 || cfg.DBMaxIdleConns != 2 {
		t.Fatalf("unexpected pool defaults: open=%d idle=%d", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime != 5*time.Minute || cfg.DBConnMaxIdleTime != time.Minute {
		t.Fatalf("unexpected connection lifetime defaults: lifetime=%s idle=%s", cfg.DBConnMaxLifetime, cfg.DBConnMaxIdleTime)
	}
	if !cfg.RunStartupMigrations || !cfg.SeedDatabase {
		t.Fatal("startup migrations and database seeding should default to enabled")
	}
}

func TestLoadConfigReadsRenderDatabaseSettings(t *testing.T) {
	t.Setenv("DB_MAX_OPEN_CONNS", "5")
	t.Setenv("DB_MAX_IDLE_CONNS", "2")
	t.Setenv("DB_CONN_MAX_LIFETIME", "5m")
	t.Setenv("DB_CONN_MAX_IDLE_TIME", "1m")
	t.Setenv("RUN_STARTUP_MIGRATIONS", "false")
	t.Setenv("SEED_DATABASE", "false")

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}

	if cfg.DBMaxOpenConns != 5 || cfg.DBMaxIdleConns != 2 {
		t.Fatalf("unexpected configured pool: open=%d idle=%d", cfg.DBMaxOpenConns, cfg.DBMaxIdleConns)
	}
	if cfg.DBConnMaxLifetime != 5*time.Minute || cfg.DBConnMaxIdleTime != time.Minute {
		t.Fatalf("unexpected configured durations: lifetime=%s idle=%s", cfg.DBConnMaxLifetime, cfg.DBConnMaxIdleTime)
	}
	if cfg.RunStartupMigrations || cfg.SeedDatabase {
		t.Fatal("startup migrations and database seeding should be disabled")
	}
}

func TestLoadConfigRejectsInvalidPoolConfiguration(t *testing.T) {
	t.Setenv("DB_MAX_OPEN_CONNS", "5")
	t.Setenv("DB_MAX_IDLE_CONNS", "6")

	if _, err := LoadConfig(); err == nil {
		t.Fatal("expected invalid idle connection limit to fail")
	}
}
