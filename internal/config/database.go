package config

import "time"

type DatabaseConfig struct {
    Mysql map[string]MysqlConfig `yaml:"mysql"`
    Redis map[string]RedisConfig `yaml:"redis"`
}

type MysqlConfig struct {
    DSN         string        `yaml:"dsn"`
    MaxOpen     int           `yaml:"max_open"`
    MaxIdle     int           `yaml:"max_idle"`
    MaxLifetime time.Duration `yaml:"max_lifetime"`
}

type RedisConfig struct {
    Addr     string `yaml:"addr"`
    Password string `yaml:"password"`
    DB       int    `yaml:"db"`
    PoolSize int    `yaml:"pool_size"`
}
