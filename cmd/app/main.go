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

	_ "testEM/docs"
)

//	@title			TestEM API
//	@version		0.0.1
//	@description	Test app for EM from Alina Kuznetsova
//	@termsOfService	http://github.com/Neeraxed

// @contact.name	Alina Kuznetsova
// @contact.email	Neeraxed@gmail.com
func main() {
	logger, _ := zap.NewDevelopment()
	err := godotenv.Load("./.env")
	if err != nil {
		logger.Fatal("Failed to load .env file",
			zap.String("message", err.Error()),
			zap.Time("time", time.Now()),
		)
	}
	logger.Info("Loaded .env file",
		zap.Time("time", time.Now()),
	)

	conf := config.ReadConfig()
	db := postgresql.NewConnection(conf.PostgresDSN, logger)
	defer db.Close()
	logger.Info("Connected to db",
		zap.Time("time", time.Now()),
	)

	songStorage := repository.NewSongStorage(db, logger)
	verseStorage := repository.NewVerseStorage(db, logger)

	cl := http.Client{}
	externalApiClient := usecase.NewDetailClient(conf.ExternalUrl, &cl, logger)

	//mock client for testing
	//externalApiClient := &delivery.MockExternal{}
	//uc := usecase.NewUsecase(songStorage, verseStorage, logger, externalApiClient)
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

	logger.Info("Server started",
		zap.Time("time", time.Now()),
	)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	if err != nil {
		logger.Fatal("Server died",
			zap.String("message", err.Error()),
		)
	}

	logger.Info("Server stopped",
		zap.Time("time", time.Now()),
	)
}
