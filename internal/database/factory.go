package database

import (
    "context"
    "errors"
    "fmt"
    "sync"
    "time"

    "github.com/go-redis/redis/v8"
    "github.com/niuzheng1314520/gin/internal/config"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
)

// 定义标准错误类型
var (
    ErrDuplicateConnection = errors.New("duplicate connection")
    ErrConnectionNotFound  = errors.New("connection not found")
    ErrInvalidConfig       = errors.New("invalid configuration")
)

type DBFactory struct {
    mysqlConns map[string]*gorm.DB
    redisConns map[string]*redis.Client
    mu         sync.RWMutex
    logger     Logger
}

// Logger 日志接口抽象
type Logger interface {
    Warnf(format string, args ...interface{})
}

// NewDBFactory 创建数据库工厂实例
func NewDBFactory(cfg *config.DatabaseConfig, logger Logger) (*DBFactory, error) {
    factory := &DBFactory{
        mysqlConns: make(map[string]*gorm.DB),
        redisConns: make(map[string]*redis.Client),
        logger:     logger,
    }

    // 初始化所有MySQL连接
    for name, mysqlCfg := range cfg.Mysql {
        if err := factory.initMySQL(name, mysqlCfg); err != nil {
            return nil, fmt.Errorf("mysql init failed for %s: %w", name, err)
        }
    }

    // 初始化所有Redis连接
    for name, redisCfg := range cfg.Redis {
        if err := factory.initRedis(name, redisCfg); err != nil {
            factory.logger.Warnf("Redis init failed for %s: %v", name, err)
        }
    }

    if len(factory.mysqlConns) == 0 {
        return nil, errors.New("no valid mysql connections initialized")
    }

    return factory, nil
}

// initMySQL 初始化MySQL连接
func (f *DBFactory) initMySQL(name string, cfg config.MysqlConfig) error {
    // 参数验证
    if cfg.DSN == "" {
        return fmt.Errorf("%w: mysql dsn is empty", ErrInvalidConfig)
    }

    // 设置默认值
    if cfg.MaxIdle <= 0 {
        cfg.MaxIdle = 10
    }
    if cfg.MaxOpen <= 0 {
        cfg.MaxOpen = 100
    }
    if cfg.MaxLifetime <= 0 {
        cfg.MaxLifetime = time.Hour
    }

    // 创建数据库连接
    db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{
        DisableForeignKeyConstraintWhenMigrating: true,
        NowFunc: func() time.Time {
            return time.Now().UTC()
        },
    })
    if err != nil {
        return fmt.Errorf("mysql connection failed: %w", err)
    }

    // 配置连接池
    sqlDB, err := db.DB()
    if err != nil {
        return fmt.Errorf("get sql.DB failed: %w", err)
    }
    sqlDB.SetMaxIdleConns(cfg.MaxIdle)
    sqlDB.SetMaxOpenConns(cfg.MaxOpen)
    sqlDB.SetConnMaxLifetime(cfg.MaxLifetime)

    // 测试连接
    if err := sqlDB.Ping(); err != nil {
        _ = sqlDB.Close()
        return fmt.Errorf("mysql ping failed: %w", err)
    }

    // 存储连接
    f.mu.Lock()
    defer f.mu.Unlock()

    if _, exists := f.mysqlConns[name]; exists {
        _ = sqlDB.Close()
        return fmt.Errorf("%w: mysql %s", ErrDuplicateConnection, name)
    }

    f.mysqlConns[name] = db
    return nil
}

// initRedis 初始化Redis连接
func (f *DBFactory) initRedis(name string, cfg config.RedisConfig) error {
    // 参数验证
    if cfg.Addr == "" {
        return fmt.Errorf("%w: redis addr is empty", ErrInvalidConfig)
    }

    // 设置默认值
    if cfg.PoolSize <= 0 {
        cfg.PoolSize = 10
    }

    // 创建客户端
    client := redis.NewClient(&redis.Options{
        Addr:     cfg.Addr,
        Password: cfg.Password,
        DB:       cfg.DB,
        PoolSize: cfg.PoolSize,
    })

    // 测试连接
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := client.Ping(ctx).Err(); err != nil {
        _ = client.Close()
        return fmt.Errorf("redis ping failed: %w", err)
    }

    // 存储连接
    f.mu.Lock()
    defer f.mu.Unlock()

    if _, exists := f.redisConns[name]; exists {
        _ = client.Close()
        return fmt.Errorf("%w: redis %s", ErrDuplicateConnection, name)
    }

    f.redisConns[name] = client
    return nil
}

// GetMySQL 获取MySQL连接
func (f *DBFactory) GetMySQL(name string) (*gorm.DB, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()

    db, ok := f.mysqlConns[name]
    if !ok {
        return nil, fmt.Errorf("%w: mysql %s", ErrConnectionNotFound, name)
    }
    return db, nil
}

// GetRedis 获取Redis连接
func (f *DBFactory) GetRedis(name string) (*redis.Client, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()

    client, ok := f.redisConns[name]
    if !ok || client == nil {
        return nil, fmt.Errorf("%w: redis %s", ErrConnectionNotFound, name)
    }
    return client, nil
}

// Close 关闭所有连接
func (f *DBFactory) Close() error {
    f.mu.Lock()
    defer f.mu.Unlock()

    var errs []error

    // 关闭MySQL连接
    for name, db := range f.mysqlConns {
        sqlDB, err := db.DB()
        if err != nil {
            errs = append(errs, fmt.Errorf("get mysql pool failed (%s): %w", name, err))
            continue
        }
        if err := sqlDB.Close(); err != nil {
            errs = append(errs, fmt.Errorf("close mysql failed (%s): %w", name, err))
        }
    }

    // 关闭Redis连接
    for name, client := range f.redisConns {
        if err := client.Close(); err != nil {
            errs = append(errs, fmt.Errorf("close redis failed (%s): %w", name, err))
        }
    }

    // 清空连接池
    f.mysqlConns = make(map[string]*gorm.DB)
    f.redisConns = make(map[string]*redis.Client)

    if len(errs) > 0 {
        return fmt.Errorf("close errors: %v", errs)
    }
    return nil
}
