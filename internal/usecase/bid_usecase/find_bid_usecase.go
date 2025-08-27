package bid_usecase

import (
	"context"
	"l03/internal/internal_error"
)

func (buc *BidUseCase) FindByAuctionId(ctx context.Context, auctionID string) (*[]BidOutputDTO, *internal_error.InternalError) {
	collection, err := buc.BidRepository.FindByAuctionId(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	listDTO := make([]BidOutputDTO, len(collection))
	for i, item := range collection {
		listDTO[i] = BidOutputDTO{
			ID:        item.ID,
			UserID:    item.UserID,
			AuctionID: item.AuctionID,
			Amount:    item.Amount,
			Timestamp: item.Timestamp,
		}
	}

	return &listDTO, nil
}

func (buc *BidUseCase) FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError) {
	entity, err := buc.BidRepository.FindWinningBidByAuctionId(ctx, auctionID)
	if err != nil {
		return nil, err
	}

	return &BidOutputDTO{
		ID:        entity.ID,
		UserID:    entity.UserID,
		AuctionID: entity.AuctionID,
		Amount:    entity.Amount,
		Timestamp: entity.Timestamp,
	}, nil
}
