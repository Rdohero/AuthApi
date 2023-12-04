package models

import "time"

type Product struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `json:"user_id"`
	Image       string    `json:"image"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Location    string    `json:"location"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	User        User      `gorm:"foreignKey:UserID"`
}
