package csvstream

import (
	"github.com/willbrid/api-gateway-sql/pkg/logger"

	"encoding/csv"
	"io"
	"mime/multipart"
	"strings"
)

type Block struct {
	StartLine int
	EndLine   int
	Lines     [][]string
}

type ICSVStream interface {
	ReadCSVInBlock(file multipart.File, blockSize int) (chan *Block, chan error)
}

type CSVStream struct {
	iLogger logger.ILogger
}

func NewCSVStream(iLogger logger.ILogger) *CSVStream {
	return &CSVStream{iLogger}
}

func (s *CSVStream) ReadCSVInBlock(file multipart.File, blockSize int) (chan *Block, chan error) {
	blockChannel := make(chan *Block)
	errorChannel := make(chan error)

	go func() {
		numLine := 0
		reader := csv.NewReader(file)

		defer close(blockChannel)
		defer close(errorChannel)

		for {
			block := &Block{
				StartLine: blockSize*numLine + 1,
				Lines:     make([][]string, 0, blockSize),
			}

			for i := range blockSize {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
					s.iLogger.Error("failed to read a line %v - start of the block: %v - error: %s", i, block.StartLine, err.Error())
					errorChannel <- err
					return
				}
				block.Lines = append(block.Lines, strings.Split(record[0], ";"))
			}

			if len(block.Lines) == 0 {
				break
			}

			block.EndLine = block.StartLine + len(block.Lines) - 1
			blockChannel <- block
			numLine++
		}
	}()

	return blockChannel, errorChannel
}
