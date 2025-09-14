package external

import "context"

type ExecutionResult struct {
	Rows         []map[string]any
	AffectedRows int64
	DurationMs   int64
}

type QueryExecutor interface {
	Execute(ctx context.Context, query string, params map[string]any) (*ExecutionResult, error)
}
