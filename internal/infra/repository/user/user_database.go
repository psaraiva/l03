package user

import (
	"l03/configuration/logger"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserEntityMongo struct {
	ID        string `bson:"_id"`
	Name      string `bson:"name"`
	Timestamp int64  `bson:"timestamp"`
}

type UserRepository struct {
	Collection *mongo.Collection
	logger     *logger.ContextualLogger
}

func NewRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{
		Collection: database.Collection("users"),
		logger:     logger.WithComponent("repository-user"),
	}
}
