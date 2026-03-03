package csvmapper

import "fmt"

// ChunkLines cuts a slice of csv lines into sub-slices of fixed size.
func ChunkLines(lines [][]string, size int) [][][]string {
	if size <= 0 {
		return nil
	}
	chunks := make([][][]string, 0, (len(lines)+size-1)/size)
	for size < len(lines) {
		lines, chunks = lines[size:], append(chunks, lines[:size])
	}
	return append(chunks, lines)
}

// MapBatchLines maps each CSV line to a key/value record
func MapBatchLines(lines [][]string, fields []string) ([]map[string]any, error) {
	records := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		rec, err := mapBatchFieldToValueLine(fields, line)
		if err != nil {
			return nil, fmt.Errorf("unable to map batch field to value line : %w", err)
		}
		records = append(records, rec)
	}
	return records, nil
}

// mapBatchFieldToValueLine allow to match csv field to his value
func mapBatchFieldToValueLine(fields []string, values []string) (map[string]any, error) {
	if len(fields) != len(values) {
		return nil, fmt.Errorf("bad mapping fields and file column")
	}

	result := make(map[string]any, 0)

	for index, field := range fields {
		result[field] = values[index]
	}

	return result, nil
}
