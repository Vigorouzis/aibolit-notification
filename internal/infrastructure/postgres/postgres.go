package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	Name     string `env:"DATABASE"`
	Host     string `env:"HOST"`
	Port     string `env:"PORT"`
}

func (cfg Config) URL() string {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		cfg.Username, cfg.Password, cfg.Host,
		cfg.Port, cfg.Name,
	)
	return dsn
}

func ConnectViaConfig(cfg Config) (*sql.DB, error) {
	db, err := sql.Open("pgx", cfg.URL())
	if err != nil {
		return nil, err
	}
	return db, nil
}
