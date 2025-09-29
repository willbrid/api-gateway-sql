package repository

import (
	"context"
	"strings"

	"gorm.io/gorm"
)

type SQLInitDatabaseRepo struct {
	db *gorm.DB
}

func NewSQLInitDatabaseRepo() *SQLInitDatabaseRepo {
	return &SQLInitDatabaseRepo{}
}

func (i *SQLInitDatabaseRepo) SetDB(db *gorm.DB) {
	i.db = db
}

func (i *SQLInitDatabaseRepo) CloseDB() {
	if cnx, err := i.db.DB(); err == nil {
		cnx.Close()
	}
}

func (i *SQLInitDatabaseRepo) ExecuteInit(ctx context.Context, sqlQueries []string) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, sqlQuery := range sqlQueries {
			query := strings.TrimSpace(sqlQuery)
			if query != "" {
				if err := tx.Exec(query).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}
