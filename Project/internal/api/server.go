package api

import (
	"Project/internal/app/handler"
	"Project/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func StartServer() {
	log.Println("Starting server")

	repo, err := repository.NewRepository()
	if err != nil {
		logrus.Error("ошибка инициализации репозитория")
	}

	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	// Основные маршруты
	r.GET("/", h.GetCatalog)            // главная страница
	r.GET("/hello", h.GetCatalog)       // старый маршрут для совместимости
	r.GET("/product/:id", h.GetProduct) // страница товара

	// Новый маршрут заявки (с карточками котлов)
	r.GET("/zayavka", h.GetZayavka)

	r.Run()
	log.Println("Server down")
}
