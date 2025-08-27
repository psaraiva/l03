package auction

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
)

func (ar *AuctionRepository) Create(ctx context.Context, auction auction_entity.Auction) *internal_error.InternalError {
	entityMongo := &AuctionEntityMongo{
		ID:          auction.ID,
		ProductName: auction.ProductName,
		Category:    auction.Category,
		Description: auction.Description,
		Condition:   auction.Condition,
		Status:      auction.Status,
		Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, entityMongo)
	if err != nil {
		logger.Error("repository.auction.Create.err", err)
		return internal_error.NewInternalServerError("error trying to create auction")
	}

	return nil
}
