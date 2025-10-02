package ds

import (
	"time"
)

type HeatersProductRequest struct {
	ID                 uint      `gorm:"primaryKey"`
	Status             string    `gorm:"not null"`
	CreatedAt          time.Time `gorm:"not null"`
	UpdatedAt          time.Time `gorm:"not null"`
	CreatorID          uint      `gorm:"not null"`
	PlaceSquare        float64   `gorm:"not null"`
	OutsideTemperature float64   `gorm:"not null"`
	InsideTemperature  float64   `gorm:"not null"`
	CarrierVolume      float64   `gorm:"not null"`
	DeletedAt          *time.Time
	RequestHeaters     []RequestHeater `gorm:"foreignKey:HeatersProductRequestID;references:ID"`
}
