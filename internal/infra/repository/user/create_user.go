package user

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"
)

func (ur *UserRepository) Create(ctx context.Context, user user_entity.User) *internal_error.InternalError {
	userEntityMongo := &UserEntityMongo{
		ID:   user.ID,
		Name: user.Name,
		// @TODO no futuro terá um campo de data de criação do registro
		//Timestamp:   auction.Timestamp.Unix(),
	}

	_, err := ur.Collection.InsertOne(ctx, userEntityMongo)
	if err != nil {
		logger.Error("user.CreateUser.err", err)
		return internal_error.NewInternalServerError("error trying to create user")
	}

	return nil
}
