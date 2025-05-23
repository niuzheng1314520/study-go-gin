package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/niuzheng1314520/gin/pkg/utils"
    "go.uber.org/zap" // 添加zap日志依赖
    "net/http"
)

// 修改函数签名，添加logger参数
func JWT(secret string, logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            logger.Warn("Missing authorization header") // 添加日志记录
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "code":    http.StatusUnauthorized,
                "message": "授权凭证缺失",
            })
            return
        }

        userID, err := utils.ParseToken(token, secret)
        if err != nil {
            logger.Warn("Invalid token", // 添加日志记录
                zap.String("token", token),
                zap.Error(err),
            )
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
                "code":    http.StatusUnauthorized,
                "message": "无效的访问令牌",
            })
            return
        }

        logger.Info("User authenticated", // 添加成功日志
            zap.Int64("userID", userID),
        )
        c.Set("userID", userID)
        c.Next()
    }
}
