package bootstrap

import (
	"fmt"
	"main/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func InitSqlxDB(cfg config.Config) (*sqlx.DB, error) {
	return sqlx.Connect(cfg.DBDriver, formatConnect(cfg))
}

func formatConnect(cfg config.Config) string {
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)
}
