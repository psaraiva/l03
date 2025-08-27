package main

import (
	"context"
	"log"

	"l03/configuration/database/mongodb"
	"l03/configuration/logger"
	"l03/internal/infra/api/controller/auction_controller"
	"l03/internal/infra/api/controller/bid_controller"
	"l03/internal/infra/api/controller/user_controller"
	"l03/internal/infra/repository/auction"
	"l03/internal/infra/repository/bid"
	"l03/internal/infra/repository/user"
	"l03/internal/usecase/auction_usecase"
	"l03/internal/usecase/bid_usecase"
	"l03/internal/usecase/user_usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	ctx := context.Background()
	defer logger.Sync()

	//if err := godotenv.Load("cmd/auction/.env"); err != nil {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file:", err)
		return
	}

	dbConnection, err := mongodb.NewMongoDbConnection(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err.Error())
		return
	}
	log.Println("Successfully connected to MongoDB:", dbConnection.Name())

	router := gin.Default()
	userCtl, auctionCtl, bidCtl := initDependencies(dbConnection)

	// Auctions
	router.GET("/auctions", auctionCtl.FindAuctions)
	router.GET("/auctions/:id", auctionCtl.FindById)
	router.POST("/auctions", auctionCtl.CreateAuction)
	router.GET("/auctions/winner/:auction-id", auctionCtl.FindWinningBidbyAuctionId)

	// Bids
	router.POST("/bids", bidCtl.CreateBid)
	router.GET("/bids/:auction-id", bidCtl.FindBidByAuctionId)

	// Users
	router.GET("/users", userCtl.FindUsers)
	router.POST("/users", userCtl.Create)
	router.GET("/users/:user-id", userCtl.FindUserById)

	router.Run(":8080")
}

func initDependencies(database *mongo.Database) (
	userController *user_controller.UserController,
	auctionController *auction_controller.AuctionController,
	bidController *bid_controller.BidController,
) {
	// Repositories
	userRepository := user.NewRepository(database)
	auctionRepository := auction.NewRepository(database)
	bidRepository := bid.NewRepository(database, auctionRepository)

	// Use cases
	userUseCase := user_usecase.NewUseCase(userRepository)
	bidUseCase := bid_usecase.NewUseCase(bidRepository)
	auctionUseCase := auction_usecase.NewUseCase(auctionRepository, bidRepository)

	// Controllers
	userController = user_controller.NewController(userUseCase)
	bidController = bid_controller.NewController(bidUseCase)
	auctionController = auction_controller.NewController(auctionUseCase)

	return
}
