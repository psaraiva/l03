package main

import (
	"context"
	"errors"
	"l03/configuration/database/mongodb"
	"l03/configuration/logger"
	"l03/internal/infra"
	"l03/internal/infra/api/controller/auction_controller"
	"l03/internal/infra/api/controller/bid_controller"
	"l03/internal/infra/api/controller/user_controller"
	"l03/internal/infra/repository/auction"
	"l03/internal/infra/repository/bid"
	"l03/internal/infra/repository/user"
	"l03/internal/usecase/auction_usecase"
	"l03/internal/usecase/bid_usecase"
	"l03/internal/usecase/user_usecase"
	"l03/worker"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	defer logger.Sync()

	if err := loadEnv(); err != nil {
		log.Fatal(err)
		return
	}

	dbConnection, err := mongodb.NewMongoDbConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err.Error())
		return
	}

	logger.Info("Successfully connected to MongoDB", zap.String("database", dbConnection.Name()))

	// Desligamento com segurança
	workerCtx, cancelWorker := context.WithCancel(ctx)
	defer cancelWorker()

	// Repositories
	userRepository := user.NewRepository(dbConnection)
	auctionRepository := auction.NewRepository(dbConnection)
	bidRepository := bid.NewRepository(dbConnection, auctionRepository)

	// Use cases
	userUseCase := user_usecase.NewUseCase(userRepository)
	bidUseCase := bid_usecase.NewUseCase(workerCtx, bidRepository)
	auctionUseCase := auction_usecase.NewUseCase(auctionRepository, bidRepository)

	// Uma go rotina para o worker trabalhar livremente.
	go worker.NewAuctionWorker(auctionUseCase).Start(workerCtx)

	userCtl, auctionCtl, bidCtl := initControllers(userUseCase, auctionUseCase, bidUseCase)
	router := infra.GetRoute(auctionCtl, bidCtl, userCtl)
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Outra go rotina para o servidor http trabalhar livremente.
	go func() {
		logger.Info("Starting http server on port :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Could not listen on %s: %v\n", server.Addr, err)
		}
	}()

	// Ponto uníco para encerrar o sistema, isso inclui todas as partes do sistema.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	logger.Info("Shutting down server and worker...")
	cancelWorker() // bye worker

	bidUseCase.Shutdown(context.Background())

	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancelShutdown()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	logger.Info("Server exiting")
}

func initControllers(
	userUseCase user_usecase.UserUseCaseInterface,
	auctionUseCase auction_usecase.AuctionUseCaseInterface,
	bidUseCase bid_usecase.BidUseCaseInterface,
) (
	userController *user_controller.UserController,
	auctionController *auction_controller.AuctionController,
	bidController *bid_controller.BidController,
) {
	userController = user_controller.NewController(userUseCase)
	bidController = bid_controller.NewController(bidUseCase)
	auctionController = auction_controller.NewController(auctionUseCase)

	return
}

func loadEnv() error {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return err
	}

	envPath := projectRoot + "/cmd/auction/.env"
	if err := godotenv.Load(envPath); err != nil {
		return errors.New("error loading .env file from " + envPath + ": " + err.Error())
	}
	return nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(dir + "/go.mod"); err == nil {
			return dir, nil
		}

		parent := dir + "/.."
		if dir == parent {
			return "", errors.New("go.mod not found, puf")
		}
		dir = parent
	}
}
