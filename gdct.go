package gdct

import (
	"database/sql"
	"fmt"
	"time"
)

type DBConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         int
	Database     string
	SslMode      *string
	MaxLifeTime  *time.Duration
	MaxIdleConns *int
	MaxOpenConns *int
}

type DataBaseConnector struct {
	*sql.DB
}

type PreparedQuery struct {
	Query  string
	Params []interface{}
}

func InitConnection(dbType DBType, cfg DBConfig) (*DataBaseConnector, error) {
	switch dbType {
	case MariaDB:
		return InitMariadbConnection("mysql", cfg)
	case Mysql:
		return InitMariadbConnection("mysql", cfg)
	case PostgreSQL:
		return InitPostgresConnection("postgres", cfg)
	case Sqlite:
		return InitSqliteConnection("sqlite", cfg)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}

func (connect *DataBaseConnector) QueryBuilderRows(queryString string, args []interface{}) (*sql.Rows, error) {
	result, err := connect.Query(queryString, args...)

	if err != nil {
		return nil, err
	}

	defer connect.Close()

	return result, nil
}

func (connect *DataBaseConnector) QueryBuilderOneRow(queryString string, args []interface{}) (*sql.Row, error) {
	result := connect.QueryRow(queryString, args...)

	if result.Err() != nil {
		return nil, result.Err()
	}

	defer connect.Close()

	return result, nil
}

func (connect *DataBaseConnector) QueryBuilderInsert(queryString string, args []interface{}) (sql.Result, error) {
	updateResult, err := connect.Exec(queryString, args...)

	if err != nil {
		return nil, err
	}

	defer connect.Close()

	return updateResult, nil
}

func (connect *DataBaseConnector) QueryBuilderUpdate(queryString string, args []interface{}) (sql.Result, error) {
	updateResult, err := connect.Exec(queryString, args...)

	if err != nil {
		return nil, err
	}

	defer connect.Close()

	return updateResult, nil
}
