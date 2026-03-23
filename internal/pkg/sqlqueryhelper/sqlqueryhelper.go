package sqlqueryhelper

import (
	"regexp"
	"strings"
)

// TransformQuery used to parse query from config target
func TransformQuery(sqlQuery string, params map[string]any) (string, []any) {
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

// IsSelectQuery used to detect whether a string query is a SELECT or no
func IsSelectQuery(query string) bool {
	trimmed := strings.TrimSpace(query)
	if len(trimmed) < 6 {
		return false
	}

	return strings.EqualFold(trimmed[:6], "SELECT")
}
