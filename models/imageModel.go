package models

type Image struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	Image []byte `json:"image" gorm:"type:blob"`
}
