package api

import (
	"Project/internal/api/cache"
	"Project/internal/app/handler"
	"Project/internal/app/middleware"
	"Project/internal/app/repository"
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "Project/docs" // обязательно для swag
)

type pingReq struct{}
type pingResp struct {
	Status string `json:"status"`
}

// @title Heater API
// @version 1.0
// @description API для управления обогревателями и заявками
// @host localhost:8001
// @BasePath /
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// Ping godoc
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @Summary      Show hello text
// @Description  very very friendly response
// @Tags         Tests
// @Produce      json
// @Success      200  {object}  pingResp
// @Router       /ping/{name} [get]
func StartServer(repo *repository.Repository) {
	rdb := cache.NewRedis("localhost:6379", "")

	// создаём Redis клиента
	h := handler.NewHandler(repo, rdb) // передаем его в handler

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/static", "./resources")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// ===================== Каталог =====================
	r.GET("/", h.GetCatalog)
	r.GET("/catalog_heaters", h.GetCatalog)
	r.GET("/heater/:id", h.GetHeaterByID)

	// ===================== Корзина =====================
	r.GET("/cart", h.GetCart)
	r.GET("/cart-icon", h.GetCartIcon)
	r.PUT("/cart/update", h.UpdateCart)
	r.POST("/clear-cart", h.ClearCart)
	r.POST("/add-to-cart/:id", h.AddToCart)

	// ===================== Товары =====================
	r.POST("/heater", h.AddHeaterProduct)
	r.PUT("/heater/:id", h.UpdateHeaterProduct)
	r.DELETE("/heater/:id", h.DeleteHeaterProduct)
	r.POST("/heater/:id/image", h.UploadHeaterImage)

	// ===================== Пользователи =====================
	r.POST("/register", h.RegisterUser)
	r.POST("/login", h.LoginUser)

	// ===================== Защищённые маршруты =====================
	auth := r.Group("/")
	auth.Use(middleware.JWTMiddleware(rdb))
	{
		auth.POST("/logout", h.LogoutUser)
	}

	auth.Use(middleware.JWTMiddleware(rdb))
	{
		auth.GET("/heaters_application", h.GetApplications)
		auth.PUT("/heaters_application/:id", h.UpdateHeatersProductRequest)
		auth.PUT("/heaters_application/submit/:id", h.SubmitHeatersProductRequest)
		auth.PUT("/heaters_application/moderate/:id", h.ModerateHeatersProductRequest)
		auth.PUT("/heaters_application/product", h.UpdateRequestHeater)
	}

	// ===================== Запуск сервера =====================
	if err := r.Run(":8001"); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
