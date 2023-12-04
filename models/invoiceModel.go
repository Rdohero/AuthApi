package models

import "time"

type Invoice struct {
	ID           uint          `gorm:"primarykey" json:"id"`
	UserID       uint          `json:"user_id"`
	TotalAmount  uint          `json:"total_amount"`
	StatusID     uint          `json:"status_id"`
	Address      string        `json:"address"`
	PaymentID    int           `json:"payment_id"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
	InvoiceItems []InvoiceItem `gorm:"foreignKey:InvoiceID"`
	Status       Status        `gorm:"foreignKey:StatusID"`
	PaymentItem  Payment       `gorm:"foreignKey:PaymentID"`
}

type Status struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
