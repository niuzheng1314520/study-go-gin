package main

import (
    "go.uber.org/zap"
    "github.com/niuzheng1314520/gin/api/routes"
    "github.com/niuzheng1314520/gin/internal/config"
    "github.com/niuzheng1314520/gin/internal/database"
    "github.com/niuzheng1314520/gin/internal/migrations"
    "github.com/niuzheng1314520/gin/internal/repositories"
    "github.com/niuzheng1314520/gin/internal/services"
    "github.com/niuzheng1314520/gin/api/controllers"
)

func main() {
    // 1. 初始化日志
    logger, err := zap.NewProduction()
    if err != nil {
        panic("failed to initialize logger: " + err.Error())
    }
    defer func() {
        if err := logger.Sync(); err != nil {
            // 在容器环境中Sync可能会失败，可以忽略
            logger.Warn("failed to sync logger", zap.Error(err))
        }
    }()

    // 2. 加载配置
    cfg, err := config.LoadConfig()
    if err != nil {
        logger.Fatal("加载配置失败", zap.Error(err))
    }

    // 3. 初始化数据库
    dbFactory, err := database.NewDBFactory(&cfg.Database, logger)
    if err != nil {
        logger.Fatal("数据库初始化失败",
            zap.Error(err),
            zap.Any("config", cfg.Database),
        )
    }

    // 4. 执行数据库迁移
    if db, err := dbFactory.GetMySQL("default"); err == nil {
        if err := migrations.AutoMigrate(db); err != nil {
            logger.Warn("数据库迁移失败",
                zap.Error(err),
                zap.String("db", "default"),
            )
        }
    } else {
        logger.Warn("获取数据库连接失败",
            zap.Error(err),
        )
    }

    // 5. 初始化应用组件
    userRepo := repositories.NewUserRepository(dbFactory)
    userService := services.NewUserService(userRepo)
    userCtrl := controllers.NewUserController(userService)

    // 6. 创建路由注册中心
    registry := routes.NewRouteRegistry(userCtrl)

    // 7. 创建路由实例
    r := routes.NewRouter(
        registry,
        cfg.JWT.Secret,
        logger,
    )

    // 8. 启动服务
    logger.Info("服务启动中",
        zap.String("port", cfg.Server.Port),
        zap.String("env", cfg.Server.Env),
    )
    if err := r.Run(":" + cfg.Server.Port); err != nil {
        logger.Fatal("服务启动失败",
            zap.Error(err),
            zap.String("port", cfg.Server.Port),
        )
    }
}
