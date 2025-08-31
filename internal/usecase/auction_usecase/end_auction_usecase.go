package auction_usecase

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/internal_error"
)

func (auc *AuctionUseCase) End(ctx context.Context, id string) *internal_error.InternalError {
	auction, err := auc.auctionRepository.FindById(ctx, id)
	if err != nil {
		auc.logger.Error("Error trying to find by id by repository", err)
		return err
	}

	if auction.Status != auction_entity.Running {
		auc.logger.Error("Error trying to end auction not in status running", err)
		return internal_error.NewBadRequestError("Auction is not in Running state to be completed")
	}

	err = auc.auctionRepository.ChangeStatus(ctx, id, auction_entity.Completed)
	if err != nil {
		auc.logger.Error("Error trying to change status by repository", err)
		return err
	}

	return nil
}
