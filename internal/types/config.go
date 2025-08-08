package types

// ConfigPostgres конфигурация подключения к БД Postgres
type ConfigPostgres struct {
	PostgresHost         string `mapstructure:"POSTGRES_HOST"`
	PostgresPort         int    `mapstructure:"POSTGRES_PORT"`
	PostgresUser         string `mapstructure:"POSTGRES_USER"`
	PostgresPassword     string `mapstructure:"POSTGRES_PASSWORD"`
	PostgresDBName       string `mapstructure:"POSTGRES_DB"`
	SSLMode              string `mapstructure:"POSTGRES_SSLMODE"`
	PostgresQueryTimeout int    `mapstructure:"POSTGRES_QUERY_TIMEOUT"`
}

// ConfigConnDB конфигурация подключения к БД
type ConfigConnDB struct {
	CfgDBMaxConn      int `mapstructure:"DB_MAX_CONN"`
	CfgDBConnIdleTime int `mapstructure:"DB_CONN_IDLE_TIME"`
	CfgDBConnLifeTime int `mapstructure:"DB_CONN_LIFE_TIME"`
}

// ConfigHTTPServer конфигурация HTTP сервера
type ConfigHTTPServer struct {
	Port            int `mapstructure:"HTTP_PORT"`
	ReadTimeout     int `mapstructure:"HTTP_READ_TIMEOUT"`     // в секундах
	WriteTimeout    int `mapstructure:"HTTP_WRITE_TIMEOUT"`    // в секундах
	IdleTimeout     int `mapstructure:"HTTP_IDLE_TIMEOUT"`     // в секундах
	ShutdownTimeout int `mapstructure:"HTTP_SHUTDOWN_TIMEOUT"` // в секундах
}

// ConfigAPIClient внешняя конфигурация клиента
type ConfigAPIClient struct {
	BaseURL      string `mapstructure:"API_BASE_URL"`
	Timeout      int    `mapstructure:"API_TIMEOUT"`
	MaxRetries   int    `mapstructure:"API_MAX_RETRIES"`
	RetryBackoff int    `mapstructure:"API_RETRY_BACKOFF"`
}

// ConfigTasks конфигурация интервалов и батчинга PostgreSQL
type ConfigTasks struct {
	FetchInterval int `mapstructure:"FETCH_INTERVAL"`
	BatchInterval int `mapstructure:"BATCH_INTERVAL"`
}

// ConfigApp конфигурация всего приложения
type ConfigApp struct {
	Postgres   ConfigPostgres   `mapstructure:"postgres"`
	ConnDB     ConfigConnDB     `mapstructure:"db"`
	HTTPServer ConfigHTTPServer `mapstructure:"http"`
	APIClient  ConfigAPIClient  `mapstructure:"api"`
	Tasks      ConfigTasks      `mapstructure:"tasks"`
}
