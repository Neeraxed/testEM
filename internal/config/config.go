package config

import (
	"go.uber.org/zap"
)

type Config struct {
	postgresDSN string
}

func ReadConfig(log *zap.Logger) string {
	//TODO use .env
	//dbuser := os.Getenv("DBUSER")
	//dbpass := os.Getenv("DBPASSWORD")

	//return &Config{
	//	connStr: "user=" + dbuser + " password=" + dbpass + " dbname=tester sslmode=disable",
	//}

	return ""
}
