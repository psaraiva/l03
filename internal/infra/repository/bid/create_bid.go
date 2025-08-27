package bid

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"sync"
)

func (br *BidRepository) Create(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, item := range bidEntities {
		wg.Add(1)
		go func(entity bid_entity.Bid) {
			defer wg.Done()
			auctionEntity, err := br.AuctionRepository.FindById(ctx, entity.AuctionID)
			if err != nil {
				logger.Error("Error trying to find auction by id", err)
				return
			}

			if auctionEntity.Status != auction_entity.Active {
				logger.Error("Error auction is not active", err)
				return
			}

			entityMongo := &BidEntityMongo{
				ID:        entity.ID,
				UserID:    entity.UserID,
				AuctionID: entity.AuctionID,
				Amount:    entity.Amount,
				Timestamp: entity.Timestamp.Unix(),
			}

			if _, err := br.Collection.InsertOne(ctx, entityMongo); err != nil {
				logger.Error("Error trying to insert bid", err)
				return
			}
		}(item)
	}

	wg.Wait()
	return nil
}
