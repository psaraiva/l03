package user_usecase

import (
	"context"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
)

func (uuc *UserUseCase) Create(ctx context.Context, userInput UserInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError) {
	entity, err := user_entity.Create(userInput.Name)
	if err != nil {
		return nil, err
	}

	if err := uuc.userRepository.Create(ctx, *entity); err != nil {
		return nil, err
	}

	return &usecase.IDOutputDTO{ID: entity.ID}, nil
}
