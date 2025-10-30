package handler

import (
	"Project/internal/app/ds"
	"Project/internal/app/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{Repository: r}
}

// ===================== Каталог =====================

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

func (h *Handler) GetApplications(ctx *gin.Context) {
	requests, err := h.Repository.GetAllRequests()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при получении заявок"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Список заявок", "data": requests})
}

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

func (h *Handler) SubmitHeatersProductRequest(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "Неверный ID заявки"})
		return
	}

	if err := h.Repository.SubmitHeatersProductRequest(uint(id)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при формировании заявки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Заявка успешно сформирована"})
}

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

// ===================== UpdateRequestHeater (для server.go) =====================

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

	// Пересчитываем стоимость заявки
	if err := h.Repository.CalculateRequestCost(uint(input.RequestID)); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "Ошибка при расчёте стоимости заявки"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно обновлён в заявке и стоимость пересчитана"})
}

// ===================== Пользователи =====================

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

	ctx.Set("userID", user.ID)
	ctx.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Пользователь успешно вошёл",
		"data": gin.H{
			"id":           user.ID,
			"login":        user.Login,
			"is_moderator": user.IsModerator,
		},
	})
}
