package gdct

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// InitSqliteConnection initializes SQLite database connection
func InitSqliteConnection(dbType string, cfg DBConfig) (*DataBaseConnector, error) {
	// For SQLite, the Database field should contain the file path
	db, err := sql.Open(dbType, cfg.Database)
	if err != nil {
		return nil, fmt.Errorf("sqlite open connection error: %w", err)
	}

	// Apply default configurations for SQLite
	cfg = decideDefaultConfigs(cfg, Sqlite)

	// Set connection pool settings
	if cfg.MaxOpenConns != nil {
		db.SetMaxOpenConns(*cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != nil {
		db.SetMaxIdleConns(*cfg.MaxIdleConns)
	}
	if cfg.MaxLifeTime != nil {
		db.SetConnMaxLifetime(*cfg.MaxLifeTime)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("sqlite ping error: %w", err)
	}

	connect := &DataBaseConnector{db}
	return connect, nil
}

// SqCheckConnection checks SQLite database connection
func (connect *DataBaseConnector) SqCheckConnection() error {
	if err := connect.Ping(); err != nil {
		return fmt.Errorf("sqlite ping error: %w", err)
	}
	return nil
}

// SqCreateTable creates tables using transaction for SQLite
func (connect *DataBaseConnector) SqCreateTable(queryList []string) error {
	ctx := context.Background()

	tx, err := connect.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	for _, queryString := range queryList {
		if _, execErr := tx.ExecContext(ctx, queryString); execErr != nil {
			err = fmt.Errorf("exec transaction context error: %w", execErr)
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction error: %w", err)
	}

	return nil
}

// SqSelectMultiple queries multiple rows from SQLite
func (connect *DataBaseConnector) SqSelectMultiple(queryString string, args ...string) (*sql.Rows, error) {
	arguments := convertArgs(args)

	result, err := connect.Query(queryString, arguments...)
	if err != nil {
		return nil, fmt.Errorf("query select multiple rows error: %w", err)
	}

	return result, nil
}

// SqSelectSingle queries single row from SQLite
func (connect *DataBaseConnector) SqSelectSingle(queryString string, args ...string) (*sql.Row, error) {
	arguments := convertArgs(args)

	result := connect.QueryRow(queryString, arguments...)
	if result.Err() != nil {
		return nil, fmt.Errorf("query single row error: %w", result.Err())
	}

	return result, nil
}

// SqInsertQuery inserts data into SQLite
func (connect *DataBaseConnector) SqInsertQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	insertResult, err := connect.Exec(queryString, arguments...)
	if err != nil {
		return nil, fmt.Errorf("exec insert query error: %w", err)
	}

	return insertResult, nil
}

// SqUpdateQuery updates data in SQLite
func (connect *DataBaseConnector) SqUpdateQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	updateResult, err := connect.Exec(queryString, arguments...)
	if err != nil {
		return nil, fmt.Errorf("exec update query error: %w", err)
	}

	return updateResult, nil
}

// SqDeleteQuery deletes data from SQLite
func (connect *DataBaseConnector) SqDeleteQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	delResult, err := connect.Exec(queryString, arguments...)
	if err != nil {
		return nil, fmt.Errorf("exec delete query error: %w", err)
	}

	return delResult, nil
}

// SqInsertMultiple inserts multiple records with transaction
func (connect *DataBaseConnector) SqInsertMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, err := connect.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			err = fmt.Errorf("prepare statement error: %w", prepareErr)
			return nil, err
		}

		txResult, execErr := stmt.ExecContext(ctx, query.Params...)
		stmt.Close()

		if execErr != nil {
			err = fmt.Errorf("exec prepared statement error: %w", execErr)
			return nil, err
		}

		txResultList = append(txResultList, txResult)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}

	return txResultList, nil
}

// SqUpdateMultiple updates multiple records with transaction
func (connect *DataBaseConnector) SqUpdateMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, err := connect.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			err = fmt.Errorf("prepare statement error: %w", prepareErr)
			return nil, err
		}

		txResult, execErr := stmt.ExecContext(ctx, query.Params...)
		stmt.Close()

		if execErr != nil {
			err = fmt.Errorf("exec prepared statement error: %w", execErr)
			return nil, err
		}

		txResultList = append(txResultList, txResult)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}

	return txResultList, nil
}

// SqDeleteMultiple deletes multiple records with transaction
func (connect *DataBaseConnector) SqDeleteMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
	ctx := context.Background()

	tx, err := connect.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction error: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		}
	}()

	var txResultList []sql.Result

	for _, query := range queryList {
		stmt, prepareErr := tx.PrepareContext(ctx, query.Query)
		if prepareErr != nil {
			err = fmt.Errorf("prepare statement error: %w", prepareErr)
			return nil, err
		}

		txResult, execErr := stmt.ExecContext(ctx, query.Params...)
		stmt.Close()

		if execErr != nil {
			err = fmt.Errorf("exec prepared statement error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("commit transaction error: %w", err)
	}

	return txResultList, nil
}

// Additional SQLite-specific methods

// SqEnableWAL enables Write-Ahead Logging for better concurrency
func (connect *DataBaseConnector) SqEnableWAL() error {
	_, err := connect.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("failed to enable WAL mode: %w", err)
	}
	return nil
}

// SqEnableForeignKeys enables foreign key constraints
func (connect *DataBaseConnector) SqEnableForeignKeys() error {
	_, err := connect.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		return fmt.Errorf("failed to enable foreign keys: %w", err)
	}
	return nil
}

// SqGetVersion returns SQLite version
func (connect *DataBaseConnector) SqGetVersion() (string, error) {
	var version string
	err := connect.QueryRow("SELECT sqlite_version()").Scan(&version)
	if err != nil {
		return "", fmt.Errorf("failed to get SQLite version: %w", err)
	}
	return version, nil
}

// SqVacuum performs VACUUM operation to reclaim space
func (connect *DataBaseConnector) SqVacuum() error {
	_, err := connect.Exec("VACUUM")
	if err != nil {
		return fmt.Errorf("failed to vacuum database: %w", err)
	}
	return nil
}

// SqAnalyze performs ANALYZE operation to update statistics
func (connect *DataBaseConnector) SqAnalyze() error {
	_, err := connect.Exec("ANALYZE")
	if err != nil {
		return fmt.Errorf("failed to analyze database: %w", err)
	}
	return nil
}
