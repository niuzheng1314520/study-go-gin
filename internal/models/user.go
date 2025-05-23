package models

type User struct {
    ID         int64  `json:"id" gorm:"primary_key"`
    Username   string `json:"username" gorm:"unique;not null"`
    Password   string `json:"password" gorm:"not null"`
    Status     int    `json:"status" gorm:"not null"`
    Realname   string `json:"realname" gorm:"not null"`
    Suffix     string `json:"suffix" gorm:"not null"`
    ParentUser string `json:"parent_user" gorm:"not null"`
    RootUser   string `json:"root_user" gorm:"not null"`
    ClientType string `json:"client_type" gorm:"not null"`
}
