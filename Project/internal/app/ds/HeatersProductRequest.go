package ds

import (
	"time"
)

type HeatersProductRequest struct {
	ID             int       `gorm:"primaryKey"`
	Status         string    `gorm:"not null"`
	CreatedAt      time.Time `gorm:"not null"`
	CreatorId      int       `gorm:"not null"`
	FormationDate  time.Time
	DeleteDate     time.Time `gorm:"default:NULL"`
	CompletionDate time.Time `gorm:"default:NULL"`
	RejectionDate  time.Time `gorm:"default:NULL"`
	ModeratorID    int       `gorm:"default:NULL"`
	TotalPower     float64
	Insolation     float64

	Heaters []RequestHeaters `gorm:"foreignKey:HeatersProductRequestID"`
}
