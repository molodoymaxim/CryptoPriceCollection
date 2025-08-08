package crypto

import (
	"CryptoPriceCollection/internal/repositories"
	"CryptoPriceCollection/internal/types"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CryptoServiceInterface interface {
	AddCurrency(ctx context.Context, coin string) error                                        // Добавление валюты в список наблюдаемых валют
	RemoveCurrency(ctx context.Context, coin string) error                                     // Удаление валюты из списка наблюдаемых валю
	GetPrice(ctx context.Context, coin string, timestamp *int64) (*types.CurrencyPrice, error) // Получение цены валюты
	StartPriceFetcher(ctx context.Context)                                                     // Фоновое получение цен
}

type CryptoService struct {
	repo          repositories.Repositories
	client        *http.Client
	fetchInterval time.Duration
	batchInterval time.Duration
	apiBaseURL    string
	prices        chan types.CurrencyPrice
}

func NewCryptoService(repo repositories.Repositories, apiBaseURL string, fetchInterval, batchInterval time.Duration) *CryptoService {
	return &CryptoService{
		repo:          repo,
		client:        &http.Client{Timeout: 10 * time.Second},
		fetchInterval: fetchInterval,
		batchInterval: batchInterval,
		apiBaseURL:    apiBaseURL,
		prices:        make(chan types.CurrencyPrice, 1000),
	}
}

// StartPriceFetcher запускает фоновое получение цен
func (s *CryptoService) StartPriceFetcher(ctx context.Context) {
	go s.batchWriter(ctx)

	ticker := time.NewTicker(s.fetchInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.fetchAndStorePrices(ctx)
		}
	}
}

// batchWriter обрабатывает пакетные вставки в БД
func (s *CryptoService) batchWriter(ctx context.Context) {
	ticker := time.NewTicker(s.batchInterval)
	defer ticker.Stop()

	var batch []types.CurrencyPrice
	for {
		select {
		case <-ctx.Done():
			return
		case price := <-s.prices:
			batch = append(batch, price)
		case <-ticker.C:
			if len(batch) > 0 {
				if err := s.repo.Crypto.Postgres.StoreBatch(ctx, batch); err != nil {
					log.Printf("Error storing batch: %v", err)
				}
				batch = nil
			}
		}
	}
}

// fetchAndStorePrices извлекает и сохраняет цены для отслеживаемых валют
func (s *CryptoService) fetchAndStorePrices(ctx context.Context) {
	coins, err := s.repo.Crypto.Postgres.GetWatchedCurrencies(ctx)
	if err != nil {
		log.Printf("Error fetching watched currencies: %v", err)
		return
	}

	if len(coins) == 0 {
		return
	}

	prices, err := s.fetchPricesFromAPI(coins)
	if err != nil {
		log.Printf("Error fetching prices: %v", err)
		return
	}

	timestamp := time.Now().Unix()
	for coin, price := range prices {
		s.prices <- types.CurrencyPrice{
			Coin:      coin,
			Price:     price,
			Timestamp: timestamp,
		}
	}
}

// fetchPricesFromAPI выводит цены на несколько монет
func (s *CryptoService) fetchPricesFromAPI(coins []string) (map[string]float64, error) {
	url := fmt.Sprintf("%s/simple/price?ids=%s&vs_currencies=usd", s.apiBaseURL, strings.Join(coins, ","))
	resp, err := s.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("HTTP request error to CoinGecko: %w", err)
	}
	if resp.StatusCode == 429 {
		log.Printf("Rate limit exceeded, pause for 60 seconds")
		time.Sleep(60 * time.Second)
		return nil, fmt.Errorf("rate limit exceeded")
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading the API response: %w", err)
	}
	log.Printf("Response CoinGecko API: %s", string(body))

	var result map[string]map[string]interface{}
	if err := json.NewDecoder(bytes.NewReader(body)).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding error JSON: %w", err)
	}

	prices := make(map[string]float64)
	for _, coin := range coins {
		if priceData, ok := result[coin]; ok {
			if usd, ok := priceData["usd"]; ok {
				switch v := usd.(type) {
				case float64:
					prices[coin] = v
				case string:
					price, err := strconv.ParseFloat(v, 64)
					if err != nil {
						log.Printf("Error converting the price for %s: %v", coin, err)
						continue
					}
					prices[coin] = price
				default:
					log.Printf("Unsupported price type for %s: %T", coin, v)
					continue
				}
			} else {
				log.Printf("The USD price was not found for %s", coin)
			}
		} else {
			log.Printf("No data was found for %s", coin)
		}
	}
	return prices, nil
}

// AddCurrency добавляет валюту в список отслеживаемых валют
func (s *CryptoService) AddCurrency(ctx context.Context, coin string) error {
	err := s.repo.Crypto.Postgres.AddCurrency(ctx, coin)
	if err != nil {
		log.Printf("Error in the repository when adding currency %s: %v", coin, err)
		return fmt.Errorf("couldn't add currency: %w", err)
	}
	return nil
}

// RemoveCurrency удаляет валюту из списка отслеживаемых валют
func (s *CryptoService) RemoveCurrency(ctx context.Context, coin string) error {
	err := s.repo.Crypto.Postgres.RemoveCurrency(ctx, coin)
	if err != nil {
		log.Printf("Error in the repository when deleting currency %s: %v", coin, err)
		return fmt.Errorf("couldn't delete currency: %w", err)
	}
	return nil
}

// GetPrice извлекает цену монеты, либо самую последнюю, либо на определенную временную метку
func (s *CryptoService) GetPrice(ctx context.Context, coin string, timestamp *int64) (*types.CurrencyPrice, error) {
	if timestamp == nil {
		return s.repo.Crypto.Postgres.GetLatestPrice(ctx, coin)
	}
	return s.repo.Crypto.Postgres.GetPrice(ctx, coin, *timestamp)
}
