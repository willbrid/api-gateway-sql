package usecase

import (
	"github.com/rs/zerolog"

	"github.com/willbrid/api-gateway-sql/config"
	"github.com/willbrid/api-gateway-sql/internal/domain"
	"github.com/willbrid/api-gateway-sql/internal/dto"
	"github.com/willbrid/api-gateway-sql/internal/pkg/confighelper"
	"github.com/willbrid/api-gateway-sql/internal/pkg/csvmapper"
	"github.com/willbrid/api-gateway-sql/internal/repository"
	"github.com/willbrid/api-gateway-sql/pkg/csvstream"
	"github.com/willbrid/api-gateway-sql/pkg/database/external"

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
	logger        zerolog.Logger
}

func NewSQLBatchQueryUsecase(sqlQueryRepo *repository.SQLQueryRepo, batchStatRepo *repository.BatchStatRepo, blockRepo *repository.BlockRepo, config *config.Config, logger zerolog.Logger) *SQLBatchQueryUsecase {
	return &SQLBatchQueryUsecase{
		sqlQueryRepo:  sqlQueryRepo,
		batchStatRepo: batchStatRepo,
		blockRepo:     blockRepo,
		config:        config,
		logger:        logger.With().Str("layer", "usecase").Str("component", "sqlbatchquery").Logger(),
	}
}

func (squ *SQLBatchQueryUsecase) ExecuteBatch(ctx context.Context, sqlbatchquery *dto.SQLBatchQueryInput) error {
	target, cfgdb, err := confighelper.GetTargetAndDatabase(squ.config, sqlbatchquery.TargetName)
	if err != nil {
		squ.logger.Error().Err(err).Msg("unable to get target and database from config")
		return err
	}

	if !target.Multi {
		squ.logger.Error().Err(err).Msg(errBatchModeNotActivated.Error())
		return errBatchModeNotActivated
	}

	batchStat, err := squ.batchStatRepo.Create(ctx, target.Name)
	if err != nil {
		squ.logger.Error().Err(err).Str("target", target.Name).Msg("failed to create a batch")
		return err
	}

	blockCh, errCh := csvstream.ReadCSVInBlock(sqlbatchquery.File, target.BufferSize)

	var wg sync.WaitGroup
	for block := range blockCh {
		wg.Add(1)
		go func(b *csvstream.Block) {
			defer wg.Done()
			squ.processBlock(ctx, &dto.BlockDataInput{
				BSInput: batchStat,
				BLInput: b,
				TGInput: target,
				DBInput: cfgdb,
			})
		}(block)
	}

	if err := <-errCh; err != nil {
		wg.Wait()
		return squ.finalizeWithError(ctx, batchStat, err)
	}

	wg.Wait()
	return squ.batchStatRepo.UpdateLastCompleted(ctx, batchStat)
}

func (squ *SQLBatchQueryUsecase) finalizeWithError(ctx context.Context, batchStat *domain.BatchStat, cause error) error {
	if err := squ.batchStatRepo.UpdateLastCompleted(ctx, batchStat); err != nil {
		squ.logger.Error().Err(err).Msg("failed to complete a batch")
		return err
	}

	return cause
}

func (squ *SQLBatchQueryUsecase) processBlock(ctx context.Context, input *dto.BlockDataInput) {
	block, err := squ.initBlock(ctx, input)
	if err != nil {
		squ.logger.Error().Err(err).Msg("failed to initialize a block")
		return
	}

	cnx, err := external.NewDatabase(*input.DBInput)
	if err != nil {
		squ.logger.Error().Err(err).Msg("failed to open database connection")
		return
	}
	squ.sqlQueryRepo.SetDB(cnx)
	defer squ.sqlQueryRepo.CloseDB()

	batchFields := strings.Split(input.TGInput.BatchFields, ";")
	batches := csvmapper.ChunkLines(input.BLInput.Lines, input.TGInput.BatchSize)

	var wg sync.WaitGroup
	for i, batch := range batches {
		wg.Add(1)
		go func(idx int, lines [][]string) {
			defer wg.Done()
			squ.processBatch(ctx, block, input, idx, lines, batchFields)
		}(i, batch)
	}
	wg.Wait()
}

func (squ *SQLBatchQueryUsecase) initBlock(ctx context.Context, input *dto.BlockDataInput) (*domain.Block, error) {
	block := domain.NewBlock(input.BLInput.StartLine, input.BLInput.EndLine)
	block, err := squ.batchStatRepo.AddBlockToBatchStat(ctx, input.BSInput, block)
	if err != nil {
		squ.logger.Error().Err(err).Msg("failed to register block")
		return nil, err
	}

	squ.logger.Info().Msg("init block completed")
	return block, nil
}

func (squ *SQLBatchQueryUsecase) processBatch(ctx context.Context, block *domain.Block, input *dto.BlockDataInput, idx int, lines [][]string, batchFields []string) {
	records, err := csvmapper.MapBatchLines(lines, batchFields)
	if err != nil {
		squ.logger.Error().Err(err).Msg("failed to map batch lines")
		return
	}

	batchSize := input.TGInput.BatchSize
	start, end := idx*batchSize, min(idx*batchSize+len(lines), len(input.BLInput.Lines))

	if execErr := squ.sqlQueryRepo.ExecuteBatch(ctx, input.TGInput.SqlQuery, records); execErr != nil {
		squ.logger.Error().Err(err).Msg("failed to execute batch")
		if err := squ.blockRepo.Update(ctx, block, domain.NewFailureRange(start, end), false); err != nil {
			squ.logger.Error().Err(err).Msg("failed to update block with failure")
		}
		return
	}

	if err := squ.blockRepo.Update(ctx, block, nil, true); err != nil {
		squ.logger.Error().Err(err).Msg("failed to update block on success")
	}
}
