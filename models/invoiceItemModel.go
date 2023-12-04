package models

import "time"

type InvoiceItem struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	InvoiceID  uint      `json:"invoice_id"`
	ProductID  uint      `json:"product_id"`
	Quantity   uint      `json:"quantity"`
	TotalPrice uint      `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	Product    Product   `gorm:"foreignKey:ProductID"`
}
