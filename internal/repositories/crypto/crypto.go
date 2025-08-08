package crypto

import (
	"CryptoPriceCollection/internal/repositories/crypto/postgresql"
	"CryptoPriceCollection/internal/system/database"
)

type Crypto struct {
	Postgres postgresql.CryptoRepository
}

func New(
	db *database.DataBase,
) *Crypto {
	return &Crypto{
		Postgres: postgresql.New(db),
	}
}
