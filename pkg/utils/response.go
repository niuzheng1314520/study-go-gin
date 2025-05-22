package utils

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
	Count int64       `json:"count,omitempty"`
}

// Success 返回成功响应
func Success(ctx *gin.Context, data interface{}) {
	resp := Response{
		Code: 0,
		Msg:  "success",
		Data: data,
	}
	ctx.JSON(http.StatusOK, resp)
}

// SuccessWithCount 返回带计数的成功响应（用于分页等场景）
func SuccessWithCount(ctx *gin.Context, data interface{}, count int64) {
	resp := Response{
		Code:  0,
		Msg:   "success",
		Data:  data,
		Count: count,
	}
	ctx.JSON(http.StatusOK, resp)
}

// Error 返回错误响应
func Error(ctx *gin.Context, httpCode, businessCode int, msg string) {
	resp := Response{
		Code: businessCode,
		Msg:  msg,
	}
	ctx.JSON(httpCode, resp)
}

// Errorf 格式化错误信息并返回错误响应
func Errorf(ctx *gin.Context, httpCode, businessCode int, format string, args ...interface{}) {
	resp := Response{
		Code: businessCode,
		Msg:  fmt.Sprintf(format, args...),
	}
	ctx.JSON(httpCode, resp)
}

// ErrorFromErr 从 error 对象生成错误响应
func ErrorFromErr(ctx *gin.Context, httpCode, businessCode int, err error) {
	resp := Response{
		Code: businessCode,
		Msg:  err.Error(),
	}
	ctx.JSON(httpCode, resp)
}
