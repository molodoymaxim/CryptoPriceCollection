package services

import (
	"CryptoPriceCollection/internal/repositories"
	"CryptoPriceCollection/internal/services/crypto"
	"time"
)

type Service struct {
	CryptoService crypto.CryptoServiceInterface
}

func NewService(repo repositories.Repositories, apiBaseURL string, fetchInterval, batchInterval time.Duration) *Service {
	return &Service{
		CryptoService: crypto.NewCryptoService(repo, apiBaseURL, fetchInterval, batchInterval),
	}
}
