package bid

import (
	"context"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func (br *BidRepository) FindByAuctionId(ctx context.Context, auctionId string) ([]bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}
	cursor, err := br.Collection.Find(ctx, filter)
	if err != nil {
		br.logger.Error("Error trying to find bids by auctionId", err, zap.String("auctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error trying to find bids by auctionId")
	}

	var bidEntityMongo []BidEntityMongo
	if err := cursor.All(ctx, &bidEntityMongo); err != nil {
		br.logger.Error("Error decoding bids from cursor", err, zap.String("auctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error trying to find bids by auctionId")
	}

	var bids []bid_entity.Bid
	for _, bid := range bidEntityMongo {
		bids = append(bids, bid_entity.Bid{
			ID:        bid.ID,
			UserID:    bid.UserID,
			AuctionID: bid.AuctionID,
			Amount:    bid.Amount,
			Timestamp: time.Unix(bid.Timestamp, 0),
		})
	}
	return bids, nil
}

func (br *BidRepository) FindWinningBidByAuctionId(ctx context.Context, auctionId string) (*bid_entity.Bid, *internal_error.InternalError) {
	filter := bson.M{"auction_id": auctionId}

	var bidEntityMongo BidEntityMongo
	opts := options.FindOne().SetSort(bson.D{{Key: "amount", Value: -1}})
	if err := br.Collection.FindOne(ctx, filter, opts).Decode(&bidEntityMongo); err != nil {
		if err == mongo.ErrNoDocuments {
			br.logger.Info("No winning bid found for auction", zap.String("auctionId", auctionId))
			return nil, internal_error.NewNotFoundError("No winning bid found for auction")
		}

		br.logger.Error("Error trying to find winning bid", err, zap.String("auctionId", auctionId))
		return nil, internal_error.NewInternalServerError("Error trying to find winning bid")
	}

	return &bid_entity.Bid{
		ID:        bidEntityMongo.ID,
		UserID:    bidEntityMongo.UserID,
		AuctionID: bidEntityMongo.AuctionID,
		Amount:    bidEntityMongo.Amount,
		Timestamp: time.Unix(bidEntityMongo.Timestamp, 0),
	}, nil
}
