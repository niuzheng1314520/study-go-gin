// api/controllers/user_controller.go
package controllers

import (
    "github.com/gin-gonic/gin"
    "github.com/niuzheng1314520/gin/internal/services"
    "strconv"
)

type UserController struct {
    BaseController
    service services.UserService
}

func NewUserController(service services.UserService) *UserController {
    return &UserController{service: service}
}

func (c *UserController) RegisterPublicRoutes(group *gin.RouterGroup) {
    group.POST("/login", c.login)
}

func (c *UserController) RegisterAuthRoutes(group *gin.RouterGroup) {
    userGroup := group.Group("/users")
    {
        userGroup.GET("/:id", c.getUserByID)
        userGroup.GET("", c.listUsers)
    }
}

func (c *UserController) getUserByID(ctx *gin.Context) {
    idStr := ctx.Param("id")
    userID, err := strconv.ParseInt(idStr, 10, 64)
    if err != nil {
        c.Error(ctx, 400, 1001, "无效的用户ID")
        return
    }

    user, err := c.service.GetUserByID(ctx, userID)
    if err != nil {
        c.ErrorFromErr(ctx, 500, 1002, err)
        return
    }

    c.Success(ctx, user)
}

// 其他方法实现...
