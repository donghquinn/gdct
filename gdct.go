package gdct

import (
	"database/sql"
	"fmt"
	"time"
)

// DBConfig holds database connection configuration.
type DBConfig struct {
	UserName     string         // Database username
	Password     string         // Database password
	Host         string         // Database host
	Port         int            // Database port
	Database     string         // Database name or file path for SQLite
	SslMode      *string        // SSL mode for PostgreSQL
	MaxLifeTime  *time.Duration // Maximum connection lifetime
	MaxIdleConns *int           // Maximum idle connections
	MaxOpenConns *int           // Maximum open connections
}

// DataBaseConnector wraps sql.DB with additional functionality.
type DataBaseConnector struct {
	*sql.DB
	dbType DBType // Store database type for query building
}

// PreparedQuery represents a prepared SQL query with parameters.
type PreparedQuery struct {
	Query  string        // SQL query string
	Params []interface{} // Query parameters
}

// InitConnection creates a new database connection based on the database type.
func InitConnection(dbType DBType, cfg DBConfig) (*DataBaseConnector, error) {
	switch dbType {
	case MariaDB:
		return InitMariadbConnection("mysql", cfg)
	case Mysql:
		return InitMariadbConnection("mysql", cfg)
	case PostgreSQL:
		return InitPostgresConnection("postgres", cfg)
	case Sqlite:
		return InitSqliteConnection("sqlite3", cfg)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}

// QueryBuilderRows executes a query that returns multiple rows.
// Note: Caller is responsible for closing the returned *sql.Rows.
func (connect *DataBaseConnector) QueryBuilderRows(queryString string, args []interface{}) (*sql.Rows, error) {
	result, err := connect.Query(queryString, args...)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	return result, nil
}

// QueryBuilderOneRow executes a query that returns at most one row.
func (connect *DataBaseConnector) QueryBuilderOneRow(queryString string, args []interface{}) *sql.Row {
	return connect.QueryRow(queryString, args...)
}

// QueryBuilderInsert executes an INSERT query.
func (connect *DataBaseConnector) QueryBuilderInsert(queryString string, args []interface{}) (sql.Result, error) {
	result, err := connect.Exec(queryString, args...)
	if err != nil {
		return nil, fmt.Errorf("insert execution failed: %w", err)
	}
	return result, nil
}

// QueryBuilderUpdate executes an UPDATE query.
func (connect *DataBaseConnector) QueryBuilderUpdate(queryString string, args []interface{}) (sql.Result, error) {
	result, err := connect.Exec(queryString, args...)
	if err != nil {
		return nil, fmt.Errorf("update execution failed: %w", err)
	}
	return result, nil
}
