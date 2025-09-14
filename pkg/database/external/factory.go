package external

import (
	"api-gateway-sql/config"
)

func NewExecutor(db config.Database) (QueryExecutor, error) {
	dbType := db.Type

	switch dbType {
	case "mariadb":
		return NewMariadbExecutor(db)
	case "mysql":
		return NewMySQLExecutor(db)
	case "postgres":
		return NewPostgresExecutor(db)
	case "sqlserver":
		return NewSQLServerExecutor(db)
	case "sqlite":
		return NewSqliteExecutor(db)
	default:
		return nil, nil
	}
}
