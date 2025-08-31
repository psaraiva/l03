package worker

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/auction_entity"
	"l03/internal/usecase/auction_usecase"
	"os"
	"time"

	"go.uber.org/zap"
)

type AuctionWorker struct {
	auctionUseCase auction_usecase.AuctionUseCaseInterface
	logger         *logger.ContextualLogger
}

func NewAuctionWorker(auctionUseCase auction_usecase.AuctionUseCaseInterface) WorkerInterface {
	return &AuctionWorker{
		auctionUseCase: auctionUseCase,
		logger:         logger.WithComponent("worker"),
	}
}

func (aw *AuctionWorker) doWork(ctx context.Context) {
	aw.logger.Info("Running auction check")

	auctions, err := aw.auctionUseCase.FindActions(ctx, auction_usecase.AuctionStatus(auction_entity.Active), "", "")
	if err != nil {
		aw.logger.Error("Error trying to find auctions", err)
		return
	}

	if len(*auctions) == 0 {
		aw.logger.Info("Auctions not found to process")
		return
	}

	for _, auction := range *auctions {
		if err := aw.auctionUseCase.Start(ctx, auction.ID); err != nil {
			aw.logger.Error("Error trying to start auction", err, zap.String("auctionId", auction.ID))
			continue
		}

		// Uma go rotina cuida do fechando no tempo certo
		go func(id string) {
			time.Sleep(getAuctionDuration())
			if err := aw.auctionUseCase.End(context.Background(), id); err != nil {
				aw.logger.Error("Error trying to end auction", err, zap.String("auctionId", id))
			}

			aw.logger.Info("Auction successfully closed", zap.String("auctionId", id))
		}(auction.ID)
	}

	aw.logger.Info("Finished processing auctions")
}

func (aw *AuctionWorker) Start(ctx context.Context) {
	aw.doWork(ctx)
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// 4(ever)
	for {
		select {
		case <-ctx.Done():
			aw.logger.Info("Shutting down by context")
			return
		case <-ticker.C:
			aw.doWork(ctx)
		}
	}
}

func getAuctionDuration() time.Duration {
	durationStr := os.Getenv("AUCTION_DURATION_DEFAULT")
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		logger.Info("AUCTION_DURATION_DEFAULT env var not set or invalid, using default of 5 minutes")
		return 5 * time.Minute
	}

	logger.Info("Using auction duration", zap.String("duration", duration.String()))
	return duration
}
