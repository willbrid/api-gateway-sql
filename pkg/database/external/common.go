package external

import (
	"context"
	"regexp"
	"time"

	"gorm.io/gorm"
)

// executeQuery used to execute query with gorm connection
func executeQuery(ctx context.Context, cnx *gorm.DB, query string, params map[string]any) (*ExecutionResult, error) {
	start := time.Now()
	var result []map[string]any
	parsedQuery, parsedParams := transformQuery(query, params)

	if err := cnx.WithContext(ctx).Raw(parsedQuery, parsedParams...).Scan(&result).Error; err == nil {
		return &ExecutionResult{
			Rows:         result,
			AffectedRows: int64(len(result)),
			DurationMs:   time.Since(start).Milliseconds(),
		}, nil
	}

	tx := cnx.Exec(parsedQuery, parsedParams...)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &ExecutionResult{
		Rows:         nil,
		AffectedRows: tx.RowsAffected,
		DurationMs:   time.Since(start).Milliseconds(),
	}, nil
}

// transformQuery used to parse query from config target
func transformQuery(sqlQuery string, params map[string]any) (string, []any) {
	re := regexp.MustCompile(`{{(\w+)}}`)
	matches := re.FindAllStringSubmatch(sqlQuery, -1)

	values := make([]any, 0, len(matches))
	transformedQuery := re.ReplaceAllStringFunc(sqlQuery, func(param string) string {
		paramName := param[2 : len(param)-2]
		if value, exists := params[paramName]; exists {
			values = append(values, value)
			return "?"
		}
		return param
	})

	return transformedQuery, values
}
