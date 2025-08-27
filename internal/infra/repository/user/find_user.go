package user

import (
	"context"
	"errors"
	"fmt"
	"l03/configuration/logger"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func (ur *UserRepository) FindById(ctx context.Context, id string) (*user_entity.User, *internal_error.InternalError) {
	filter := bson.M{"_id": id}

	var user UserEntityMongo
	err := ur.Collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		logger.Error(fmt.Sprintf("repository.user.FindById.err.id='%s'", id), err)
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, internal_error.NewNotFoundError(fmt.Sprintf("user not found by id = %s", id))
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
		logger.Error("repository.user.FindUsers.err", err)
		return nil, internal_error.NewInternalServerError("error trying to find users")
	}
	defer cursor.Close(ctx)

	var userEntityMongo []UserEntityMongo
	if err := cursor.All(ctx, &userEntityMongo); err != nil {
		logger.Error("repository.user.FindUsers.err", err)
		return nil, internal_error.NewInternalServerError("Error trying to find users")
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
