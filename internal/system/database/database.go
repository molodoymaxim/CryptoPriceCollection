package database

import (
	"CryptoPriceCollection/internal/system/database/postgresql"
	"CryptoPriceCollection/internal/types"
	"context"
	"fmt"
	"time"
)

type DataBase struct {
	Psql postgresql.Postgreser
}

func New(cfgPostgres *types.ConfigPostgres, cfgConn *types.ConfigConnDB) (*DataBase, error) {
	psql := postgresql.New(cfgPostgres)
	err := psql.NewPoolConfig(
		cfgConn.CfgDBMaxConn,
		time.Duration(cfgConn.CfgDBConnIdleTime)*time.Second,
		time.Duration(cfgConn.CfgDBConnLifeTime)*time.Second,
	)
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	err = psql.ConnectionPool(context.Background())
	if err != nil {
		return nil, fmt.Errorf("postgres: %v", err)
	}

	return &DataBase{
		Psql: psql,
	}, nil
}
