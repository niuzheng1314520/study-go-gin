package models

type User struct {
	ID       int64  `json:"id" gorm:"primary_key"`
	Username string `json:"username" gorm:"unique;not null"`
	Password string `json:"password" gorm:"not null"`
	Email    string `json:"email" gorm:"unique;not null"`
}
