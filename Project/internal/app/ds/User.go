package ds

type User struct {
	ID          uint   `gorm:"primaryKey"`
	Login       string `gorm:"type:varchar(255);unique"`
	Password    string `gorm:"type:varchar(255)"`
	IsModerator bool
}
