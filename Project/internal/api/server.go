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

	// ===================== Каталог =====================
	r.GET("/", h.GetCatalog)
	r.GET("/catalog_heaters", h.GetCatalog) // поиск можно реализовать через query param

	r.GET("/heater/:id", h.GetHeaterByID)

	// ===================== Корзина =====================
	r.POST("/clear-cart", h.ClearCart)
	r.POST("/add-to-cart/:id", h.AddToCart)

	// ===================== Товары =====================
	r.POST("/heater", h.AddHeaterProduct)          // Raw JSON
	r.PUT("/heater/:id", h.UpdateHeaterProduct)    // Raw JSON
	r.DELETE("/heater/:id", h.DeleteHeaterProduct) // JSON
	r.POST("/heater/:id/image", h.UploadHeaterImage)

	// ===================== Заявки =====================
	r.GET("/heaters_application", h.GetApplications)
	r.PUT("/heaters_application/:id", h.UpdateHeatersProductRequest)            // Raw JSON
	r.PUT("/heaters_application/submit/:id", h.SubmitHeatersProductRequest)     // JSON
	r.PUT("/heaters_application/moderate/:id", h.ModerateHeatersProductRequest) // JSON

	r.DELETE("/heaters_application/product/:id", h.DeleteHeaterProduct) // JSON
	r.PUT("/heaters_application/product", h.UpdateRequestHeater)        // PUT с JSON {request_id, product_id, area}

	// ===================== Пользователи =====================
	r.POST("/register", h.RegisterUser) // Raw JSON
	r.POST("/login", h.LoginUser)       // Raw JSON

	// ===================== Запуск сервера =====================
	if err := r.Run(":8001"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
