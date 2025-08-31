package infra

import (
	"l03/internal/infra/api/controller/auction_controller"
	"l03/internal/infra/api/controller/bid_controller"
	"l03/internal/infra/api/controller/user_controller"

	"github.com/gin-gonic/gin"
)

func GetRoute(
	auctionCtl *auction_controller.AuctionController,
	bidCtl *bid_controller.BidController,
	userCtl *user_controller.UserController,
) *gin.Engine {
	router := gin.Default()

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

	return router
}
