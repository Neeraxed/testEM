package main

import (
	"net/http"
	"os"
	"testEM/internal/config"
	"testEM/internal/delivery"
	"testEM/internal/repository/song"
	"testEM/internal/repository/verse"
	"testEM/internal/usecase"
	"testEM/pkg/middleware"
	"testEM/pkg/postgresql"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	godotenv.Load("./.env")
	logger, _ := zap.NewDevelopment()
	config := config.ReadConfig(logger)
	//TODO change config
	db := postgresql.NewConnection(config)

	songStorage := song.NewStorage(db)
	verseStorage := verse.NewStorage(db)
	defer db.Close()

	//TODO create client
	uc := usecase.NewUsecase(songStorage, verseStorage, logger, nil)
	app := delivery.NewHandler(logger, uc)

	onion := middleware.NewOnion(logger)

	onion.AppendMiddleware(
		onion.Timer,
		onion.LogRequestResponse)

	router := app.ApplyRoutes(onion)

	port := os.Getenv("PORT")
	err := http.ListenAndServe(port, router)
	logger.Fatal("Server died",
		zap.String("message", err.Error()),
	)
}
