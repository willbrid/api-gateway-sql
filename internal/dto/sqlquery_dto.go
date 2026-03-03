package dto

import "mime/multipart"

type SQLQueryOutput struct {
	Rows         []map[string]any
	AffectedRows int64
	DurationMs   int64
}

type SQLQueryInput struct {
	TargetName string
	PostParams map[string]any
}

type SQLBatchQueryInput struct {
	TargetName string
	File       multipart.File
}

type SQLInitDatabaseInput struct {
	Datasource     string
	SQLFileContent string
}
