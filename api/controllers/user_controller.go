package controllers

import "github.com/gin-gonic/gin"

type UserController struct {
	BaseController
}

func (u *UserController) GetUserById(ctx *gin.Context) {
	userData := "user-GetUser"
	u.Success(ctx, userData)
}
