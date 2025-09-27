package mapperfieldshelper

import "fmt"

func MapBatchFieldToValueLine(fields []string, values []string) (map[string]any, error) {
	if len(fields) != len(values) {
		return nil, fmt.Errorf("bad mapping fields and file column")
	}

	var result map[string]any = make(map[string]any, 0)

	for index, field := range fields {
		result[field] = values[index]
	}

	return result, nil
}
