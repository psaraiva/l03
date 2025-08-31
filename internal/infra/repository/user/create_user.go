package user

import (
	"context"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"
	"time"

	"go.uber.org/zap"
)

func (ur *UserRepository) Create(ctx context.Context, user user_entity.User) *internal_error.InternalError {
	userEntityMongo := &UserEntityMongo{
		ID:        user.ID,
		Name:      user.Name,
		Timestamp: time.Now().Unix(),
	}

	_, err := ur.Collection.InsertOne(ctx, userEntityMongo)
	if err != nil {
		ur.logger.Error("Error trying to create user in database", err,
			zap.String("userId", user.ID))
		return internal_error.NewInternalServerError("error trying to create user")
	}

	return nil
}
