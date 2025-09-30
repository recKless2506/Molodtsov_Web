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

	ctx.HTML(http.StatusOK, "catalog.html", gin.H{
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

	if product.Description == "" {
		product.Description = "Описание отсутствует"
	}
	if product.Specs == "" {
		product.Specs = "Технические характеристики отсутствуют"
	}

	ctx.HTML(http.StatusOK, "heater.html", gin.H{
		"Product": product,
	})
}

type ZayavkaItem struct {
	Product  repository.Product
	Defaults repository.ZayavkaDefaults
}

func (h *Handler) GetZayavka(ctx *gin.Context) {
	products, err := h.Repository.GetProducts()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Ошибка получения товаров")
		return
	}

	defaults := h.Repository.GetZayavkaDefaults()

	var items []ZayavkaItem
	for _, p := range products {
		if strings.Contains(p.Title, "Электрический котёл") || strings.Contains(p.Title, "Газовый котёл") {
			d, ok := defaults[p.ID]
			if !ok {
				d = repository.ZayavkaDefaults{}
			}
			items = append(items, ZayavkaItem{
				Product:  p,
				Defaults: d,
			})
		}
	}

	var inputDefaults repository.ZayavkaDefaults
	for _, v := range defaults {
		inputDefaults = v
		break
	}

	ctx.HTML(http.StatusOK, "application.html", gin.H{
		"products":      items,
		"inputDefaults": inputDefaults,
	})
}
