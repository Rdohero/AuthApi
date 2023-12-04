package models

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Foto      string    `form:"foto" json:"foto"`
	Fullname  string    `gorm:"type:varchar(255)" json:"fullname"`
	Username  string    `gorm:"unique;type:varchar(255);" json:"username"`
	Saldo     int       `json:"saldo"`
	Email     string    `gorm:"unique;type:varchar(255);" json:"email"`
	Password  string    `json:"password"`
	Active    bool      `json:"is_active" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type TestTransaction struct {
	gorm.Model
	Name    string
	Balance int
}
