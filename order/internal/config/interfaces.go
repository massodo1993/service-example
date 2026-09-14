package config

import "time"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type OrderHTTPConfig interface {
	Address() string
	ReadTimeout() time.Duration
}

type PostgresConfig interface {
	URI() string
	MigrationsDir() string
}

type InventoryGRPCConfig interface {
	Address() string
}

type PaymentGRPCConfig interface {
	Address() string
}
