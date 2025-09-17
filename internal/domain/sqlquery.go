package domain

type SQLQueryOutput struct {
	Rows         []map[string]any
	AffectedRows int64
	DurationMs   int64
}

type SQLQuery struct {
	targetName string
	postParams map[string]any
}
