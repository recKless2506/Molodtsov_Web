package api

import (
	"Project/internal/app/handler"
	"Project/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
)

func StartServer(repo *repository.Repository) {
	h := handler.NewHandler(repo)

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")

	r.GET("/", h.GetCatalog)
	r.GET("/hello", h.GetCatalog)
	r.GET("/heater/:id", h.GetHeaterByID)
	r.GET("/zayavka", h.GetApplication) // отображение страницы с заявками
	// <-- страница с заявками

	if err := r.Run(":8001"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
