package routes

import (
    "github.com/gin-gonic/gin"
    "github.com/niuzheng1314520/gin/internal/middleware"
    ginSwagger "github.com/swaggo/gin-swagger"
    swaggerFiles "github.com/swaggo/files"
    "go.uber.org/zap"
)

func NewRouter(
    registry *RouteRegistry,
    jwtSecret string,
    logger *zap.Logger,
) *gin.Engine {
    r := gin.Default()

    // Swagger文档
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    // 全局中间件
    r.Use(
        middleware.RequestLogger(logger),
        middleware.Recovery(logger),
    )

    // API路由组
    apiGroup := r.Group("/api")
    {
        // 公共路由
        publicGroup := apiGroup.Group("")
        registry.RegisterPublicRoutes(publicGroup)

        // 认证路由
        authGroup := apiGroup.Group("")
        authGroup.Use(middleware.JWT(jwtSecret, logger))
        registry.RegisterAuthRoutes(authGroup)
    }

    return r
}
