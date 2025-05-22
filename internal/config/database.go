package config

import "time"

type DatabaseConfig struct {
	Mysql struct {
		Default struct {
			DSN         string        `yaml:"dsn"`
			MaxOpen     int           `yaml:"max_open"`
			MaxIdle     int           `yaml:"max_idle"`
			MaxLifetime time.Duration `yaml:"max_lifetime"`
		} `yaml:"default"`
	} `yaml:"mysql"`

	Redis struct {
		Default struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
			PoolSize int    `yaml:"pool_size"`
		} `yaml:"default"`
	} `yaml:"redis"`
}
