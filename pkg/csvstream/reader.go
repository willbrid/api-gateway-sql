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

func ReadCSVInBlock(file multipart.File, blockSize int) (chan *Block, chan error) {
	var (
		blockChannel chan *Block = make(chan *Block)
		errorChannel chan error  = make(chan error)
	)

	go func() {
		var numLine int = 0
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
					logger.Error("failed to read a line %v - start of the block: %v - error: %s", i, block.StartLine, err.Error())
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
