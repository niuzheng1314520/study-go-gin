package router

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()

	user := r.Group("/user")

	user.GET("/list", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"code": 200,
			"msg":  "success",
			"data": "list",
		})
	})

	return r
}
