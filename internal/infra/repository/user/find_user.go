package user

import (
	"context"
	"errors"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"
)

func (ur *UserRepository) FindById(ctx context.Context, id string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var user UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		ur.logger.Error("Error trying to find user by id", err, zap.String("userId", id))
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal_error.NewNotFoundError("user not found")
		}
		return nil, internal_error.NewInternalServerError("error trying to find user by id")
	}

	return &user_entity.User{
		ID:   user.ID,
		Name: user.Name,
	}, nil
}

func (ur *UserRepository) FindUsers(ctx context.Context) ([]user_entity.User, *internal_error.InternalError) {
	cursor, err := ur.Collection.Find(ctx, bson.M{})
	if err != nil {
		ur.logger.Error("Error trying to find users", err)
		return nil, internal_error.NewInternalServerError("error trying to find users")
	}
	defer cursor.Close(ctx)

	var userEntityMongo []UserEntityMongo
	if err := cursor.All(ctx, &userEntityMongo); err != nil {
		ur.logger.Error("Error decoding users from cursor", err)
		return nil, internal_error.NewInternalServerError("error trying to find users")
	}

	list := make([]user_entity.User, len(userEntityMongo))
	for key, item := range userEntityMongo {
		list[key] = user_entity.User{
			ID:   item.ID,
			Name: item.Name,
		}
	}

	return list, nil
}
