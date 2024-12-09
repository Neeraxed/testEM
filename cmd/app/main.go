package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testEM/internal/config"
	"testEM/internal/delivery"
	"testEM/internal/repository"
	"testEM/internal/usecase"
	"testEM/pkg/middleware"
	"testEM/pkg/postgresql"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewDevelopment()

	err := godotenv.Load("./.env")
	if err != nil {
		logger.Debug("Failed to load .env file",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}

	conf := config.ReadConfig()
	db := postgresql.NewConnection(conf.PostgresDSN, logger)
	songStorage := repository.NewSongStorage(db, logger)
	verseStorage := repository.NewVerseStorage(db, logger)
	defer db.Close()

	cl := http.Client{}
	externalApiClient := usecase.NewDetailClient(conf.ExternalUrl, &cl, logger)

	//mock client for testing
	//m := delivery.MockExternal{}
	uc := usecase.NewUsecase(songStorage, verseStorage, logger, externalApiClient)
	app := delivery.NewHandler(logger, uc)
	onion := middleware.NewOnion(logger)
	onion.AppendMiddleware(
		onion.Timer,
		onion.LogRequestResponse)
	router := app.ApplyRoutes(onion)
	go func() {
		err = http.ListenAndServe(conf.Port, router)
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	if err != nil {
		logger.Fatal("Server died",
			zap.String("message", err.Error()),
		)
	}
}
