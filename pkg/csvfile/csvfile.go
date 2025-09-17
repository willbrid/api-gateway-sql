package csvfile

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

func ReadCSVInBlock(file multipart.File, blockSize int) ([]Block, error) {
	var (
		blocks []Block
	)

	reader := csv.NewReader(file)
	numLine := 0

	for {
		block := Block{
			StartLine: blockSize*numLine + 1,
			Lines:     make([][]string, 0, blockSize),
		}

		for i := 0; i < blockSize; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			block.Lines = append(block.Lines, strings.Split(record[0], ";"))
		}

		if len(block.Lines) == 0 {
			break
		}

		block.EndLine = block.StartLine + len(block.Lines) - 1
		blocks = append(blocks, block)
		numLine++
	}

	return blocks, nil
}
