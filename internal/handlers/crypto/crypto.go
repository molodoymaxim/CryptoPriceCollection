package crypto

import (
	"CryptoPriceCollection/internal/services/crypto"
	"CryptoPriceCollection/internal/types"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
)

type CryptoHandler interface {
	AddCurrencyHandler(c *gin.Context)
	RemoveCurrencyHandler(c *gin.Context)
	GetPriceHandler(c *gin.Context)
}

type cryptoHandler struct {
	service crypto.CryptoServiceInterface
}

func New(service crypto.CryptoServiceInterface) CryptoHandler {
	return &cryptoHandler{service: service}
}

// AddCurrencyHandler godoc
// @Summary      Добавить валюту
// @Description  Добавляет криптовалюту в список отслеживаемых (watched_currencies).
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        body body types.AddCurrencyRequest true "Запрос на добавление валюты"
// @Success      200 {object} map[string]string "status: success"
// @Failure      400 {object} map[string]string "error: Invalid request body"
// @Failure      500 {object} map[string]string "error: Failed to add currency: <details>"
// @Router       /currency/add [post]
func (h *cryptoHandler) AddCurrencyHandler(c *gin.Context) {
	var req types.AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.AddCurrency(c.Request.Context(), req.Coin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// RemoveCurrencyHandler godoc
// @Summary      Удалить валюту
// @Description  Удаляет криптовалюту из списка отслеживаемых (watched_currencies), сохраняя исторические цены в currency_prices.
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        body body types.AddCurrencyRequest true "Запрос на удаление валюты"
// @Success      200 {object} map[string]string "status: success"
// @Failure      400 {object} map[string]string "error: Invalid request body"
// @Failure      500 {object} map[string]string "error: Failed to remove currency: <details>"
// @Router       /currency/remove [post]
func (h *cryptoHandler) RemoveCurrencyHandler(c *gin.Context) {
	var req types.AddCurrencyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if err := h.service.RemoveCurrency(c.Request.Context(), req.Coin); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to remove currency"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// GetPriceHandler godoc
// @Summary      Получить цену валюты
// @Description  Возвращает последнюю цену валюты (без timestamp) или ближайшую цену к указанному времени (с timestamp).
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Param        body body types.PriceRequest true "Запрос на получение цены"
// @Success      200 {object} types.CurrencyPrice "Успешное получение цены"
// @Failure      400 {object} map[string]string "error: Invalid request body"
// @Failure      404 {object} map[string]string "error: Price not found"
// @Failure      500 {object} map[string]string "error: Failed to fetch price: <details>"
// @Router       /currency/price [post]
func (h *cryptoHandler) GetPriceHandler(c *gin.Context) {
	var req types.PriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	price, err := h.service.GetPrice(c.Request.Context(), req.Coin, req.Timestamp)
	if err == pgx.ErrNoRows {
		log.Printf("Цена не найдена для %s с timestamp=%v", req.Coin, req.Timestamp)
		c.JSON(http.StatusNotFound, gin.H{"error": "Price not found"})
		return
	}
	if err != nil {
		log.Printf("Ошибка получения цены для %s с timestamp=%v: %v", req.Coin, req.Timestamp, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch price"})
		return
	}

	c.JSON(http.StatusOK, price)
}
