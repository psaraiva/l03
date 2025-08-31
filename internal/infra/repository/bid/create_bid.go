package bid

import (
	"context"
	"l03/internal/entity/auction_entity"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"sync"

	"go.uber.org/zap"
)

func (br *BidRepository) Create(ctx context.Context, bidEntities []bid_entity.Bid) *internal_error.InternalError {
	var wg sync.WaitGroup

	for _, item := range bidEntities {
		wg.Add(1)
		go func(entity bid_entity.Bid) {
			defer wg.Done()
			auctionEntity, err := br.AuctionRepository.FindById(ctx, entity.AuctionID)
			if err != nil {
				br.logger.Error("Error trying to find auction by id when creating bid", err,
					zap.String("auctionId", entity.AuctionID),
					zap.String("bidId", entity.ID))
				return
			}

			if auctionEntity.Status != auction_entity.Running {
				br.logger.Error("Auction is not runnig, cannot create bid", nil,
					zap.String("auctionId", entity.AuctionID),
					zap.String("bidId", entity.ID),
					zap.Int("auctionStatus", int(auctionEntity.Status)))
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
				br.logger.Error("Error trying to insert bid into database", err,
					zap.String("bidId", entity.ID))
				return
			}
		}(item)
	}

	wg.Wait()
	return nil
}
