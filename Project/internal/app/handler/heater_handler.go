package handler

import (
	"Project/internal/api/cache"
	"Project/internal/app/ds"
	"Project/internal/app/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// Handler struct
type Handler struct {
	Repository *repository.Repository
	Redis      *cache.RedisClient
}

func NewHandler(r *repository.Repository, redis *cache.RedisClient) *Handler {
	return &Handler{Repository: r, Redis: redis}
}

// ===================== Каталог =====================

// GetCatalog godoc
// @Summary Get catalog of heaters
// @Description Get list of heaters with cart count
// @Tags catalog
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router / [get]
func (h *Handler) GetCatalog(ctx *gin.Context) {
	products, err := h.Repository.GetHeaterProducts()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка получения товаров"})
		return
	}

	count, err := h.Repository.GetRequestsCount()
	if err != nil {
		count = 0
	}

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Каталог товаров",
		"data": gin.H{
			"cart_count": count,
			"products":   products,
		},
	})
}

// GetHeaterByID godoc
// @Summary Get heater by ID
// @Param id path int true "Heater ID"
// @Produce json
// @Success 200 {object} ds.HeaterProduct
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /heater/{id} [get]
func (h *Handler) GetHeaterByID(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID товара"})
		return
	}

	product, err := h.Repository.GetHeaterProductByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Товар не найден"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар найден", "data": product})
}

// ===================== Корзина =====================

// ClearCart godoc
// @Summary Clear cart
// @Produce json
// @Success 303 {string} string "Redirect"
// @Failure 500 {object} map[string]interface{}
// @Router /clear-cart [post]
func (h *Handler) ClearCart(ctx *gin.Context) {
	if err := h.Repository.DB().
		Model(&ds.HeatersProductRequest{}).
		Where("status = ?", "черновик").
		Update("status", "удален").Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при очистке корзины"})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/")
}

// AddToCart godoc
// @Summary Add heater to cart
// @Param id path int true "Heater ID"
// @Produce json
// @Success 303 {string} string "Redirect"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /add-to-cart/{id} [post]
func (h *Handler) AddToCart(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID товара"})
		return
	}

	if err := h.Repository.AddProductToCart(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при добавлении товара"})
		return
	}

	ctx.Redirect(http.StatusSeeOther, "/catalog_heaters")
}

// ===================== CRUD Товары =====================

// AddHeaterProduct godoc
// @Summary Add new heater product
// @Accept json
// @Produce json
// @Param product body ds.HeaterProduct true "Heater Product"
// @Success 200 {object} ds.HeaterProduct
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heater [post]
func (h *Handler) AddHeaterProduct(ctx *gin.Context) {
	var input ds.HeaterProduct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверные данные"})
		return
	}

	if err := h.Repository.CreateHeaterProduct(&input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при добавлении товара"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно добавлен", "data": input})
}

// UpdateHeaterProduct godoc
// @Summary Update heater product
// @Accept json
// @Produce json
// @Param id path int true "Heater ID"
// @Param product body ds.HeaterProduct true "Heater Product"
// @Success 200 {object} ds.HeaterProduct
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heater/{id} [put]
func (h *Handler) UpdateHeaterProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID товара"})
		return
	}

	var input ds.HeaterProduct
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверные данные"})
		return
	}

	if err := h.Repository.UpdateHeaterProduct(uint(id), &input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при обновлении товара"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно обновлён", "data": input})
}

// DeleteHeaterProduct godoc
// @Summary Delete heater product (soft)
// @Param id path int true "Heater ID"
// @Produce json
// @Success 200 {string} string "Deleted"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heater/{id} [delete]
func (h *Handler) DeleteHeaterProduct(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID товара"})
		return
	}

	if err := h.Repository.DeleteHeaterProduct(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при удалении товара"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно удалён"})
}

// UploadHeaterImage godoc
// @Summary Upload image for heater product
// @Param id path int true "Heater ID"
// @Param image formData file true "Image file"
// @Produce json
// @Success 200 {string} string "Uploaded"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heater/{id}/image [post]
func (h *Handler) UploadHeaterImage(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID товара"})
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Ошибка при получении файла"})
		return
	}

	dst := "./resources/images/" + file.Filename
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при сохранении файла"})
		return
	}

	if err := h.Repository.UpdateHeaterProduct(uint(id), &ds.HeaterProduct{Image: "/static/images/" + file.Filename}); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при обновлении карточки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Изображение успешно загружено"})
}

// ===================== Заявки =====================

// @Summary Получить список заявок
// @Description Только для авторизованных пользователей
// @Tags Applications
// @Produce json
// @Success 200 {object} []ds.HeatersProductRequest
// @Failure 401 {object} map[string]interface{}
// @Security BearerAuth
// @Router /heaters_application [get]

func (h *Handler) GetApplications(ctx *gin.Context) {
	// Получаем userID и isModerator из контекста
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
		return
	}
	userID := userIDVal.(uint)

	isModeratorVal, exists := ctx.Get("isModerator")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
		return
	}
	isModerator := isModeratorVal.(bool)

	// Получаем заявки с учетом роли пользователя
	requests, err := h.Repository.GetAllRequests(userID, isModerator)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при получении заявок"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Список заявок", "data": requests})
}

// UpdateHeatersProductRequest godoc
// @Summary Update heater request
// @Accept json
// @Produce json
// @Param id path int true "Request ID"
// @Param request body ds.HeatersProductRequest true "Request data"
// @Success 200 {object} ds.HeatersProductRequest
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heaters_application/{id} [put]
func (h *Handler) UpdateHeatersProductRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID заявки"})
		return
	}

	var input ds.HeatersProductRequest
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверные данные"})
		return
	}

	if err := h.Repository.UpdateHeatersProductRequest(uint(id), &input); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при обновлении заявки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Заявка успешно обновлена", "data": input})
}

// SubmitHeatersProductRequest godoc
// @Summary Submit heater request
// @Produce json
// @Param id path int true "Request ID"
// @Success 200 {string} string "Submitted"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heaters_application/submit/{id} [put]
func (h *Handler) SubmitHeatersProductRequest(ctx *gin.Context) {
	// Получаем userID и isModerator из контекста
	userIDVal, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
		return
	}
	userID := userIDVal.(uint)

	isModeratorVal, exists := ctx.Get("isModerator")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Не авторизован"})
		return
	}
	isModerator := isModeratorVal.(bool)

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID заявки"})
		return
	}

	// Получаем заявку из БД
	request, err := h.Repository.GetHeatersProductRequestByID(uint(id))
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"success": false, "message": "Заявка не найдена"})
		return
	}

	// Проверяем, является ли пользователь создателем и не модератором
	if request.UserID == userID && !isModerator {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Только модератор может завершить заявку"})
		return
	}

	// Только модератор может завершить заявку
	if !isModerator {
		ctx.JSON(http.StatusForbidden, gin.H{"success": false, "message": "Нет прав для завершения заявки"})
		return
	}

	// Завершаем заявку
	if err := h.Repository.SubmitHeatersProductRequest(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при формировании заявки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Заявка успешно завершена"})
}

// ModerateHeatersProductRequest godoc
// @Summary Moderate heater request (change status)
// @Produce json
// @Param id path int true "Request ID"
// @Param status query string true "New status: завершено или отклонено"
// @Success 200 {string} string "Status updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heaters_application/moderate/{id} [put]
func (h *Handler) ModerateHeatersProductRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID заявки"})
		return
	}

	status := ctx.Query("status")
	if status != "завершено" && status != "отклонено" {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный статус"})
		return
	}

	if err := h.Repository.SetRequestStatus(uint(id), status); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при изменении статуса заявки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Статус заявки успешно обновлён"})
}

// UpdateRequestHeater godoc
// @Summary Update heater in request (area, cost will be recalculated)
// @Accept json
// @Produce json
// @Param request body object{request_id=int, product_id=int, area=number} true "Request-Product update"
// @Success 200 {string} string "Updated"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /heaters_application/product [put]
func (h *Handler) UpdateRequestHeater(ctx *gin.Context) {
	var input struct {
		RequestID int     `json:"request_id"`
		ProductID int     `json:"product_id"`
		Area      float64 `json:"area"`
	}

	if err := ctx.ShouldBindJSON(&input); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверные данные"})
		return
	}

	if err := h.Repository.UpdateRequestHeater(uint(input.RequestID), uint(input.ProductID), input.Area); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при обновлении товара в заявке"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно обновлён в заявке"})
}

// ===================== Пользователи =====================

// RegisterUser godoc
// @Summary Register new user
// @Accept json
// @Produce json
// @Param user body object{login=string,password=string} true "User credentials"
// @Success 200 {string} string "User registered"
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func (h *Handler) RegisterUser(ctx *gin.Context) {
	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil || input.Login == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Логин и пароль обязательны"})
		return
	}

	if err := h.Repository.CreateUser(input.Login, input.Password, false); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при создании пользователя"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Пользователь успешно зарегистрирован"})
}

// LoginUser godoc
// @Summary Вход пользователя
// @Description Авторизация по логину и паролю, возвращает JWT или session
// @Tags Users
// @Accept json
// @Produce json
// @Param login body object{login=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
var jwtSecret = []byte("secret_key") // секретный ключ для подписи JWT, можно хранить в env

// LoginUser godoc
// @Summary Вход пользователя
// @Description Авторизация по логину и паролю, возвращает JWT или session
// @Tags Users
// @Accept json
// @Produce json
// @Param login body object{login=string,password=string} true "Login credentials"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /login [post]
func (h *Handler) LoginUser(ctx *gin.Context) {
	var input struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := ctx.ShouldBindJSON(&input); err != nil || input.Login == "" || input.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Логин и пароль обязательны"})
		return
	}

	user, err := h.Repository.AuthenticateUser(input.Login, input.Password)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Неверный логин или пароль"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": user.ID,
		"exp":    time.Now().Add(time.Hour * 72).Unix(),
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка генерации токена"})
		return
	}
	ctx.SetCookie("session_token", tokenString, 3600*24*3, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Пользователь успешно вошёл",
		"data": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
			"token":        tokenString,
		},
	})
}

// LogoutUser godoc
// @Summary Logout user
// @Description Добавляет токен пользователя в blacklist
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Router /logout [post]
func (h *Handler) LogoutUser(c *gin.Context) {
	// Получаем токен из контекста (устанавливается в JWTMiddleware)
	tokenRaw, exists := c.Get("userToken")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "Токен не найден"})
		return
	}
	tokenString, ok := tokenRaw.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка обработки токена"})
		return
	}

	// Заносим токен в blacklist на срок его жизни (например, 24 часа)
	if err := h.Redis.BlacklistToken(tokenString, 24*time.Hour); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка занесения токена в черный список"})
		return
	}

	// Можно удалить cookie, если используешь сессии в cookie
	c.SetCookie("session_token", "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Вы успешно вышли",
	})
}
