package postgresql

import (
	"database/sql"
	migrate "github.com/rubenv/sql-migrate"
	"go.uber.org/zap"
)

func NewConnection(dsn string, log *zap.Logger) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}

	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)

	if err != nil {
		log.Fatal("Failed to run migrations", zap.Error(err))
	}

	return db
}
