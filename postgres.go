package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// mockResult implements sql.Result interface for RETURNING queries
type mockResult struct{}

func (m *mockResult) LastInsertId() (int64, error) {
	return 0, fmt.Errorf("LastInsertId is not supported when using RETURNING clause")
}

func (m *mockResult) RowsAffected() (int64, error) {
	return 1, nil // Assume 1 row was affected for RETURNING queries
}

// InitPostgresConnection initializes a PostgreSQL database connection.
func InitPostgresConnection(dbType string, cfg DBConfig) (*DataBaseConnector, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		*cfg.SslMode,
	)

	db, err := sql.Open(dbType, dbUrl)

	if err != nil {
		return nil, fmt.Errorf("postgres open connection error: %w", err)
	}

	cfg = decideDefaultConfigs(cfg, PostgreSQL)

	if cfg.MaxOpenConns != nil {
		db.SetMaxOpenConns(*cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != nil {
		db.SetMaxIdleConns(*cfg.MaxIdleConns)
	}
	if cfg.MaxLifeTime != nil {
		db.SetConnMaxLifetime(*cfg.MaxLifeTime)
	}

	connect := &DataBaseConnector{DB: db, dbType: PostgreSQL}

	return connect, nil
}

// PgCheckConnection checks the PostgreSQL database connection.
func (connect *DataBaseConnector) PgCheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.Ping()

	if pingErr != nil {
		return fmt.Errorf("postgres ping error: %w", pingErr)
	}

	return nil
}

func (connect *DataBaseConnector) PgCreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return fmt.Errorf("bigin transaction error: %w", txErr)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && txErr != sql.ErrTxDone {
			log.Printf("[CREATE_TABLE] Transaction rollback error: %v", txErr)
		}
	}()

	for _, queryString := range queryList {
		if _, execErr := tx.ExecContext(ctx, queryString); execErr != nil {
			return fmt.Errorf("exec transaction context error: %w", execErr)
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return nil
}

// PgSelectMultiple executes a query that returns multiple rows.
// Note: Caller is responsible for closing the returned *sql.Rows.
func (connect *DataBaseConnector) PgSelectMultiple(queryString string, args []interface{}) (*sql.Rows, error) {
	result, err := connect.Query(queryString, args...)

	if err != nil {
		return nil, fmt.Errorf("query select multiple rows error: %w", err)
	}

	return result, nil
}

// PgSelectSingle executes a query that returns at most one row.
func (connect *DataBaseConnector) PgSelectSingle(queryString string, args []interface{}) (*sql.Row, error) {
	result := connect.QueryRow(queryString, args...)

	if result.Err() != nil {
		return nil, fmt.Errorf("query single row error: %w", result.Err())
	}

	return result, nil
}

// PgInsertQuery executes an INSERT query with optional RETURNING clause.
func (connect *DataBaseConnector) PgInsertQuery(queryString string, returns []interface{}, args []interface{}) (sql.Result, error) {
	// If returns is provided and not empty, we need to handle RETURNING clause
	if len(returns) > 0 {
		// Use QueryRow for RETURNING clause to scan the returned values
		row := connect.QueryRow(queryString, args...)
		if err := row.Scan(returns...); err != nil {
			return nil, fmt.Errorf("scan returning values error: %w", err)
		}
		
		// Create a mock Result since we can't get the actual sql.Result from QueryRow
		// This is a limitation when using RETURNING - you get the returned values but lose Result info
		return &mockResult{}, nil
	}

	// No RETURNING clause, use normal Exec
	insertResult, queryErr := connect.Exec(queryString, args...)
	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return insertResult, nil
}

// PgUpdateQuery executes an UPDATE query.
func (connect *DataBaseConnector) PgUpdateQuery(queryString string, args []interface{}) (sql.Result, error) {
	updateResult, queryErr := connect.Exec(queryString, args...)

	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return updateResult, nil
}

// PgDeleteQuery executes a DELETE query.
func (connect *DataBaseConnector) PgDeleteQuery(queryString string, args []interface{}) (sql.Result, error) {
	deleteResult, queryErr := connect.Exec(queryString, args...)

	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return deleteResult, nil
}

// PgInsertMultiple executes multiple INSERT queries within a transaction.
func (connect *DataBaseConnector) PgInsertMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("begin transaction error: %w", txErr)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		// Prepared statement
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			return nil, fmt.Errorf("prepare statement error: %w", prepareErr)
		}

		// PreparedStatement
		txResult, execErr := stmt.ExecContext(ctx, query.Params...)

		// Statement
		stmt.Close()

		if execErr != nil {
			return nil, fmt.Errorf("exec prepared statement error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}

// PgUpdateMultiple executes multiple UPDATE queries within a transaction.
func (connect *DataBaseConnector) PgUpdateMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("begin transaction error: %w", txErr)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		// Prepared statement
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			return nil, fmt.Errorf("prepare statement error: %w", prepareErr)
		}

		// PreparedStatement
		txResult, execErr := stmt.ExecContext(ctx, query.Params...)

		// Statement
		stmt.Close()

		if execErr != nil {
			return nil, fmt.Errorf("exec prepared statement error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}

// PgDeleteMultiple executes multiple DELETE queries within a transaction.
func (connect *DataBaseConnector) PgDeleteMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("begin transaction error: %w", txErr)
	}

	defer func() {
		_ = tx.Rollback()
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		// Prepared statement
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			return nil, fmt.Errorf("prepare statement error: %w", prepareErr)
		}

		// PreparedStatement
		txResult, execErr := stmt.ExecContext(ctx, query.Params...)

		// Statement
		stmt.Close()

		if execErr != nil {
			return nil, fmt.Errorf("exec prepared statement error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}
