package bid_usecase

import (
	"context"
	"l03/configuration/logger"
	"l03/internal/entity/bid_entity"
	"l03/internal/internal_error"
	"l03/internal/usecase"
	"os"
	"strconv"
	"sync"
	"time"

	"go.uber.org/zap"
)

type bidBatchProcessor struct {
	bidRepository bid_entity.BidRepositoryInterface
	timer         *time.Timer
	batchSize     int
	batchDuration time.Duration
	bidChannel    chan bid_entity.Bid
	wg            sync.WaitGroup
}

func newBidBatchProcessor(
	ctx context.Context,
	bidRepository bid_entity.BidRepositoryInterface) *bidBatchProcessor {
	batchSize := getBatchSize()
	batchDuration := getBatchDuration()

	processor := &bidBatchProcessor{
		bidRepository: bidRepository,
		batchSize:     batchSize,
		batchDuration: batchDuration,
		bidChannel:    make(chan bid_entity.Bid, batchSize),
	}
	processor.timer = time.NewTimer(processor.batchDuration)

	processor.wg.Add(1)
	processor.triggerCreateRoutine(ctx)
	return processor
}

func (p *bidBatchProcessor) shutdown() {
	logger.Info("Closing bid channel and flushing remaining bids...")
	close(p.bidChannel)
	p.wg.Wait()
	logger.Info("Bid processor has been shut down.")
}

func (p *bidBatchProcessor) triggerCreateRoutine(ctx context.Context) {
	go func() {
		defer p.wg.Done()

		bidBatch := make([]bid_entity.Bid, 0, p.batchSize)

	loop:
		// maybe 4(ever)
		for {
			select {
			case bidEntity, ok := <-p.bidChannel:
				if !ok {
					if len(bidBatch) > 0 {
						logger.Info("Flushing final bid batch before shutdown", zap.Int("batch_size", len(bidBatch)))
						if err := p.bidRepository.Create(context.Background(), bidBatch); err != nil {
							logger.Error("error trying to process final bid batch list", err)
						}
					}
					logger.Info("Bid processing routine has finished.")
					return
				}

				bidBatch = append(bidBatch, bidEntity)

				if len(bidBatch) >= p.batchSize {
					logger.Info("Processing bid batch due to size", zap.Int("batch_size", len(bidBatch)))
					if err := p.bidRepository.Create(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					bidBatch = make([]bid_entity.Bid, 0, p.batchSize)
					p.timer.Reset(p.batchDuration)
				}
			case <-p.timer.C:
				if len(bidBatch) > 0 {
					logger.Info("Processing bid batch due to timeout", zap.Int("batch_size", len(bidBatch)))
					if err := p.bidRepository.Create(ctx, bidBatch); err != nil {
						logger.Error("error trying to process bid batch list", err)
					}
					bidBatch = make([]bid_entity.Bid, 0, p.batchSize)
				}
				p.timer.Reset(p.batchDuration)
			case <-ctx.Done():
				p.timer.Stop()
				logger.Info("Main context canceled in bid routine, switching to shutdown mode.")
				break loop
			}
		}

		// Modo de desligamento: esvazia o canal, envia os lances restantes e boa sorte
		logger.Info("Ops... Closing bid channel prepare for final flush...")
		for bidEntity := range p.bidChannel {
			bidBatch = append(bidBatch, bidEntity)
		}

		if len(bidBatch) > 0 {
			logger.Info("Flushing final bid batch before shutdown", zap.Int("batch_size", len(bidBatch)))
			if err := p.bidRepository.Create(context.Background(), bidBatch); err != nil {
				logger.Error("error trying to process final bid batch list", err)
			}
		}
		logger.Info("Bid processing routine has finished.")
	}()
}

func (buc *BidUseCase) Create(ctx context.Context, bidInputDTO BidInputDTO) (*usecase.IDOutputDTO, *internal_error.InternalError) {
	entity, err := bid_entity.Create(bidInputDTO.UserID, bidInputDTO.AuctionID, bidInputDTO.Amount)
	if err != nil {
		buc.logger.Error("Error trying to create by entity", err,
			zap.String("error_origin", "BidEntity.Create"))
		return nil, err
	}

	buc.batchProcessor.bidChannel <- *entity
	return &usecase.IDOutputDTO{ID: entity.ID}, nil
}

// usado para definir o timeout/intervalo de execução de bid batch
func getBatchDuration() time.Duration {
	batchInterval := os.Getenv("BID_BATCH_BUFFER_DURATION")
	duration, err := time.ParseDuration(batchInterval)
	if err != nil {
		return 3 * time.Minute
	}

	return duration
}

// usado para definir o tamanho do batch de bid batch
func getBatchSize() int {
	batchSize, err := strconv.Atoi(os.Getenv("BID_BATCH_BUFFER_SIZE"))
	if err != nil {
		return 5
	}

	return batchSize
}
