package bid_usecase

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
	"os"
	"strconv"
	"time"
)

var bidBatch []bid_entity.Bid

type BidUseCaseInterface interface {
	Create(ctx context.Context, bidInputDTO BidInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError)
	FindWinningBidByAuctionId(ctx context.Context, auctionID string) (*BidOutputDTO, *internal_error.InternalError)
	FindByAuctionId(ctx context.Context, auctionID string) (*[]BidOutputDTO, *internal_error.InternalError)
}

func (buc *BidUseCase) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer close(buc.bidChannel)

		// for(ever)
		for {
			select {
			// case normal operation
			case bidEntity, ok := <-buc.bidChannel:
				// error
				if !ok {
					if len(bidBatch) > 0 {
						if err := buc.BidRepository.Create(ctx, bidBatch); err != nil {
							logger.Error("error trying to process bid batch list", err)
						}
					}
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				// bid batch overflow
				if len(bidBatch) >= buc.maxBatchSize {
					if err := buc.BidRepository.Create(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					bidBatch = nil
					buc.timer.Reset(buc.batchInsertInterval)
				}
			// case timeout operation
			case <-buc.timer.C:
				if err := buc.BidRepository.Create(ctx, bidBatch); err != nil {
					logger.Error("error trying to process bid batch list", err)
				}
				bidBatch = nil
				buc.timer.Reset(buc.batchInsertInterval)
			}
		}
	}()
}

func (buc *BidUseCase) Create(ctx context.Context, bidInputDTO BidInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError) {
	entity, err := bid_entity.Create(bidInputDTO.UserID, bidInputDTO.AuctionID, bidInputDTO.Amount)
	if err != nil {
		return nil, err
	}

	buc.bidChannel <- *entity
	return &usecase.IDOutputDTO{ID: entity.ID}, nil
}

func getMaxBathSizeInterval() time.Duration {
	batchInsertInterval := os.Getenv("BATCH_INSERT_INTERVAL")
	duration, err := time.ParseDuration(batchInsertInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

func getMaxBatchSize() int {
	batchSize, err := strconv.Atoi(os.Getenv("MAX_BATCH_SIZE"))
	if err != nil {
		return 5
	}

	return batchSize
}
