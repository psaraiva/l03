package bid_controller

import "l03/internal/usecase/bid_usecase"

type BidController struct {
	BidUseCase bid_usecase.BidUseCaseInterface
}

func NewController(bidUseCase bid_usecase.BidUseCaseInterface) *BidController {
	return &BidController{
		BidUseCase: bidUseCase,
	}
}
