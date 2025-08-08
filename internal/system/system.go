package system

import (
	"CryptoPriceCollection/internal/system/database"
	"CryptoPriceCollection/internal/types"
	"fmt"
)

type Systems struct {
	DB *database.DataBase
}

func New(cfgPostgres *types.ConfigPostgres, cfgConn *types.ConfigConnDB) (*Systems, error) {
	db, err := database.New(cfgPostgres, cfgConn)
	if err != nil {
		return nil, fmt.Errorf("database: %v", err)
	}
	return &Systems{
		DB: db,
	}, nil
}
