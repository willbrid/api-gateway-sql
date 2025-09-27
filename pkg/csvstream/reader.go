package csvstream

import (
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
		blockChannel chan *Block
		errorChannel chan error
	)

	go func() {
		reader := csv.NewReader(file)
		numLine := 0

		defer close(blockChannel)
		defer close(errorChannel)

		for {
			block := &Block{
				StartLine: blockSize*numLine + 1,
				Lines:     make([][]string, 0, blockSize),
			}

			for i := 0; i < blockSize; i++ {
				record, err := reader.Read()
				if err == io.EOF {
					break
				}
				if err != nil {
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
