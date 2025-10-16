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
	r.POST("/add-to-cart/:id", h.AddToCart)
	r.POST("/heater", h.AddHeaterProduct)
	r.PUT("/heater/:id", h.UpdateHeaterProduct)
	r.DELETE("/heater/:id", h.DeleteHeaterProduct)
	r.POST("/heater/:id/image", h.UploadHeaterImage)
	r.PUT("/heaters_application/:id", h.UpdateHeatersProductRequest)
	r.PUT("/heaters_application/submit/:id", h.SubmitHeatersProductRequest)
	r.PUT("/heaters_application/moderate/:id", h.ModerateHeatersProductRequest)
	r.DELETE("/heaters_application/product", h.RemoveProductFromRequest)
	r.PUT("/heaters_application/product", h.UpdateRequestHeater)
	r.POST("/register", h.RegisterUser)
	r.GET("/me", h.GetCurrentUser)
	r.PUT("/me", h.UpdateCurrentUser)
	r.POST("/login", h.LoginUser)
	r.POST("/logout", h.LogoutUser)

	if err := r.Run(":8001"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
