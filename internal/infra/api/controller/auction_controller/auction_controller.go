package auction_controller

import "l03/internal/usecase/auction_usecase"

type AuctionController struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
}

func NewController(auctionUseCase auction_usecase.AuctionUseCaseInterface) *AuctionController {
	return &AuctionController{
		auctionUseCase: auctionUseCase,
	}
}
