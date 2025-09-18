package external

import (
	"api-gateway-sql/config"

	"errors"

	"gorm.io/gorm"
)

const (
	mariadbType   string = "mariadb"
	mysqlType     string = "mysql"
	postgresType  string = "postgres"
	sqlserverType string = "sqlserver"
	sqliteType    string = "sqlite"
)

var (
	errUnknownDatabaseType error = errors.New("unknown database type")
)

type IDatabase interface {
	Connect(dbConfig config.Database) (*gorm.DB, error)
}

func NewDatabase(db config.Database) (*gorm.DB, error) {
	var dbInstance IDatabase
	dbType := db.Type

	switch dbType {
	case mariadbType:
		dbInstance = &MariadbDatabase{}
	case mysqlType:
		dbInstance = &MySQLDatabase{}
	case postgresType:
		dbInstance = &PostgresDatabase{}
	case sqliteType:
		dbInstance = &SqliteDatabase{}
	case sqlserverType:
		dbInstance = &SQLServerDatabase{}
	default:
		dbInstance = nil
	}

	if dbInstance == nil {
		return nil, errUnknownDatabaseType
	}

	return dbInstance.Connect(db)
}
