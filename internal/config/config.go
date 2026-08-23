package config

import "time"

type Config struct {
	ConnURL           string
	AppPort           string
	DBMinConn         int
	DBMaxConn         int
	DBMaxConnLifetime time.Duration
}

func NewConfig(url, port string, minConn, maxConn int, maxLifetime time.Duration) *Config {
	return &Config{
		ConnURL:           url,
		AppPort:           port,
		DBMinConn:         minConn,
		DBMaxConn:         maxConn,
		DBMaxConnLifetime: maxLifetime,
	}
}
