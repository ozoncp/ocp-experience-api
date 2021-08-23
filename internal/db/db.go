package db

import (
	"github.com/rs/zerolog/log"

	sql "github.com/jmoiron/sqlx"

	_ "github.com/jackc/pgx/v4/stdlib"
)

// Connect connects to dsn and returns DB instance.
// Panics on connection failure.
func Connect(dsn string) *sql.DB {
	db, err := sql.Connect("pgx", dsn)

	if err != nil {
		log.Panic().AnErr("error", err).Msg("failed to connect to db")
	}

	return db
}
