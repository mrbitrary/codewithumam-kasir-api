package config

import "time"

type Config struct {
	Port     string         `mapstructure:"PORT"`
	Postgres PostgresConfig `mapstructure:"POSTGRES"`
}

type PostgresConfig struct {
	Host            string        `mapstructure:"HOST"`
	Port            string        `mapstructure:"PORT"`
	User            string        `mapstructure:"USER"`
	Password        string        `mapstructure:"PASSWORD"`
	DBName          string        `mapstructure:"DB_NAME"`
	SSLMode         string        `mapstructure:"SSL_MODE"`
	MaxConns        int32         `mapstructure:"MAX_CONNS"`
	MaxIdleConnTime time.Duration `mapstructure:"MAX_IDLE_CONN_TIME"`
	PingTimeout     time.Duration `mapstructure:"PING_TIMEOUT"`
}
