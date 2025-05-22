package database

import (
	"context"
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/niuzheng1314520/gin/internal/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
)

type DBFactory struct {
	mysqlDbs map[string]*gorm.DB
	redisDbs map[string]redis.UniversalClient
	mu       sync.RWMutex
}

func NewDbFactory(cfg *config.DatabaseConfig) (*DBFactory, error) {
	factory := &DBFactory{
		mysqlDbs: make(map[string]*gorm.DB),
		redisDbs: make(map[string]redis.UniversalClient),
	}

	if err := factory.initMysql(cfg); err != nil {
		return nil, err
	}
	if err := factory.initRedis(cfg); err != nil {
		return nil, err
	}
	return factory, nil
}

func (f *DBFactory) initMysql(cfg *config.DatabaseConfig) error {
	defaultDB, err := gorm.Open(mysql.Open(cfg.Mysql.Default.DSN), &gorm.Config{})
	if err != nil {
		return errors.New("初始化默认 MySQL 连接失败: " + err.Error())
	}
	sqlDB, err := defaultDB.DB()
	if err != nil {
		return errors.New("获取默认 MySQL 连接池失败: " + err.Error())
	}
	sqlDB.SetMaxOpenConns(cfg.Mysql.Default.MaxOpen)
	sqlDB.SetMaxIdleConns(cfg.Mysql.Default.MaxIdle)
	sqlDB.SetConnMaxLifetime(cfg.Mysql.Default.MaxLifetime)

	f.mu.Lock()
	f.mysqlDbs["default"] = defaultDB
	f.mu.Unlock()

	return nil
}

func (f *DBFactory) initRedis(cfg *config.DatabaseConfig) error {
	defaultClient := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Default.Addr,
		Password: cfg.Redis.Default.Password,
		DB:       cfg.Redis.Default.DB,
		PoolSize: cfg.Redis.Default.PoolSize,
	})

	if _, err := defaultClient.Ping(context.Background()).Result(); err != nil {
		return errors.New("初始化默认 Redis 连接失败: " + err.Error())
	}

	f.mu.Lock()
	f.redisDbs["default"] = defaultClient
	f.mu.Unlock()

	return nil
}

func (f *DBFactory) GetMysql(name string) (*gorm.DB, bool) {
	f.mu.RLock()
	db, ok := f.mysqlDbs[name]
	f.mu.RUnlock()
	return db, ok
}

func (f *DBFactory) Close() error {
	f.mu.Lock()
	defer f.mu.Unlock()
	for _, db := range f.mysqlDbs {
		sqlDB, _ := db.DB()
		_ = sqlDB.Close()
	}
	for _, client := range f.redisDbs {
		_ = client.Close()
	}
	return nil
}
