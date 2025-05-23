package repositories

import (
    "github.com/gin-gonic/gin"
    "github.com/niuzheng1314520/gin/internal/database"
    "github.com/niuzheng1314520/gin/internal/models"
)

type UserRepositoryImpl struct {
    dbFactory *database.DBFactory
}

func NewUserRepository(dbFactory *database.DBFactory) repositories.UserRepository {
    return &UserRepositoryImpl{dbFactory: dbFactory}
}

func (r *UserRepositoryImpl) GetByID(ctx *gin.Context, userID int64) (*models.User, error) {
    db, err := r.dbFactory.GetMySQL("default")
    if err != nil {
        return nil, err
    }

    var user models.User
    result := db.WithContext(ctx).First(&user, userID)
    if result.Error != nil {
        return nil, result.Error
    }
    return &user, nil
}
