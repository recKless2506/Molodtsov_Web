package repository

import (
	"Project/internal/app/ds"
	"errors"
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

	// Автоматические миграции, чтобы гарантировать наличие нужных таблиц
	if err := db.AutoMigrate(
		&ds.HeaterProduct{},
		&ds.HeatersProductRequest{},
		&ds.RequestHeater{},
		&ds.User{},
	); err != nil {
		return nil, fmt.Errorf("не удалось выполнить миграции: %w", err)
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

func (r *Repository) GetAllRequests(userID uint, isModerator bool) ([]ds.HeatersProductRequest, error) {
	var requests []ds.HeatersProductRequest
	query := r.db.Preload("RequestHeaters.HeaterProduct").Where("status != ?", "удален")

	if !isModerator {
		// Если не модератор, показываем только заявки пользователя
		query = query.Where("creator_id = ?", userID)
	}

	if err := query.Find(&requests).Error; err != nil {
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
	// Считаем количество товаров в корзине:
	// количество строк в request_heaters для заявок со статусом "черновик"
	if err := r.db.Model(&ds.RequestHeater{}).
		Joins("JOIN heaters_product_requests hpr ON hpr.id = request_heaters.heaters_product_request_id").
		Where("hpr.status = ?", "черновик").
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

// GetCartItems возвращает товары в корзине (заявки со статусом "черновик")
func (r *Repository) GetCartItems() ([]ds.HeatersProductRequest, error) {
	var requests []ds.HeatersProductRequest
	if err := r.db.Preload("RequestHeaters.HeaterProduct").
		Where("status = ?", "черновик").
		Find(&requests).Error; err != nil {
		return nil, err
	}
	return requests, nil
}

// Новый метод поиска
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
	// Проверяем, что товар существует
	var product ds.HeaterProduct
	if err := r.db.First(&product, productID).Error; err != nil {
		return fmt.Errorf("товар с ID %d не найден: %w", productID, err)
	}

	// Ищем существующую заявку в статусе "черновик" (корзина) для пользователя (пока используем CreatorID = 1)
	var request ds.HeatersProductRequest
	err := r.db.Where("creator_id = ? AND status = ?", 1, "черновик").First(&request).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			//  Если черновика нет — создаём новую заявку
			request = ds.HeatersProductRequest{
				Status:             "черновик",
				CreatorID:          1,
				PlaceSquare:        0,
				OutsideTemperature: 0,
				InsideTemperature:  0,
				CarrierVolume:      0,
			}

			if err := r.db.Create(&request).Error; err != nil {
				log.Println("Ошибка создания заявки:", err)
				return fmt.Errorf("не удалось создать заявку: %w", err)
			}

			log.Println("Создана новая заявка ID:", request.ID)

			if request.ID == 0 {
				return fmt.Errorf("request.ID = 0 после создания заявки")
			}
		} else {
			return fmt.Errorf("ошибка при поиске черновика: %w", err)
		}
	} else {
		log.Println("Найдена существующая заявка (черновик) ID:", request.ID)
	}

	// Создаем связь с товаром в request_heaters
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
func (r *Repository) ClearCart(requestId int) error {
	return r.db.Exec("UPDATE heaters_product_requests SET status='удален' WHERE id = $1", requestId).Error
}

func (r *Repository) GetHeatersRequestId(userId int) (int, error) {
	var id int
	err := r.db.Model(&ds.HeatersProductRequest{}).
		Where("creator_id = ? AND status = 'черновик'", userId).Select("id").First(&id).Error
	if err != nil {
		return 0, err
	}
	return id, nil
}
func (r *Repository) CreateHeaterProduct(product *ds.HeaterProduct) error {
	return r.db.Create(product).Error
}
func (r *Repository) UpdateHeaterProduct(id uint, updated *ds.HeaterProduct) error {
	var product ds.HeaterProduct
	if err := r.db.First(&product, id).Error; err != nil {
		return err
	}
	return r.db.Model(&product).Updates(updated).Error
}

// Логическое удаление товара
func (r *Repository) DeleteHeaterProduct(id uint) error {
	return r.db.Model(&ds.HeaterProduct{}).
		Where("id = ?", id).
		Update("is_delete", true).Error
}
func (r *Repository) UpdateHeatersProductRequest(id uint, updated *ds.HeatersProductRequest) error {
	var request ds.HeatersProductRequest
	if err := r.db.First(&request, id).Error; err != nil {
		return err
	}
	return r.db.Model(&request).Updates(updated).Error
}
func (r *Repository) SubmitHeatersProductRequest(id uint) error {
	return r.db.Model(&ds.HeatersProductRequest{}).
		Where("id = ? AND status = ?", id, "черновик").
		Update("status", "создано").Error
}
func (r *Repository) SetRequestStatus(id uint, status string) error {
	return r.db.Model(&ds.HeatersProductRequest{}).
		Where("id = ?", id).
		Update("status", status).Error
}
func (r *Repository) RemoveProductFromRequest(requestID, productID uint) error {
	return r.db.Where("heaters_product_request_id = ? AND heaters_product_id = ?", requestID, productID).
		Delete(&ds.RequestHeater{}).Error
}
func (r *Repository) UpdateRequestHeater(requestID, productID uint, area float64) error {
	return r.db.Model(&ds.RequestHeater{}).
		Where("heaters_product_request_id = ? AND heaters_product_id = ?", requestID, productID).
		Update("area", area).Error
}
func (r *Repository) CalculateRequestCost(requestID uint) error {
	var request ds.HeatersProductRequest
	if err := r.db.Preload("RequestHeaters").First(&request, requestID).Error; err != nil {
		log.Println("Ошибка при загрузке заявки:", err)
		return err
	}

	if len(request.RequestHeaters) == 0 {
		log.Println("В заявке нет товаров, расчёт cost невозможен")
		return nil
	}

	totalArea := 0.0
	for _, rh := range request.RequestHeaters {
		totalArea += rh.Area
	}

	cost := 2160 * 0.2 * totalArea * 8.49 * (request.InsideTemperature - request.OutsideTemperature)
	log.Printf("RequestID: %d, TotalArea: %.2f, Inside: %.2f, Outside: %.2f, Cost: %.2f",
		request.ID, totalArea, request.InsideTemperature, request.OutsideTemperature, cost)

	// Обновляем cost
	if err := r.db.Model(&ds.HeatersProductRequest{}).
		Where("id = ?", request.ID).
		Update("cost", cost).Error; err != nil {
		log.Println("Ошибка при обновлении cost:", err)
		return err
	}

	return nil
}

func (r *Repository) CreateUser(login, password string, isModerator bool) error {
	user := ds.User{
		Login:       login,
		Password:    password, // позже можно хэшировать
		IsModerator: isModerator,
	}
	return r.db.Create(&user).Error
}
func (r *Repository) GetUserByID(userID uint) (*ds.User, error) {
	var user ds.User
	if err := r.db.First(&user, userID).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) UpdateUser(userID uint, updated *ds.User) error {
	return r.db.Model(&ds.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"login":    updated.Login,
			"password": updated.Password, // позже можно хэшировать
		}).Error
}
func (r *Repository) AuthenticateUser(login, password string) (*ds.User, error) {
	var user ds.User
	if err := r.db.Where("login = ? AND password = ?", login, password).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *Repository) GetHeatersProductRequestByID(id uint) (*ds.HeatersProductRequest, error) {
	var request ds.HeatersProductRequest
	if err := r.DB().First(&request, id).Error; err != nil {
		return nil, err
	}
	return &request, nil
}
