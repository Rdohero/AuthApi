package models

import "time"

type Payment struct {
	ID        int       `gorm:"primarykey" json:"id"`
	Payment   string    `json:"payment"`
	Deskripsi string    `json:"deskripsi"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
