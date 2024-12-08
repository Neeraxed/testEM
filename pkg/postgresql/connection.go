package postgresql

import (
	"database/sql"
	migrate "github.com/rubenv/sql-migrate"
)

func NewConnection(dsn string) *sql.DB {
	db, err := sql.Open("postgres", dsn)
	if err == nil {
	}

	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	_, err = migrate.Exec(db, "postgres", migrations, migrate.Up)

	if err == nil {
	}
	
	return db
}

func Close(db *sql.DB) error {
	return db.Close()
}
