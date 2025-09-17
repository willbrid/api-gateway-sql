package sqlqueryhelper

import "regexp"

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
