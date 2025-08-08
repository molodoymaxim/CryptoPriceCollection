package postgresql

import (
	"CryptoPriceCollection/internal/types"
	"context"
	"database/sql"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"net/url"
	"runtime"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Postgreser interface {
	NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error                  // Создание конфигурации пула
	ConnectionPool(ctx context.Context) error                                                   // Подключаемся с помощью пула к Postgres
	GetSQL(sqlFunc func(db *sql.DB) error) error                                                // Выполнение функции от имени драйвера sql.DB
	Ping(ctx context.Context) error                                                             // Проверяем соединение
	Transact(ctxParent context.Context, txFunc func(context.Context, pgx.Tx) error) (err error) // Обработчик транзакций
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)          // Exec запрос
	Close()                                                                                     // Закрытие соединения
	Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error)                  // Query запрос
	QueryRow(ctxParent context.Context, sql string, arguments ...any) pgx.Row                   // QueryRow запрос
}

type postgres struct {
	conn         *pgxpool.Pool
	connStr      string
	poolConfig   *pgxpool.Config
	queryTimeout time.Duration
}

func New(cfg *types.ConfigPostgres) Postgreser {
	connStr := fmt.Sprintf("%s://%s:%s@%s:%d/%s?sslmode=disable&connect_timeout=%d",
		"postgres",
		url.QueryEscape(cfg.PostgresUser),
		url.QueryEscape(cfg.PostgresPassword),
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDBName,
		cfg.PostgresQueryTimeout)

	return &postgres{
		connStr:      connStr,
		queryTimeout: time.Duration(cfg.PostgresQueryTimeout) * time.Second,
	}
}

func (d *postgres) NewPoolConfig(maxConn int, connIdleTime, connLifeTime time.Duration) error {
	// Создание конфигурации пула
	poolConfig, err := pgxpool.ParseConfig(d.connStr)
	if err != nil {
		return err
	}

	// Проверка
	cpu := runtime.NumCPU()
	if maxConn > cpu {
		maxConn = cpu
	}

	poolConfig.MaxConns = int32(maxConn)
	poolConfig.MaxConnIdleTime = connIdleTime
	poolConfig.MaxConnLifetime = connLifeTime
	d.poolConfig = poolConfig
	return nil
}

func (d *postgres) ConnectionPool(ctx context.Context) error {
	conn, err := pgxpool.NewWithConfig(ctx, d.poolConfig)
	if err != nil {
		return err
	}
	d.conn = conn
	err = d.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *postgres) Ping(ctx context.Context) error {
	return d.conn.Ping(ctx)
}

func (d *postgres) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return d.conn.Exec(ctx, sql, arguments...)
}

func (d *postgres) Query(ctx context.Context, sql string, arguments ...any) (pgx.Rows, error) {
	return d.conn.Query(ctx, sql, arguments...)
}

func (d *postgres) QueryRow(ctxParent context.Context, sql string, arguments ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(ctxParent, d.queryTimeout)
	defer cancel()
	return d.conn.QueryRow(ctx, sql, arguments...)
}

func (d *postgres) Close() {
	d.conn.Close()
}

func (d *postgres) Transact(ctxParent context.Context, txFunc func(context.Context, pgx.Tx) error) (err error) {
	ctx, cancel := context.WithTimeout(ctxParent, d.queryTimeout)
	defer cancel()

	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return
	}
	defer func() {
		if p := recover(); p != nil {
			if errRoll := tx.Rollback(ctx); errRoll != nil {
				log.Printf("Failed rollback database TX: %v\n", errRoll)
			}
			log.Printf("Failed rollback database: %v\n", p)
		} else if err != nil {
			if errRoll := tx.Rollback(ctx); errRoll != nil {
				log.Printf("Failed rollback database TX: %v\n", errRoll)
			}
		} else {
			err = tx.Commit(ctx)
		}
	}()
	err = txFunc(ctx, tx)
	return err
}

func (d *postgres) GetSQL(sqlFunc func(db *sql.DB) error) error {
	return sqlFunc(stdlib.OpenDBFromPool(d.conn))
}
