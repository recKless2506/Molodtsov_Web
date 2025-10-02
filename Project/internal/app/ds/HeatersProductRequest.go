package ds

import (
	"database/sql"

	"gorm.io/gorm"
)

type HeatersProductRequest struct {
	gorm.Model
	Status             string           `gorm:"column:status"`
	CreatorID          uint             `gorm:"column:creator_id"`
	FormationDate      sql.NullTime     `gorm:"column:formation_date"`
	DeleteDate         sql.NullTime     `gorm:"column:delete_date"`
	CompletionDate     sql.NullTime     `gorm:"column:completion_date"`
	RejectionDate      sql.NullTime     `gorm:"column:rejection_date"`
	ModeratorID        uint             `gorm:"column:moderator_id"`
	PlaceSquare        float64          `gorm:"column:place_square"`
	OutsideTemperature float64          `gorm:"column:outside_temperature"`
	InsideTemperature  float64          `gorm:"column:inside_temperature"`
	CarrierVolume      float64          `gorm:"column:carrier_volume"`
	Products           []RequestHeaters `gorm:"foreignKey:HeatersProductRequestID"`
}

// Указываем точное имя таблицы
func (HeatersProductRequest) TableName() string {
	return "heaters_product_requests"
}
