package ds

type RequestHeaters struct {
	HeatersProductRequestID uint `gorm:"primaryKey;auto_increment:false"`
	HeatersProductID        uint `gorm:"primaryKey;auto_increment:false"`
	Area                    float64

	Heaters HeatersProduct `gorm:"foreignKey:HeatersProductID"`
}
