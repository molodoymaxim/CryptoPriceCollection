package handlers

import (
	"CryptoPriceCollection/internal/handlers/crypto"
	"CryptoPriceCollection/internal/services"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	crypto crypto.CryptoHandler
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{
		crypto: crypto.New(services.CryptoService),
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()

	router.POST("/currency/add", h.crypto.AddCurrencyHandler)
	router.POST("/currency/remove", h.crypto.RemoveCurrencyHandler)
	router.POST("/currency/price", h.crypto.GetPriceHandler)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	return router
}
