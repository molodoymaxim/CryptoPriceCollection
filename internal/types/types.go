package types

// CurrencyPrice содержит информацию по монете
type CurrencyPrice struct {
	Coin      string  `json:"coin"`
	Price     float64 `json:"price"`
	Timestamp int64   `json:"timestamp"`
}

// AddCurrencyRequest содержит список использующзихся монет
type AddCurrencyRequest struct {
	Coin string `json:"coin" binding:"required"`
}

// PriceRequest запрос на получение цены
type PriceRequest struct {
	Coin      string `json:"coin" binding:"required"`
	Timestamp *int64 `json:"timestamp"`
}
