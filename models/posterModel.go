package models

import "time"

type Poster struct {
	ID        uint `gorm:"primarykey" json:"id"`
	Poster    string
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
