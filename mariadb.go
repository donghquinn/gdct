package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// InitMariadbConnection initializes a MariaDB/MySQL database connection.
func InitMariadbConnection(dbType string, cfg DBConfig) (*DataBaseConnector, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open(dbType, dbUrl)

	if err != nil {
		return nil, fmt.Errorf("mariadb open connection error: %w", err)
	}

	cfg = decideDefaultConfigs(cfg, MariaDB)

	if cfg.MaxOpenConns != nil {
		db.SetMaxOpenConns(*cfg.MaxOpenConns)
	}
	if cfg.MaxIdleConns != nil {
		db.SetMaxIdleConns(*cfg.MaxIdleConns)
	}
	if cfg.MaxLifeTime != nil {
		db.SetConnMaxLifetime(*cfg.MaxLifeTime)
	}

	connect := &DataBaseConnector{DB: db, dbType: MariaDB}

	return connect, nil
}

// MrCheckConnection checks the MariaDB/MySQL database connection.
func (connect *DataBaseConnector) MrCheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.Ping()

	if pingErr != nil {
		return fmt.Errorf("mariadb ping error: %w", pingErr)
	}

	return nil
}

func (connect *DataBaseConnector) MrCreateTable(queryList []string) error {
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
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			return fmt.Errorf("exec transaction context error: %w", execErr)
		}
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return nil
}

// MrSelectMultiple executes a query that returns multiple rows.
// Note: Caller is responsible for closing the returned *sql.Rows.
func (connect *DataBaseConnector) MrSelectMultiple(queryString string, args []interface{}) (*sql.Rows, error) {
	result, err := connect.Query(queryString, args...)

	if err != nil {
		return nil, fmt.Errorf("query select multiple rows error: %w", err)
	}

	return result, nil
}

// MrSelectSingle executes a query that returns at most one row.
func (connect *DataBaseConnector) MrSelectSingle(queryString string, args []interface{}) (*sql.Row, error) {
	result := connect.QueryRow(queryString, args...)

	if result.Err() != nil {
		return nil, fmt.Errorf("query single row error: %w", result.Err())
	}

	return result, nil
}

// MrInsertQuery executes an INSERT query.
func (connect *DataBaseConnector) MrInsertQuery(queryString string, args []interface{}) (sql.Result, error) {
	insertResult, insertErr := connect.Exec(queryString, args...)

	if insertErr != nil {
		return nil, fmt.Errorf("exec insert query error: %w", insertErr)
	}

	return insertResult, nil
}

// MrUpdateQuery executes an UPDATE query.
func (connect *DataBaseConnector) MrUpdateQuery(queryString string, args []interface{}) (sql.Result, error) {
	updateResult, updateErr := connect.Exec(queryString, args...)

	if updateErr != nil {
		return nil, fmt.Errorf("exec update query error: %w", updateErr)
	}

	return updateResult, nil
}

// MrDeleteQuery executes a DELETE query.
func (connect *DataBaseConnector) MrDeleteQuery(queryString string, args []interface{}) (sql.Result, error) {
	delResult, delErr := connect.Exec(queryString, args...)

	if delErr != nil {
		return nil, fmt.Errorf("exec delete query error: %w", delErr)
	}

	return delResult, nil
}

// MrInsertMultiple executes multiple INSERT queries within a transaction.
func (connect *DataBaseConnector) MrInsertMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
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

// MrUpdateMultiple executes multiple UPDATE queries within a transaction.
func (connect *DataBaseConnector) MrUpdateMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
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

// MrDeleteMultiple executes multiple DELETE queries within a transaction.
func (connect *DataBaseConnector) MrDeleteMultiple(queryList []PreparedQuery) ([]sql.Result, error) {
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
