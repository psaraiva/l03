package auction

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
)

func (ar *AuctionRepository) ChangeStatus(ctx context.Context, id string, status auction_entity.AuctionStatus) *internal_error.InternalError {
	_, err := ar.Collection.UpdateOne(
		ctx,
		bson.M{"_id": id},
		bson.M{"$set": bson.M{"status": status}},
	)

	if err != nil {
		logger.Error("repository.auction.ChangeStatus.err", err)
		return internal_error.NewInternalServerError("error trying auction status change")
	}

	return nil
}
