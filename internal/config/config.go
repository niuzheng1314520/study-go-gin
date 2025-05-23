// internal/config/config.go
package config

import (
    "github.com/spf13/viper"
    "path/filepath"
    "runtime"
)

type AppConfig struct {
    Database DatabaseConfig `yaml:"database"`
    JWT      JWTConfig      `yaml:"jwt"`
    Server   ServerConfig   `yaml:"server"`
}

func LoadConfig() (*AppConfig, error) {
    _, filename, _, _ := runtime.Caller(0)
    dir := filepath.Dir(filename)

    viper.AddConfigPath(dir)
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg AppConfig
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }
    return &cfg, nil
}
