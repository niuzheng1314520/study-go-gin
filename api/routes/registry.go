package routes

import (
    "github.com/gin-gonic/gin"
)

type RouteHandler interface {
    RegisterPublicRoutes(group *gin.RouterGroup)
    RegisterAuthRoutes(group *gin.RouterGroup)
}

type RouteRegistry struct {
    handlers []RouteHandler
}

func NewRouteRegistry(handlers ...RouteHandler) *RouteRegistry {
    return &RouteRegistry{
        handlers: handlers,
    }
}

func (r *RouteRegistry) RegisterPublicRoutes(group *gin.RouterGroup) {
    for _, handler := range r.handlers {
        handler.RegisterPublicRoutes(group)
    }
}

func (r *RouteRegistry) RegisterAuthRoutes(group *gin.RouterGroup) {
    for _, handler := range r.handlers {
        handler.RegisterAuthRoutes(group)
    }
}
