package auction_usecase

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase/bid_usecase"
)

func (au *AuctionUseCase) FindById(ctx context.Context, id string) (*AuctionOutputDTO, *internal_error.InternalError) {
	entity, err := au.auctionRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &AuctionOutputDTO{
		ID:          entity.ID,
		ProductName: entity.ProductName,
		Category:    entity.Category,
		Description: entity.Description,
		Condition:   ProductCondition(entity.Condition),
		Status:      AuctionStatus(entity.Status),
		Timestamp:   entity.Timestamp,
	}, nil
}

func (au *AuctionUseCase) FindActions(ctx context.Context, status AuctionStatus, category, productName string) (*[]AuctionOutputDTO, *internal_error.InternalError) {
	collection, err := au.auctionRepository.FindAuctions(ctx, auction_entity.AuctionStatus(status), category, productName)
	if err != nil {
		return nil, err
	}

	listDTO := make([]AuctionOutputDTO, len(*collection))
	for i, item := range *collection {
		listDTO[i] = AuctionOutputDTO{
			ID:          item.ID,
			ProductName: item.ProductName,
			Category:    item.Category,
			Description: item.Description,
			Condition:   ProductCondition(item.Condition),
			Status:      AuctionStatus(item.Status),
			Timestamp:   item.Timestamp,
		}
	}

	return &listDTO, nil
}

func (au *AuctionUseCase) FindWinnigBidByAuctionId(ctx context.Context, auctionId string) (*WinningInfoOutputDTO, *internal_error.InternalError) {
	auctionEntity, err := au.auctionRepository.FindById(ctx, auctionId)
	if err != nil {
		return nil, err
	}

	auctionOutputDTO := AuctionOutputDTO{
		ID:          auctionEntity.ID,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   ProductCondition(auctionEntity.Condition),
		Status:      AuctionStatus(auctionEntity.Status),
		Timestamp:   auctionEntity.Timestamp,
	}

	bidEntity, err := au.bidRepository.FindWinningBidByAuctionId(ctx, auctionEntity.ID)
	if err != nil {
		return &WinningInfoOutputDTO{
			Auction: auctionOutputDTO,
			Bid:     nil,
		}, nil
	}

	bidOutputDTO := &bid_usecase.BidOutputDTO{
		ID:        bidEntity.ID,
		AuctionID: bidEntity.AuctionID,
		UserID:    bidEntity.UserID,
		Amount:    bidEntity.Amount,
		Timestamp: bidEntity.Timestamp,
	}

	return &WinningInfoOutputDTO{
		Auction: auctionOutputDTO,
		Bid:     bidOutputDTO,
	}, nil
}
