package bid_usecase

import (
	"context"
	"l03/internal/internal_error"

	"go.uber.org/zap"
)

func (buc *BidUseCase) FindByAuctionId(
	ctx context.Context,
	auctionId string,
) (*[]BidOutputDTO, *internal_error.InternalError) {
	bids, err := buc.BidRepository.FindByAuctionId(ctx, auctionId)
	if err != nil {
		buc.logger.Error("Error trying to find bids by auction id", err,
			zap.String("auctionId", auctionId),
			zap.String("error_origin", "BidRepository.FindByAuctionId"))
		return nil, err
	}

	var bidOutputs []BidOutputDTO
	for _, bid := range bids {
		bidOutputs = append(bidOutputs, BidOutputDTO{
			ID:        bid.ID,
			UserID:    bid.UserID,
			AuctionID: bid.AuctionID,
			Amount:    bid.Amount,
			Timestamp: bid.Timestamp,
		})
	}

	return &bidOutputs, nil
}

func (buc *BidUseCase) FindWinningBidByAuctionId(
	ctx context.Context,
	auctionId string,
) (*BidOutputDTO, *internal_error.InternalError) {
	bid, err := buc.BidRepository.FindWinningBidByAuctionId(ctx, auctionId)
	if err != nil {
		buc.logger.Error("Error trying to find winning bid", err,
			zap.String("auctionId", auctionId),
			zap.String("error_origin", "BidRepository.FindWinningBidByAuctionId"))
		return nil, err
	}

	return &BidOutputDTO{
		ID:        bid.ID,
		UserID:    bid.UserID,
		AuctionID: bid.AuctionID,
		Amount:    bid.Amount,
		Timestamp: bid.Timestamp,
	}, nil
}
