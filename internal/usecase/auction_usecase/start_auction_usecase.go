package auction_usecase

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
)

func (auc *AuctionUseCase) Start(ctx context.Context, id string) *internal_error.InternalError {
	auction, err := auc.auctionRepository.FindById(ctx, id)
	if err != nil {
		auc.logger.Error("Error trying to find by id by repository", err)
		return err
	}

	if auction.Status != auction_entity.Active {
		auc.logger.Error("Error trying to start auction not in status active", err)
		return internal_error.NewBadRequestError("Auction is not in Active state")
	}

	return auc.auctionRepository.ChangeStatus(ctx, id, auction_entity.Running)
}
