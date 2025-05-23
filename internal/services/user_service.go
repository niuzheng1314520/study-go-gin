package services

import (
    "github.com/gin-gonic/gin"
    "github.com/niuzheng1314520/gin/internal/models"
    "github.com/niuzheng1314520/gin/internal/repositories"
)

type UserServiceImpl struct {
    repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) services.UserService {
    return &UserServiceImpl{repo: repo}
}

func (s *UserServiceImpl) GetUserByID(ctx *gin.Context, id int64) (*models.User, error) {
    return s.repo.GetByID(ctx, id)
}
