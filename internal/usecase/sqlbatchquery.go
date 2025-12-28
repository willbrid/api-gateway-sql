package usecase

import (
	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/pkg/confighelper"
	"github.com/willbrid/api-gateway-sql/internal/pkg/mapperfieldshelper"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/pkg/csvstream"
	"github.com/willbrid/api-gateway-sql/pkg/database/external"
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"context"
	"errors"
	"strings"
	"sync"
)

var (
	errBatchModeNotActivated = errors.New("attribut multi for batch mode is not activate for this target")
)

type SQLBatchQueryUsecase struct {
	sqlQueryRepo  *repository.SQLQueryRepo
	batchStatRepo *repository.BatchStatRepo
	blockRepo     *repository.BlockRepo
	config        *config.Config
}

func NewSQLBatchQueryUsecase(sqlQueryRepo *repository.SQLQueryRepo, batchStatRepo *repository.BatchStatRepo, blockRepo *repository.BlockRepo, config *config.Config) *SQLBatchQueryUsecase {
	return &SQLBatchQueryUsecase{sqlQueryRepo, batchStatRepo, blockRepo, config}
}

func (squ *SQLBatchQueryUsecase) ExecuteBatch(ctx context.Context, sqlbatchquery *domain.SQLBatchQueryInput) error {
	target, cfgdb, err := confighelper.GetTargetAndDatabase(squ.config, sqlbatchquery.TargetName)
	if err != nil {
		return err
	}

	if !target.Multi {
		return errBatchModeNotActivated
	}

	batchStat, err := squ.batchStatRepo.Create(ctx, target.Name)
	if err != nil {
		return err
	}

	blockChannel, errorChannel := csvstream.ReadCSVInBlock(sqlbatchquery.File, target.BufferSize)
	var openChannels int = 2
	var wg sync.WaitGroup

	for openChannels > 0 {
		select {
		case block, open := <-blockChannel:
			if !open {
				openChannels--
				continue
			}

			blockDataInput := &domain.BlockDataInput{
				BSInput: batchStat,
				BLInput: block,
				TGInput: target,
				DBInput: cfgdb,
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				squ.processBlock(ctx, blockDataInput)
			}()

		case err, open := <-errorChannel:
			if !open {
				openChannels--
				continue
			}

			if err != nil {
				if updateErr := squ.batchStatRepo.UpdateLastCompleted(ctx, batchStat); updateErr != nil {
					return updateErr
				}

				return err
			}
		}
	}

	wg.Wait()
	if updateErr := squ.batchStatRepo.UpdateLastCompleted(ctx, batchStat); updateErr != nil {
		return updateErr
	}

	return nil
}

func (squ *SQLBatchQueryUsecase) processBlock(ctx context.Context, blockDataInput *domain.BlockDataInput) {
	newBlock := domain.NewBlock(blockDataInput.BLInput.StartLine, blockDataInput.BLInput.EndLine)
	newBlock, err := squ.batchStatRepo.AddBlockToBatchStat(ctx, blockDataInput.BSInput, newBlock)
	if err != nil {
		logger.Error("failed to process block : %v", err.Error())
		return
	}

	cnx, err := external.NewDatabase(*blockDataInput.DBInput)
	if err != nil {
		return
	}

	squ.sqlQueryRepo.SetDB(cnx)
	defer squ.sqlQueryRepo.CloseDB()

	var (
		wg                sync.WaitGroup
		batchSize         int      = blockDataInput.TGInput.BatchSize
		batchFields       []string = strings.Split(blockDataInput.TGInput.BatchFields, ";")
		currentBufferSize int      = len(blockDataInput.BLInput.Lines)
		numBatches        int      = currentBufferSize / batchSize
	)

	if currentBufferSize%batchSize != 0 || currentBufferSize < batchSize {
		numBatches++
	}

	for i := 0; i < numBatches; i++ {
		start := i * batchSize
		end := start + batchSize
		if end > currentBufferSize {
			end = currentBufferSize
		}

		batch := blockDataInput.BLInput.Lines[start:end]
		wg.Add(1)

		go func() {
			defer wg.Done()
			var (
				record  map[string]any
				records []map[string]any = make([]map[string]any, 0, len(batch))
				err     error
			)

			for _, line := range batch {
				record, err = mapperfieldshelper.MapBatchFieldToValueLine(batchFields, line)

				if err != nil {
					logger.Error("failed to process batch in block : %v", err.Error())
					break
				} else {
					records = append(records, record)
				}
			}

			if len(records) > 0 {
				err = squ.sqlQueryRepo.ExecuteBatch(ctx, blockDataInput.TGInput.SqlQuery, records)
				if err != nil {
					logger.Error("failed to process batch in block : %v", err.Error())
					failureRange := domain.NewFailureRange(start, end)
					if err = squ.blockRepo.Update(ctx, newBlock, failureRange, false); err != nil {
						logger.Error("failed to process batch in block : %v", err.Error())
						return
					}
				} else {
					if err = squ.blockRepo.Update(ctx, newBlock, nil, true); err != nil {
						logger.Error("failed to process batch in block : %v", err.Error())
						return
					}
				}
			}
		}()
	}

	wg.Wait()
}
