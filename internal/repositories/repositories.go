package repositories

import (
	"CryptoPriceCollection/internal/repositories/crypto"
	"CryptoPriceCollection/internal/system"
)

type Repositories struct {
	Crypto *crypto.Crypto
}

func New(
	sys *system.Systems,
) *Repositories {
	return &Repositories{
		Crypto: crypto.New(sys.DB),
	}
}
