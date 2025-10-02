package handler

import (
	"Project/internal/app/repository"
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

// Список товаров
func (h *Handler) GetCatalog(ctx *gin.Context) {
	products, err := h.Repository.GetHeaterProducts()
	if err != nil {
		log.Println("Ошибка получения продуктов:", err)
		ctx.String(http.StatusInternalServerError, "Ошибка получения товаров")
		return
	}

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
		"products": products,
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

	ctx.HTML(http.StatusOK, "heater.html", gin.H{
		"Product": product,
	})
}

// Страница с заявками (корзина)
func (h *Handler) GetApplication(ctx *gin.Context) {
	requests, err := h.Repository.GetAllRequests()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка получения заявок")
		return
	}

	ctx.HTML(http.StatusOK, "application.html", gin.H{
		"requests": requests,
	})
}
