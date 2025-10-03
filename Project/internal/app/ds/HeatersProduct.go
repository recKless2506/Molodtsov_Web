package ds

import "time"

type HeaterProduct struct {
	ID          uint   `gorm:"primaryKey"`
	Title       string `gorm:"size:100;not null"`
	Type        string `gorm:"size:30;not null"`
	Description string `gorm:"size:500"`
	Power       string `gorm:"not null"`
	Efficiency  string `gorm:"not null"`
	Image       string
	IsDelete    bool `gorm:"default:false;not null"`
	DeletedAt   *time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (HeaterProduct) TableName() string {
	return "heaters_products"
}
