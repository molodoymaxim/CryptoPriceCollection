CREATE TABLE IF NOT EXISTS watched_currencies (
    coin VARCHAR(10) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS currency_prices (
    id SERIAL PRIMARY KEY,
    coin VARCHAR(10) NOT NULL,
    price FLOAT NOT NULL,
    timestamp BIGINT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_currency_timestamp ON currency_prices(coin, timestamp);