package handler

import (
	"Project/internal/app/ds"
	"Project/internal/app/repository"
	"fmt"
	"log"
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

// Вспомогательная функция: определить, нужен ли JSON
func wantsJSON(ctx *gin.Context) bool {
	return ctx.Query("format") == "json" || ctx.GetHeader("Accept") == "application/json"
}

// ===================== Основные обработчики =====================

// Список товаров (каталог)
func (h *Handler) GetCatalog(ctx *gin.Context) {
	products, err := h.Repository.GetHeaterProducts()
	if err != nil {
		log.Println("Ошибка получения продуктов:", err)
		ctx.String(http.StatusInternalServerError, "Ошибка получения товаров")
		return
	}

	count, err := h.Repository.GetRequestsCount()
	if err != nil {
		log.Println("Ошибка получения количества заявок:", err)
		count = 0
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Каталог товаров",
			"data": gin.H{
				"cart_count": count,
				"products":   products,
			},
		})
		return
	}

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"products":   products,
		"cart_count": count,
	})
}

// Конкретный товар
func (h *Handler) GetHeaterByID(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	product, err := h.Repository.GetHeaterProductByID(uint(id))
	if err != nil {
		ctx.String(http.StatusNotFound, "Товар не найден")
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Товар найден",
			"data":    product,
		})
		return
	}

	ctx.HTML(http.StatusOK, "heater.html", gin.H{
		"Product": product,
	})
}

// Страница с заявками
func (h *Handler) GetApplications(ctx *gin.Context) {
	requests, err := h.Repository.GetAllRequests()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при получении заявок: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Список заявок",
			"data":    requests,
		})
		return
	}

	ctx.HTML(http.StatusOK, "application.html", gin.H{
		"requests": requests,
	})
}

// Очистка корзины (POST-запрос)
func (h *Handler) ClearCart(ctx *gin.Context) {
	if err := h.Repository.DB().
		Model(&ds.HeatersProductRequest{}).
		Where("status = ?", "черновик").
		Update("status", "удален").Error; err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при очистке корзины: %v", err)
		return
	}

	// Редирект оставляем
	ctx.Redirect(http.StatusSeeOther, "/")
}

// Поиск товаров
func (h *Handler) SearchCatalog(ctx *gin.Context) {
	query := ctx.Query("query")
	if query == "" {
		ctx.Redirect(http.StatusSeeOther, "/")
		return
	}

	products, err := h.Repository.SearchHeaterProducts(query)
	if err != nil {
		log.Println("Ошибка поиска продуктов:", err)
		ctx.String(http.StatusInternalServerError, "Ошибка поиска товаров")
		return
	}

	count, err := h.Repository.GetRequestsCount()
	if err != nil {
		log.Println("Ошибка получения количества заявок:", err)
		count = 0
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": fmt.Sprintf("Результаты поиска по запросу '%s'", query),
			"data": gin.H{
				"cart_count": count,
				"products":   products,
			},
		})
		return
	}

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"products":   products,
		"cart_count": count,
	})
}

// Добавление товара в корзину
func (h *Handler) AddToCart(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	if err := h.Repository.AddProductToCart(uint(id)); err != nil {
		log.Println("Ошибка при добавлении товара в корзину:", err)
		ctx.String(http.StatusInternalServerError, "Ошибка при добавлении товара в корзину: "+err.Error())
		return
	}

	// Редирект оставляем
	ctx.Redirect(http.StatusSeeOther, "/catalog_heaters")
}

// ===================== CRUD для товаров =====================

func (h *Handler) AddHeaterProduct(ctx *gin.Context) {
	var input ds.HeaterProduct
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.String(http.StatusBadRequest, "Неверные данные: %v", err)
		return
	}

	if err := h.Repository.CreateHeaterProduct(&input); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при добавлении товара: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно добавлен", "data": input})
		return
	}

	ctx.String(http.StatusOK, "Товар успешно добавлен")
}

func (h *Handler) UpdateHeaterProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	var input ds.HeaterProduct
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.String(http.StatusBadRequest, "Неверные данные: %v", err)
		return
	}

	if err := h.Repository.UpdateHeaterProduct(uint(id), &input); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при обновлении товара: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно обновлён", "data": input})
		return
	}

	ctx.String(http.StatusOK, "Товар успешно обновлён")
}

func (h *Handler) DeleteHeaterProduct(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	if err := h.Repository.DeleteHeaterProduct(uint(id)); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при удалении товара: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно удалён"})
		return
	}

	ctx.String(http.StatusOK, "Товар успешно удалён")
}

// ===================== Работа с изображениями =====================

func (h *Handler) UploadHeaterImage(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	file, err := ctx.FormFile("image")
	if err != nil {
		ctx.String(http.StatusBadRequest, "Ошибка при получении файла: %v", err)
		return
	}

	dst := fmt.Sprintf("./resources/images/%s", file.Filename)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при сохранении файла: %v", err)
		return
	}

	if err := h.Repository.UpdateHeaterProduct(uint(id), &ds.HeaterProduct{Image: "/static/images/" + file.Filename}); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при обновлении карточки: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Изображение успешно загружено"})
		return
	}

	ctx.String(http.StatusOK, "Изображение успешно загружено")
}

// ===================== Работа с заявками =====================

func (h *Handler) UpdateHeatersProductRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	var input ds.HeatersProductRequest
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.String(http.StatusBadRequest, "Неверные данные: %v", err)
		return
	}

	if err := h.Repository.UpdateHeatersProductRequest(uint(id), &input); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при обновлении заявки: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Заявка успешно обновлена", "data": input})
		return
	}

	ctx.String(http.StatusOK, "Заявка успешно обновлена")
}

func (h *Handler) SubmitHeatersProductRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	if err := h.Repository.SubmitHeatersProductRequest(uint(id)); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при формировании заявки: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Заявка успешно сформирована"})
		return
	}

	ctx.String(http.StatusOK, "Заявка успешно сформирована")
}

func (h *Handler) ModerateHeatersProductRequest(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	status := ctx.Query("status")
	if status != "завершено" && status != "отклонено" {
		ctx.String(http.StatusBadRequest, "Неверный статус")
		return
	}

	if err := h.Repository.SetRequestStatus(uint(id), status); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при изменении статуса заявки: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Статус заявки успешно обновлён"})
		return
	}

	ctx.String(http.StatusOK, "Статус заявки успешно обновлён")
}

// ===================== Работа с пользователями =====================

func (h *Handler) RegisterUser(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := ctx.PostForm("password")

	if login == "" || password == "" {
		ctx.String(http.StatusBadRequest, "Логин и пароль обязательны")
		return
	}

	if err := h.Repository.CreateUser(login, password, false); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при создании пользователя: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Пользователь успешно зарегистрирован"})
		return
	}

	ctx.String(http.StatusOK, "Пользователь успешно зарегистрирован")
}

func (h *Handler) GetCurrentUser(ctx *gin.Context) {
	userIDInterface, exists := ctx.Get("userID")
	if !exists {
		ctx.String(http.StatusUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.String(http.StatusInternalServerError, "Ошибка обработки ID пользователя")
		return
	}

	user, err := h.Repository.GetUserByID(userID)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при получении данных пользователя: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Данные пользователя",
			"data": gin.H{
				"id":           user.ID,
				"login":        user.Login,
				"is_moderator": user.IsModerator,
			},
		})
		return
	}

	ctx.String(http.StatusOK, fmt.Sprintf("Пользователь: %s", user.Login))
}

func (h *Handler) UpdateCurrentUser(ctx *gin.Context) {
	userIDInterface, exists := ctx.Get("userID")
	if !exists {
		ctx.String(http.StatusUnauthorized, "Пользователь не аутентифицирован")
		return
	}

	userID, ok := userIDInterface.(uint)
	if !ok {
		ctx.String(http.StatusInternalServerError, "Ошибка обработки ID пользователя")
		return
	}

	var input ds.User
	if err := ctx.ShouldBind(&input); err != nil {
		ctx.String(http.StatusBadRequest, "Неверные данные: %v", err)
		return
	}

	if err := h.Repository.UpdateUser(userID, &input); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при обновлении данных пользователя: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Данные пользователя успешно обновлены"})
		return
	}

	ctx.String(http.StatusOK, "Данные пользователя успешно обновлены")
}

func (h *Handler) LoginUser(ctx *gin.Context) {
	login := ctx.PostForm("login")
	password := ctx.PostForm("password")

	if login == "" || password == "" {
		ctx.String(http.StatusBadRequest, "Логин и пароль обязательны")
		return
	}

	user, err := h.Repository.AuthenticateUser(login, password)
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Неверный логин или пароль")
		return
	}

	ctx.Set("userID", user.ID)

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Пользователь успешно вошёл",
			"data": gin.H{
				"id":           user.ID,
				"login":        user.Login,
				"is_moderator": user.IsModerator,
			},
		})
		return
	}

	ctx.String(http.StatusOK, "Пользователь успешно вошёл")
}

func (h *Handler) LogoutUser(ctx *gin.Context) {
	ctx.Set("userID", nil)

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Пользователь успешно вышел"})
		return
	}

	ctx.String(http.StatusOK, "Пользователь успешно вышел")
}

// Удаление товара из заявки
func (h *Handler) RemoveProductFromRequest(ctx *gin.Context) {
	requestParam := ctx.Query("request_id")
	productParam := ctx.Query("product_id")

	requestID, err := strconv.Atoi(requestParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	productID, err := strconv.Atoi(productParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	if err := h.Repository.RemoveProductFromRequest(uint(requestID), uint(productID)); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при удалении товара из заявки: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно удалён из заявки"})
		return
	}

	ctx.String(http.StatusOK, "Товар успешно удалён из заявки")
}

// Обновление товара в заявке (например, площадь)
func (h *Handler) UpdateRequestHeater(ctx *gin.Context) {
	requestParam := ctx.Query("request_id")
	productParam := ctx.Query("product_id")
	areaParam := ctx.Query("area")

	requestID, err := strconv.Atoi(requestParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID заявки")
		return
	}

	productID, err := strconv.Atoi(productParam)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверный ID товара")
		return
	}

	area, err := strconv.ParseFloat(areaParam, 64)
	if err != nil {
		ctx.String(http.StatusBadRequest, "Неверное значение площади")
		return
	}

	if err := h.Repository.UpdateRequestHeater(uint(requestID), uint(productID), area); err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка при обновлении товара в заявке: %v", err)
		return
	}

	if wantsJSON(ctx) {
		ctx.JSON(http.StatusOK, gin.H{"success": true, "message": "Товар успешно обновлён в заявке"})
		return
	}

	ctx.String(http.StatusOK, "Товар успешно обновлён в заявке")
}
