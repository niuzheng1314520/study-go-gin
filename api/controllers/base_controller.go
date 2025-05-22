package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/niuzheng1314520/gin/pkg/utils"
)

type BaseController struct{}

func (b *BaseController) Success(ctx *gin.Context, data interface{}) {
	utils.Success(ctx, data)
}

func (b *BaseController) SuccessWithCount(ctx *gin.Context, data interface{}, count int64) {
	utils.SuccessWithCount(ctx, data, count)
}

func (b *BaseController) Error(ctx *gin.Context, httpCode, businessCode int, msg string) {
	utils.Error(ctx, httpCode, businessCode, msg)
}

func (b *BaseController) Errorf(ctx *gin.Context, httpCode, businessCode int, format string, args ...interface{}) {
	utils.Errorf(ctx, httpCode, businessCode, format, args...)
}

func (b *BaseController) ErrorFromErr(ctx *gin.Context, httpCode, businessCode int, err error) {
	utils.ErrorFromErr(ctx, httpCode, businessCode, err)
}
