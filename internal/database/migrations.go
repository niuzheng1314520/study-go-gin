package database

import (
    "github.com/niuzheng1314520/gin/internal/models"
    "gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
    return db.AutoMigrate(
        &models.User{},
    )
}
