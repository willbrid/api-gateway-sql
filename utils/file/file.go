package file

import (
	"encoding/csv"
	"io"
	"mime/multipart"
	"strings"
)

type Buffer struct {
	StartLine int
	EndLine   int
	Lines     [][]string
}

func ReadCSVInBuffer(file multipart.File, bufferSize int) ([]Buffer, error) {
	var (
		buffers []Buffer
	)

	reader := csv.NewReader(file)
	numLine := 0

	for {
		buffer := Buffer{
			StartLine: bufferSize*numLine + 1,
			Lines:     make([][]string, 0, bufferSize),
		}

		for i := 0; i < bufferSize; i++ {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err
			}
			buffer.Lines = append(buffer.Lines, strings.Split(record[0], ";"))
		}

		if len(buffer.Lines) == 0 {
			break
		}

		buffer.EndLine = buffer.StartLine + len(buffer.Lines) - 1
		buffers = append(buffers, buffer)
		numLine++
	}

	return buffers, nil
}
