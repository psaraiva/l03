package user_controller

import "l03/internal/usecase/user_usecase"

type UserController struct {
	userUseCase user_usecase.UserUseCaseInterface
}

func NewController(userUseCase user_usecase.UserUseCaseInterface) *UserController {
	return &UserController{
		userUseCase: userUseCase,
	}
}
