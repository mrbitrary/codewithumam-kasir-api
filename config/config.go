package config

import "time"

type Config struct {
	Port     string         `mapstructure:"PORT"`
	Postgres PostgresConfig `mapstructure:"POSTGRES"`
}

type PostgresConfig struct {
	ConnString      string        `mapstructure:"CONN_STRING"`
	MaxConns        int32         `mapstructure:"MAX_CONNS"`
	MaxIdleConnTime time.Duration `mapstructure:"MAX_IDLE_CONN_TIME"`
	PingTimeout     time.Duration `mapstructure:"PING_TIMEOUT"`
}
