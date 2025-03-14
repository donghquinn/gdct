package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB 연결 인스턴스
func InitPostgresConnection(cfg DBConfig) (*DataBaseConnector, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.SslMode,
	)

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		return nil, fmt.Errorf("postgres open connection error: %w", err)
	}

	cfg = decideDefaultConfigs(cfg)

	db.SetConnMaxLifetime(cfg.MaxLifeTime)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxIdleConns)

	connect := &DataBaseConnector{db}

	return connect, nil
}

/*
Check Connection
*/
func (connect *DataBaseConnector) PgCheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.Ping()

	if pingErr != nil {
		return fmt.Errorf("postgres ping error: %w", pingErr)
	}

	defer connect.Close()

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

/*
Query Multiple Rows

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Multiple Row Result
*/
func (connect *DataBaseConnector) PgSelectMultiple(queryString string, args ...string) (*sql.Rows, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result, err := connect.Query(queryString, arguments...)

	if err != nil {
		return nil, fmt.Errorf("query select multiple rows error: %w", err)
	}

	defer connect.Close()

	return result, nil
}

/*
Query Single Row

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Single Row Result
*/
func (connect *DataBaseConnector) PgSelectSingle(queryString string, args ...string) (*sql.Row, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result := connect.QueryRow(queryString, arguments...)

	defer connect.Close()

	if result.Err() != nil {
		return nil, fmt.Errorf("query single row error: %w", result.Err())
	}

	return result, nil
}

/*
Insert Single Data

@queryString: Query String with prepared statement
@returns: Return Value by RETURNING <Column_name>;
@args: Query Parameters
*/
func (connect *DataBaseConnector) PgInsertQuery(queryString string, returns []interface{}, args ...string) error {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	queryResult := connect.QueryRow(queryString, arguments...)

	defer connect.Close()

	if returns != nil {
		// Insert ID
		if scanErr := queryResult.Scan(returns...); scanErr != nil {
			return fmt.Errorf("scan returning values by insert: %w", scanErr)
		}
	}

	return nil
}

/*
Update Single Data

@ queryString: Query String with prepared statement
@ args: Query Parameters
@ Return: Affected Rows
*/
func (connect *DataBaseConnector) PgUpdateQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	updateResult, queryErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if queryErr != nil {
		return -999, fmt.Errorf("exec query error: %w", queryErr)
	}

	// Insert ID
	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		return -999, fmt.Errorf("get affected row error: %w", afftedRowErr)
	}

	return affectedRow, nil
}

/*
DELETE Single Data

@ queryString: Query String with prepared statement
@ args: Query Parameters
@ Return: Affected Rows
*/
func (connect *DataBaseConnector) PgDeleteQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	updateResult, queryErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if queryErr != nil {
		return -999, fmt.Errorf("exec query error: %w", queryErr)
	}

	// Insert ID
	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		return -999, fmt.Errorf("get affected row error: %w", afftedRowErr)
	}

	return affectedRow, nil
}

/*
INSERT Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
func (connect *DataBaseConnector) PgInsertMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("bigin transaction error: %w", txErr)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && txErr != sql.ErrTxDone {
			log.Printf("[INSERT_MULTIPLE] Transaction rollback error: %v", txErr)
		}
	}()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			return nil, fmt.Errorf("exec insert multiple transaction context error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}

/*
UPDATE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
func (connect *DataBaseConnector) PgUpdateMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("bigin transaction error: %w", txErr)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && txErr != sql.ErrTxDone {
			log.Printf("[UPDATE_MULTIPLE] Transaction rollback error: %v", txErr)
		}
	}()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			return nil, fmt.Errorf("exec update multiple transaction context error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}

/*
DELETE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
func (connect *DataBaseConnector) PgDeleteMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, fmt.Errorf("bigin transaction error: %w", txErr)
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && txErr != sql.ErrTxDone {
			log.Printf("[DELETE_MULTIPLE] Transaction rollback error: %v", txErr)
		}
	}()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			return nil, fmt.Errorf("exec delete multiple transaction context error: %w", execErr)
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}
