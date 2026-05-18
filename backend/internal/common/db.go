package common

import (
	"fmt"

	"github.com/seki18/lingdu-feed/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// DB is the global database connection pool.
var DB *sqlx.DB

// Init initializes the PostgreSQL connection using the provided config.
func Init(cfg config.Config) {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost,
		cfg.DBPort,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		panic(err)
	}

	DB = db
}
