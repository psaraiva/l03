package user_usecase

import (
	"context"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"

	"go.uber.org/zap"
)

func (uuc *UserUseCase) Create(ctx context.Context, userInput UserInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError) {
	entity, err := user_entity.Create(userInput.Name)
	if err != nil {
		uuc.logger.Error("Error trying to create user by entity", err,
			zap.String("name", userInput.Name),
			zap.String("error_origin", "UserEntity.Create"))
		return nil, err
	}

	if err := uuc.userRepository.Create(ctx, *entity); err != nil {
		uuc.logger.Error("Error trying to create user by repository", err,
			zap.String("error_origin", "UserRepository.Create"))
		return nil, err
	}

	return &usecase.IDOutputDTO{ID: entity.ID}, nil
}
