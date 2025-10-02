package ds

import "time"

type RequestHeater struct {
	HeatersProductRequestID uint `gorm:"column:heaters_product_request_id;primaryKey"`
	HeatersProductID        uint `gorm:"column:heaters_product_id;primaryKey"`
	Area                    float64
	DeletedAt               *time.Time
	HeaterProduct           HeaterProduct `gorm:"foreignKey:HeatersProductID;references:ID"`
}

func (RequestHeater) TableName() string {
	return "request_heaters"
}
