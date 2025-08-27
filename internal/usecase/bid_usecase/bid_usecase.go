package bid_usecase

import (
	"context"
	"l03/internal/entity/bid_entity"
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
	BidRepository       bid_entity.BidRepositoryInterface
	timer               *time.Timer
	maxBatchSize        int
	batchInsertInterval time.Duration
	bidChannel          chan bid_entity.Bid
}

func NewUseCase(bidRepository bid_entity.BidRepositoryInterface) *BidUseCase {
	maxBathSizeInterval := getMaxBathSizeInterval()
	maxBatchSize := getMaxBatchSize()

	bidUseCase := &BidUseCase{
		BidRepository:       bidRepository,
		timer:               time.NewTimer(maxBathSizeInterval),
		maxBatchSize:        maxBatchSize,
		batchInsertInterval: maxBathSizeInterval,
		bidChannel:          make(chan bid_entity.Bid, maxBatchSize),
	}

	bidUseCase.triggerCreateRoutine(context.Background())
	return bidUseCase
}
