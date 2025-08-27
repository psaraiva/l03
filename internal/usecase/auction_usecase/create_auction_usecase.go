package auction_usecase

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
)

func (au *AuctionUseCase) Create(ctx context.Context, auctionInput AuctionInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError) {
	entity, err := auction_entity.Create(
		auctionInput.ProductName,
		auctionInput.Category,
		auctionInput.Description,
		auction_entity.ProductCondition(auctionInput.Condition))

	if err != nil {
		return nil, err
	}

	if err := au.auctionRepository.Create(ctx, *entity); err != nil {
		return nil, err
	}

	return &usecase.IDOutputDTO{ID: entity.ID}, nil
}
