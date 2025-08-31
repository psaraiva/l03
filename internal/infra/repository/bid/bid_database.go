package bid

import (
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"

	"go.mongodb.org/mongo-driver/mongo"
)

type BidEntityMongo struct {
	ID        string  `bson:"_id"`
	UserID    string  `bson:"user_id"`
	AuctionID string  `bson:"auction_id"`
	Amount    float64 `bson:"amount"`
	Timestamp int64   `bson:"timestamp"`
}

type BidRepository struct {
	Collection        *mongo.Collection
	AuctionRepository auction_entity.AuctionRepositoryInterface
	logger            *logger.ContextualLogger
}

func NewRepository(database *mongo.Database, auctionRepository auction_entity.AuctionRepositoryInterface) *BidRepository {
	return &BidRepository{
		Collection:        database.Collection("bids"),
		AuctionRepository: auctionRepository,
		logger:            logger.WithComponent("repository-bid"),
	}
}
