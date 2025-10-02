package repository

import (
	"Project/internal/app/ds"
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(dsn string) (*Repository, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("не удалось подключиться к БД: %w", err)
	}
	return &Repository{db: db}, nil
}

func (r *Repository) DB() *gorm.DB {
	return r.db
}

func (r *Repository) GetHeaterProducts() ([]ds.HeatersProduct, error) {
	var products []ds.HeatersProduct
	if err := r.db.Where("is_delete = ?", false).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) GetHeaterProductByID(id uint) (*ds.HeatersProduct, error) {
	var product ds.HeatersProduct
	if err := r.db.Where("id = ? AND is_delete = ?", id, false).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

// Получаем все заявки и связанные товары
func (r *Repository) GetAllRequests() ([]ds.HeatersProductRequest, error) {
	var requests []ds.HeatersProductRequest
	if err := r.db.Preload("Products.Product").Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}
