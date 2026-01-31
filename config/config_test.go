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
			ConnString:      "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
			MaxConns:        10,
			MaxIdleConnTime: 5 * time.Minute,
			PingTimeout:     3 * time.Second,
		},
	}

	assert.Equal(t, "8080", cfg.Port)
	assert.Equal(t, "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable", cfg.Postgres.ConnString)
	assert.Equal(t, int32(10), cfg.Postgres.MaxConns)
	assert.Equal(t, 5*time.Minute, cfg.Postgres.MaxIdleConnTime)
	assert.Equal(t, 3*time.Second, cfg.Postgres.PingTimeout)
}

func TestPostgresConfig(t *testing.T) {
	pgConfig := PostgresConfig{
		ConnString:      "postgres://admin:secret@db.example.com:5433/production?sslmode=require",
		MaxConns:        25,
		MaxIdleConnTime: 10 * time.Minute,
		PingTimeout:     5 * time.Second,
	}

	assert.NotEmpty(t, pgConfig.ConnString)
	assert.Greater(t, pgConfig.MaxConns, int32(0))
	assert.Greater(t, pgConfig.MaxIdleConnTime, time.Duration(0))
	assert.Greater(t, pgConfig.PingTimeout, time.Duration(0))
}
