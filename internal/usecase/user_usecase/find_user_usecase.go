package user_usecase

import (
	"context"
	"l03/internal/internal_error"
)

func (uuc *UserUseCase) FindById(ctx context.Context, id string) (*UserOutputDTO, *internal_error.InternalError) {
	entity, err := uuc.userRepository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &UserOutputDTO{
		ID:   entity.ID,
		Name: entity.Name,
	}, nil
}

func (uuc *UserUseCase) FindUsers(ctx context.Context) (*[]UserOutputDTO, *internal_error.InternalError) {
	collection, err := uuc.userRepository.FindUsers(ctx)
	if err != nil {
		return nil, err
	}

	list := make([]UserOutputDTO, len(collection))
	for i, item := range collection {
		list[i] = UserOutputDTO{
			ID:   item.ID,
			Name: item.Name,
		}
	}

	return &list, nil
}
