package bid_entity

import (
	"context"
	"l03/internal/internal_error"
	"time"

	"github.com/google/uuid"
)

type Bid struct {
	ID        string
	UserID    string
	AuctionID string
	Amount    float64
	Timestamp time.Time
}

type BidRepositoryInterface interface {
	Create(ctx context.Context, bidEntities []Bid) *internal_error.InternalError
	FindByAuctionId(ctx context.Context, auctionID string) ([]Bid, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*Bid, *internal_error.InternalError)
}

func Create(userId, auctionId string, amount float64) (*Bid, *internal_error.InternalError) {
	bid := &Bid{
		ID:        uuid.New().String(),
		UserID:    userId,
		AuctionID: auctionId,
		Amount:    amount,
		Timestamp: time.Now(),
	}

	if err := bid.Validade(); err != nil {
		return nil, err
	}

	return bid, nil
}

func (be *Bid) Validade() *internal_error.InternalError {
	if err := uuid.Validate(be.UserID); err != nil {
		return internal_error.NewBadRequestError("Invalid User ID")
	}

	if err := uuid.Validate(be.AuctionID); err != nil {
		return internal_error.NewBadRequestError("Invalid Action ID")
	}

	if be.Amount <= 0 {
		return internal_error.NewBadRequestError("Invalid Amount value")
	}

	return nil
}
