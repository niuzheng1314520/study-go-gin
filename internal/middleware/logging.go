package middleware

import (
    "github.com/gin-gonic/gin"
    "go.uber.org/zap"
    "time"
)

func RequestLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()

        latency := time.Since(start)
        logger.Info("Request handled",
            zap.Int("status", c.Writer.Status()),
            zap.String("method", c.Request.Method),
            zap.String("path", c.Request.URL.Path),
            zap.Duration("latency", latency),
        )
    }
}
