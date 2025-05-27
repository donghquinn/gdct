package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Initiate Mariadb Connection
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
		return nil, fmt.Errorf("postgres open connection error: %w", err)
	}

	cfg = decideDefaultConfigs(cfg, MariaDB)

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
func (connect *DataBaseConnector) MrCheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.Ping()

	if pingErr != nil {
		return fmt.Errorf("postgres ping error: %w", pingErr)
	}

	defer connect.Close()

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

/*
Query Multiple Rows

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Multiple Row Result
*/
func (connect *DataBaseConnector) MrSelectMultiple(queryString string, args ...string) (*sql.Rows, error) {
	arguments := convertArgs(args)

	result, err := connect.Query(queryString, arguments...)

	defer connect.Close()

	if err != nil {
		return nil, fmt.Errorf("query select multiple rows error: %w", err)
	}

	return result, nil
}

/*
Query Single Row

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Single Row Result
*/
func (connect *DataBaseConnector) MrSelectSingle(queryString string, args ...string) (*sql.Row, error) {
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
@Return: Insert ID
*/
func (connect *DataBaseConnector) MrInsertQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	insertResult, insertErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if insertErr != nil {
		return nil, fmt.Errorf("exec insert query error: %w", insertErr)
	}

	return insertResult, nil
}

/*
Update Single Data

@ queryString: Query String with prepared statement
@ args: Query Parameters
@ Return: Affected Rows
*/
func (connect *DataBaseConnector) MrUpdateQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	updateResult, updateErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if updateErr != nil {
		return nil, fmt.Errorf("exec update query error: %w", updateErr)
	}

	return updateResult, nil
}

/*
Delete Single Data

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Affected Rows
*/
func (connect *DataBaseConnector) MrDeleteQuery(queryString string, args ...string) (sql.Result, error) {
	arguments := convertArgs(args)

	delResult, delErr := connect.Exec(queryString, arguments...)

	defer connect.Close()

	if delErr != nil {
		return nil, fmt.Errorf("exec delete query error: %w", delErr)
	}

	return delResult, nil
}

/*
INSERT Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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

/*
UPDATE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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

/*
DELETE Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
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
