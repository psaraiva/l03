package bid_usecase

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
	"time"
)

type BidInputDTO struct {
	UserID    string  `json:"user_id"`
	AuctionID string  `json:"auction_id"`
	Amount    float64 `json:"amount"`
}

type BidOutputDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	AuctionID string    `json:"auction_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp" time_format:"2006-01-02T15:04:05Z07:00"`
}

type BidUseCase struct {
	BidRepository  bid_entity.BidRepositoryInterface
	batchProcessor *bidBatchProcessor
	logger         *logger.ContextualLogger
}

type BidUseCaseInterface interface {
	Create(ctx context.Context, bidInputDTO BidInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError)
	FindByAuctionId(ctx context.Context, auctionID string) (*[]BidOutputDTO, *internal_error.InternalError)
	Shutdown(ctx context.Context)
}

func NewUseCase(
	ctx context.Context,
	bidRepository bid_entity.BidRepositoryInterface) *BidUseCase {
	bidUseCase := &BidUseCase{
		BidRepository:  bidRepository,
		batchProcessor: newBidBatchProcessor(ctx, bidRepository),
		logger:         logger.WithComponent("usecase-bid"),
	}
	return bidUseCase
}

func (bu *BidUseCase) Shutdown(ctx context.Context) {
	if bu.batchProcessor != nil {
		bu.batchProcessor.shutdown()
	}
}
