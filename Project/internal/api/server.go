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
	r.GET("/catalog_heaters", h.SearchCatalog)
	r.GET("/heater/:id", h.GetHeaterByID)

	r.GET("/heaters_application", h.GetApplications)
	r.POST("/clear-cart", h.ClearCart)
	r.POST("/add-to-cart/:id", h.AddToCart) // новый POST-роу

	if err := r.Run(":8001"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
