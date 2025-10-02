package ds

import "gorm.io/gorm"

type HeatersProduct struct {
	gorm.Model         // автоматически добавит ID, CreatedAt, UpdatedAt, DeletedAt
	Title       string `gorm:"column:title"`
	Type        string `gorm:"column:type"`
	Description string `gorm:"column:description"`
	Power       string `gorm:"column:power"`
	Efficiency  string `gorm:"column:efficiency"`
	Image       string `gorm:"column:image"`
	IsDelete    bool   `gorm:"column:is_delete"`
}

// указываем точное имя таблицы
func (HeatersProduct) TableName() string {
	return "heaters_products"
}
