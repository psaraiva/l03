package auction

import (
	"context"
	"fmt"
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (ar *AuctionRepository) FindById(ctx context.Context, id string) (*auction_entity.Auction, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var auctionEntityMongo AuctionEntityMongo
	if err := ar.Collection.FindOne(ctx, filter).Decode(&auctionEntityMongo); err != nil {
		logger.Error(fmt.Sprintf("repository.auction.FindAuctionById.err.id = %s", id), err)
		return nil, internal_error.NewInternalServerError("error trying to find auction")
	}

	if auctionEntityMongo.ID == "" {
		return nil, internal_error.NewNotFoundError("auction not found")
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
		logger.Error("repository.auction.FindAuction.err", err)
		return nil, internal_error.NewInternalServerError("error trying to find auctions")
	}
	defer cursor.Close(ctx)

	var auctionEntityMongo []AuctionEntityMongo
	if err := cursor.All(ctx, &auctionEntityMongo); err != nil {
		logger.Error("repository.auction.FindAuction.err", err)
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
