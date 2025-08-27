package auction_usecase

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
	"l03/internal/usecase/bid_usecase"
	"time"
)

type AuctionInputDTO struct {
	ProductName string           `json:"product_name" binding:"required,min=5,max=100"`
	Category    string           `json:"category" binding:"required,min=5,max=100"`
	Description string           `json:"description" binding:"required,min=5,max=100"`
	Condition   ProductCondition `json:"condition"`
}

type AuctionOutputDTO struct {
	ID          string           `json:"id"`
	ProductName string           `json:"product_name"`
	Category    string           `json:"category"`
	Description string           `json:"description"`
	Condition   ProductCondition `json:"condition"`
	Status      AuctionStatus    `json:"status"`
	Timestamp   time.Time        `json:"timestamp" time_format:"2006-01-02 15:04:05"`
}

type WinningInfoOutputDTO struct {
	Auction AuctionOutputDTO          `json:"auction"`
	Bid     *bid_usecase.BidOutputDTO `json:"bid"`
}

type ProductCondition int
type AuctionStatus int

type AuctionUseCase struct {
	auctionRepository auction_entity.AuctionRepositoryInterface
	bidRepository     bid_entity.BidRepositoryInterface
}

type AuctionUseCaseInterface interface {
	Create(ctx context.Context, auctionInput AuctionInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError)
	FindById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError)
	FindActions(ctx context.Context, status AuctionStatus, category, productName string) (*[]AuctionOutputDTO, *internal_error.InternalError)
	FindWinnigBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError)
}

func NewUseCase(
	auctionRepository auction_entity.AuctionRepositoryInterface,
	bidRepository bid_entity.BidRepositoryInterface) AuctionUseCaseInterface {
	return &AuctionUseCase{
		auctionRepository: auctionRepository,
		bidRepository:     bidRepository,
	}
}
