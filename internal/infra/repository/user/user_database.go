package user

import "go.mongodb.org/mongo-driver/mongo"

type UserEntityMongo struct {
	ID   string `bson:"_id"`
	Name string `bson:"name"`
	// @TODO adicionar timestamp de quando o requistro foi criado (metricas)
}

type UserRepository struct {
	Collection *mongo.Collection
}

func NewRepository(database *mongo.Database) *UserRepository {
	return &UserRepository{Collection: database.Collection("users")}
}
