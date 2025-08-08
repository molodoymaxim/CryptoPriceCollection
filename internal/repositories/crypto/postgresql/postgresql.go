package postgresql

import (
	"CryptoPriceCollection/internal/system/database"
	"CryptoPriceCollection/internal/types"
	"context"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"log"
)

type CryptoRepository interface {
	AddCurrency(ctx context.Context, coin string) error                                       // Добавление валюты в список наблюдаемых валют
	RemoveCurrency(ctx context.Context, coin string) error                                    // Удаление валюты из списка наблюдаемых валю
	GetPrice(ctx context.Context, coin string, timestamp int64) (*types.CurrencyPrice, error) // Получение цены валюты с указанием времени (если такой нет, то возьмется ближайшее время к заданному)
	GetLatestPrice(ctx context.Context, coin string) (*types.CurrencyPrice, error)            // Получение последней цены валюты
	GetWatchedCurrencies(ctx context.Context) ([]string, error)                               // Получение всех валют, которые наблюдаются
	StoreBatch(ctx context.Context, batch []types.CurrencyPrice) error                        // Пакетная вставка цен
}

type cryptoRepository struct {
	db *database.DataBase
}

func New(db *database.DataBase) CryptoRepository {
	return &cryptoRepository{
		db: db,
	}
}

// AddCurrency добавление валюты в список наблюдаемых валют
func (r *cryptoRepository) AddCurrency(ctx context.Context, coin string) error {
	query := "INSERT INTO watched_currencies (coin) VALUES ($1) ON CONFLICT DO NOTHING"
	tag, err := r.db.Psql.Exec(ctx, query, coin)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		log.Printf("The %s currency already exists or has not been added", coin)
	}

	return nil
}

// RemoveCurrency удаление валюты из списка наблюдаемых валют
func (r *cryptoRepository) RemoveCurrency(ctx context.Context, coin string) error {
	query := "DELETE FROM watched_currencies WHERE coin = $1"
	tag, err := r.db.Psql.Exec(ctx, query, coin)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		log.Printf("Currency %s not found for deletion", coin)
	}

	return nil
}

// GetPrice получение цены валюты с указанием времени (если такой нет, то возьмется ближайшее время к заданному)
func (r *cryptoRepository) GetPrice(ctx context.Context, coin string, timestamp int64) (*types.CurrencyPrice, error) {
	query := `SELECT coin, price, timestamp
			  FROM currency_prices
			  WHERE coin = $1
			  ORDER BY ABS(timestamp - $2)
			  LIMIT 1`
	rows, err := r.db.Psql.Query(ctx, query, coin, timestamp)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	currencyPrice := &types.CurrencyPrice{}
	err = pgxscan.ScanOne(currencyPrice, rows)
	if err != nil {
		log.Printf("Error scanning the price for %s with timestamp=%d: %v", coin, timestamp, err)
		return nil, err
	}
	return currencyPrice, nil
}

// GetLatestPrice получение последней цены валюты
func (r *cryptoRepository) GetLatestPrice(ctx context.Context, coin string) (*types.CurrencyPrice, error) {
	query := `SELECT coin, price, timestamp
			  FROM currency_prices
			  WHERE coin = $1
			  ORDER BY timestamp DESC 
			  LIMIT 1`
	rows, err := r.db.Psql.Query(ctx, query, coin)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	leatestPrice := &types.CurrencyPrice{}
	err = pgxscan.ScanOne(leatestPrice, rows)
	if err != nil {
		log.Printf("Error scanning the last price for %s: %v", coin, err)
		return nil, err
	}
	return leatestPrice, nil
}

// GetWatchedCurrencies получение всех валют, которые наблюдаются
func (r *cryptoRepository) GetWatchedCurrencies(ctx context.Context) ([]string, error) {
	query := "SELECT coin FROM watched_currencies"
	rows, err := r.db.Psql.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var coins []string
	for rows.Next() {
		var coin string
		if err := rows.Scan(&coin); err != nil {
			log.Printf("Error scanning coin: %v", err)
			continue
		}
		coins = append(coins, coin)
	}
	return coins, nil
}

// StoreBatch вставка пакеты с ценами
func (r *cryptoRepository) StoreBatch(ctx context.Context, batch []types.CurrencyPrice) error {
	return r.db.Psql.Transact(ctx, func(ctx context.Context, tx pgx.Tx) error {
		for _, price := range batch {
			_, err := tx.Exec(ctx, "INSERT INTO currency_prices (coin, price, timestamp) VALUES ($1, $2, $3)",
				price.Coin, price.Price, price.Timestamp)
			if err != nil {
				return err
			}
		}
		return nil
	})
}
