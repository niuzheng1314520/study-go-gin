package repositories

import (
	"github.com/gin-gonic/gin"
	"github.com/niuzheng1314520/gin/internal/models"
)

type UserRepository interface {
	GetByID(ctx *gin.Context, userId int64) (*models.User, error)
}
