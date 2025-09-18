package domain

type SQLQueryOutput struct {
	Rows         []map[string]any
	AffectedRows int64
	DurationMs   int64
}

type SQLQueryInput struct {
	TargetName string
	PostParams map[string]any
}
