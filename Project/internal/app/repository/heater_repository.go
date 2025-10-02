package repository

import (
	"Project/internal/app/ds"
	"fmt"
	"log"

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

func (r *Repository) GetHeaterProducts() ([]ds.HeaterProduct, error) {
	var products []ds.HeaterProduct
	if err := r.db.Where("is_delete = ?", false).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) GetHeaterProductByID(id uint) (*ds.HeaterProduct, error) {
	var product ds.HeaterProduct
	if err := r.db.Where("id = ? AND is_delete = ?", id, false).First(&product).Error; err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) GetAllRequests() ([]ds.HeatersProductRequest, error) {
	var requests []ds.HeatersProductRequest

	err := r.db.
		Preload("RequestHeaters.HeaterProduct"). // загружаем товары внутри заявки
		Where("status != ?", "удален").
		Find(&requests).Error

	if err != nil {
		return nil, err
	}
	return requests, nil
}

func (r *Repository) ClearRequests() error {
	return r.db.Model(&ds.HeatersProductRequest{}).
		Where("status = ?", "черновик").
		Update("status", "удален").Error
}

func (r *Repository) GetRequestsCount() (int64, error) {
	var count int64
	if err := r.db.Model(&ds.HeatersProductRequest{}).
		Where("status != ?", "удален").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// 🔥 Новый метод поиска
func (r *Repository) SearchHeaterProducts(query string) ([]ds.HeaterProduct, error) {
	var products []ds.HeaterProduct
	if err := r.db.Where(
		"is_delete = ? AND (title ILIKE ? OR type ILIKE ? OR description ILIKE ?)",
		false,
		"%"+query+"%",
		"%"+query+"%",
		"%"+query+"%",
	).Find(&products).Error; err != nil {
		return nil, err
	}
	return products, nil
}

func (r *Repository) AddProductToCart(productID uint) error {
	// 1️⃣ Проверяем, что товар существует
	var product ds.HeaterProduct
	if err := r.db.First(&product, productID).Error; err != nil {
		return fmt.Errorf("товар с ID %d не найден: %w", productID, err)
	}

	// 2️⃣ Создаем заявку с статусом "черновик"
	request := ds.HeatersProductRequest{
		Status:             "черновик",
		CreatorID:          1, // можно заменить на текущего пользователя
		PlaceSquare:        0,
		OutsideTemperature: 0,
		InsideTemperature:  0,
		CarrierVolume:      0,
	}

	if err := r.db.Create(&request).Error; err != nil {
		log.Println("Ошибка создания заявки:", err)
		return fmt.Errorf("не удалось создать заявку: %w", err)
	}

	log.Println("Создана заявка ID:", request.ID)

	if request.ID == 0 {
		return fmt.Errorf("request.ID = 0 после создания заявки")
	}

	// 3️⃣ Создаем связь с товаром в request_heaters
	link := ds.RequestHeater{
		HeatersProductRequestID: request.ID,
		HeatersProductID:        productID,
		Area:                    0,
	}

	if err := r.db.Create(&link).Error; err != nil {
		log.Println("Ошибка создания связи request_heaters:", err)
		return fmt.Errorf("не удалось добавить товар к заявке: %w", err)
	}

	log.Println("Товар успешно добавлен в корзину:", product.Title)

	return nil
}
