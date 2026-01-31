package config

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	cfg := Config{
		Port: "8080",
		Postgres: PostgresConfig{
			Host:            "localhost",
			Port:            "5432",
			User:            "testuser",
			Password:        "testpass",
			DBName:          "testdb",
			SSLMode:         "disable",
			MaxConns:        10,
			MaxIdleConnTime: 5 * time.Minute,
			PingTimeout:     3 * time.Second,
		},
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "localhost", cfg.Postgres.Host)
	assert.Equal(t, "5432", cfg.Postgres.Port)
	assert.Equal(t, "testuser", cfg.Postgres.User)
	assert.Equal(t, "testpass", cfg.Postgres.Password)
	assert.Equal(t, "testdb", cfg.Postgres.DBName)
	assert.Equal(t, "disable", cfg.Postgres.SSLMode)
	assert.Equal(t, int32(10), cfg.Postgres.MaxConns)
	assert.Equal(t, 5*time.Minute, cfg.Postgres.MaxIdleConnTime)
	assert.Equal(t, 3*time.Second, cfg.Postgres.PingTimeout)
}

func TestPostgresConfig(t *testing.T) {
	pgConfig := PostgresConfig{
		Host:            "db.example.com",
		Port:            "5433",
		User:            "admin",
		Password:        "secret",
		DBName:          "production",
		SSLMode:         "require",
		MaxConns:        25,
		MaxIdleConnTime: 10 * time.Minute,
		PingTimeout:     5 * time.Second,
	}

	assert.NotEmpty(t, pgConfig.Host)
	assert.NotEmpty(t, pgConfig.Port)
	assert.NotEmpty(t, pgConfig.User)
	assert.NotEmpty(t, pgConfig.Password)
	assert.NotEmpty(t, pgConfig.DBName)
	assert.NotEmpty(t, pgConfig.SSLMode)
	assert.Greater(t, pgConfig.MaxConns, int32(0))
	assert.Greater(t, pgConfig.MaxIdleConnTime, time.Duration(0))
	assert.Greater(t, pgConfig.PingTimeout, time.Duration(0))
}
