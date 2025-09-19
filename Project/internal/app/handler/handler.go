package handler

import (
	"Project/internal/app/repository"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Handler struct {
	Repository *repository.Repository
}

func NewHandler(r *repository.Repository) *Handler {
	return &Handler{
		Repository: r,
	}
}

func (h *Handler) GetCatalog(ctx *gin.Context) {
	var products []repository.Product
	var err error

	searchQuery := ctx.Query("query")
	if searchQuery == "" {
		products, err = h.Repository.GetProducts()
		if err != nil {
			logrus.Error(err)
		}
	} else {
		products, err = h.Repository.GetProductsByTitle(searchQuery)
		if err != nil {
			logrus.Error(err)
		}
	}

	ctx.HTML(http.StatusOK, "index.html", gin.H{
		"time":     time.Now().Format("15:04:05"),
		"products": products,
		"query":    searchQuery,
	})
}

func (h *Handler) GetProduct(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusBadRequest, "Некорректный ID")
		return
	}

	product, err := h.Repository.GetProduct(id)
	if err != nil {
		logrus.Error(err)
		ctx.String(http.StatusNotFound, "Товар не найден")
		return
	}

	// Добавляем описание по умолчанию, если пустое
	if product.Description == "" {
		product.Description = "Описание отсутствует"
	}
	if product.Specs == "" {
		product.Specs = "Технические характеристики отсутствуют"
	}

	ctx.HTML(http.StatusOK, "order.html", gin.H{
		"Product": product,
	})
}

// Новый метод для заявки
func (h *Handler) GetZayavka(ctx *gin.Context) {
	products, err := h.Repository.GetProducts()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка получения товаров")
		return
	}

	// фильтруем только электрический и газовый котлы
	var filtered []repository.Product
	for _, p := range products {
		if strings.Contains(p.Title, "Электрический котёл") || strings.Contains(p.Title, "Газовый котёл") {
			filtered = append(filtered, p)
		}
	}

	ctx.HTML(http.StatusOK, "zayavka.html", gin.H{
		"products": filtered,
	})
}
