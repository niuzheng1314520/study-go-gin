package middleware

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
)

func Recovery(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        defer func() {
            if err := recover(); err != nil {
                logger.Error("Recovered from panic",
                    zap.Any("error", err),
                    zap.String("path", c.Request.URL.Path),
                    zap.String("method", c.Request.Method),
                )

                c.AbortWithStatusJSON(500, gin.H{
                    "code": 5000,
                    "msg":  "Internal server error",
                })
            }
        }()
        c.Next()
    }
}
