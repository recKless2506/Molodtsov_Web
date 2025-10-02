package ds

import "gorm.io/gorm"

type RequestHeaters struct {
	gorm.Model
	HeatersProductRequestID uint
	HeatersProductID        uint
	Area                    float64

	Product HeatersProduct `gorm:"foreignKey:HeatersProductID"`
}
