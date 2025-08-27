package auction

import (
	"l03/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	ID          string                          `bson:"_id"`
	UserId      string                          `bson:"user_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewRepository(database *mongo.Database) *AuctionRepository {
	return &AuctionRepository{Collection: database.Collection("auctions")}
}
