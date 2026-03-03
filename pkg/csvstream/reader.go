package csvstream

import (
	"encoding/csv"
	"fmt"
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
					errorChannel <- fmt.Errorf("failed to read a line %v - start of the block: %v - error: %w", i, block.StartLine, err)
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
