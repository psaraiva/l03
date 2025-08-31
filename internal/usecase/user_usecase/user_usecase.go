package user_usecase

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/user_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
)

type UserUseCase struct {
	userRepository user_entity.UserRepositoryInterface
	logger         *logger.ContextualLogger
}

type UserInputDTO struct {
	Name string `json:"name"`
}

type UserOutputDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserUseCaseInterface interface {
	Create(ctx context.Context, userInput UserInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError)
	FindById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError)
	FindUsers(ctx context.Context) (*[]UserOutputDTO, *internal_error.InternalError)
}

func NewUseCase(userRepository user_entity.UserRepositoryInterface) UserUseCaseInterface {
	return &UserUseCase{
		userRepository: userRepository,
		logger:         logger.WithComponent("usecase-user"),
	}
}
