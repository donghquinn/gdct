package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB 연결 인스턴스
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

	if cfg.MaxIdleConns != nil {
		db.SetMaxOpenConns(*cfg.MaxIdleConns)
	}
	if cfg.MaxLifeTime != nil {
		db.SetConnMaxLifetime(*cfg.MaxLifeTime)
	}

	if cfg.MaxOpenConns != nil {
		db.SetMaxIdleConns(*cfg.MaxIdleConns)
	}

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
		if _, execErr := tx.ExecContext(ctx, queryString); execErr != nil {
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
	arguments := convertArgs(args)

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
	arguments := convertArgs(args)

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
@args: Query Parameters
*/
func (connect *DataBaseConnector) PgInsertQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	insertResult, queryErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return insertResult, nil
}

/*
Insert Single Data With Returning

@queryString: Query String with prepared statement
@returns: Return Value by RETURNING <Column_name>;
@args: Query Parameters
*/
func (connect *DataBaseConnector) PgInsertQueryReturning(queryString string, returns []interface{}, args ...string) error {
	arguments := convertArgs(args)

	insertResult := connect.QueryRow(queryString, arguments...)

	defer connect.Close()

	if returns != nil {
		if scanErr := insertResult.Scan(returns...); scanErr != nil {
			return fmt.Errorf("exec insert query with returning error: %v", scanErr)
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
func (connect *DataBaseConnector) PgUpdateQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	updateResult, queryErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return updateResult, nil
}

/*
DELETE Single Data

@ queryString: Query String with prepared statement
@ args: Query Parameters
@ Return: Affected Rows
*/
func (connect *DataBaseConnector) PgDeleteQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	deleteResult, queryErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if queryErr != nil {
		return nil, fmt.Errorf("exec query error: %w", queryErr)
	}

	return deleteResult, nil
}

/*
INSERT Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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

/*
UPDATE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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

/*
DELETE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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
