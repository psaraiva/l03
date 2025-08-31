package auction

import (
	"context"
	"errors"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (ar *AuctionRepository) FindById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		ar.logger.Error("error trying to find auction by id", err, zap.String("auctionId", id))
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal_error.NewNotFoundError("auction not found")
		}
		return nil, internal_error.NewInternalServerError("error trying to find auction by id")
	}

	return &auction_entity.Auction{
		ID:          auctionEntityMongo.ID,
		ProductName: auctionEntityMongo.ProductName,
		Category:    auctionEntityMongo.Category,
		Description: auctionEntityMongo.Description,
		Condition:   auctionEntityMongo.Condition,
		Status:      auctionEntityMongo.Status,
		Timestamp:   time.Unix(auctionEntityMongo.Timestamp, 0),
	}, nil
}

func (ar *AuctionRepository) FindAuctions(
	ctx context.Context,
	status auction_entity.AuctionStatus,
	category, productName string) (*[]auction_entity.Auction, *internal_error.InternalError) {

	filter := bson.M{}
	if status != 0 {
		filter["status"] = status
	}

	if category != "" {
		filter["category"] = category
	}

	if productName != "" {
		filter["product_name"] = primitive.Regex{
			Pattern: productName,
			Options: "i",
		}
	}

	cursor, err := ar.Collection.Find(ctx, filter)
	if err != nil {
		ar.logger.Error("error trying to find auctions", err)
		return nil, internal_error.NewInternalServerError("error trying to find auctions")
	}
	defer cursor.Close(ctx)

	var auctionEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionEntityMongo); err != nil {
		ar.logger.Error("error decoding auctions from cursor", err)
		return nil, internal_error.NewInternalServerError("Error trying to find auctions")
	}

	var auctions []auction_entity.Auction
	for _, auction := range auctionEntityMongo {
		auctions = append(auctions, auction_entity.Auction{
			ID:          auction.ID,
			ProductName: auction.ProductName,
			Category:    auction.Category,
			Description: auction.Description,
			Condition:   auction.Condition,
			Status:      auction.Status,
			Timestamp:   time.Unix(auction.Timestamp, 0),
		})
	}

	return &auctions, nil
}
