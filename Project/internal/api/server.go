package api

import (
	"Project/internal/app/handler"
	"Project/internal/app/repository"
	"log"
	"net/http"

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

	// Главная страница редирект на /hello
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusTemporaryRedirect, "/hello")
	})

	// Основные маршруты
	r.GET("/hello", h.GetCatalog)       // главная страница с каталогом
	r.GET("/product/:id", h.GetProduct) // страница товара
	r.GET("/zayavka", h.GetZayavka)     // страница заявки с карточками

	r.Run(":8080")
	log.Println("Server down")
}
