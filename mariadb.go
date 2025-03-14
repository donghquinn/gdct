package gdct

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Initiate Mariadb Connection
func InitMariadbConnection(cfg DBConfig) (*DataBaseConnector, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("mysql", dbUrl)

	if err != nil {
		log.Printf("[INIT] Start Database Connection Error: %v", err)

		return nil, err
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
func (connect *DataBaseConnector) MrCheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.Ping()

	if pingErr != nil {
		log.Printf("[CHECK] Database Ping Error: %v", pingErr)
		return pingErr
	}

	defer connect.Close()

	return nil
}

func (connect *DataBaseConnector) MrCreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return txErr
	}

	defer func() {
		if txErr := tx.Rollback(); txErr != nil && txErr != sql.ErrTxDone {
			log.Printf("[CREATE_TABLE] Transaction rollback error: %v", txErr)
		}
	}()

	for _, queryString := range queryList {
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			return execErr
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
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result, err := connect.Query(queryString, arguments...)

	defer connect.Close()

	if err != nil {
		return nil, err
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
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result := connect.QueryRow(queryString, arguments...)

	defer connect.Close()

	if result.Err() != nil {
		return nil, result.Err()
	}

	return result, nil
}

/*
Insert Single Data

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Insert ID
*/
func (connect *DataBaseConnector) MrInsertQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	insertResult, insertErr := connect.Exec(queryString, arguments...)

	if insertErr != nil {
		return -99999, insertErr
	}

	defer connect.Close()

	// Insert ID
	insertId, insertIdErr := insertResult.LastInsertId()

	if insertIdErr != nil {
		return -999999, insertIdErr
	}

	return insertId, nil
}

/*
Update Single Data

@ queryString: Query String with prepared statement
@ args: Query Parameters
@ Return: Affected Rows
*/
func (connect *DataBaseConnector) MrUpdateQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	updateResult, updateErr := connect.Exec(queryString, arguments...)

	if updateErr != nil {
		return -99999, updateErr
	}

	defer connect.Close()

	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		return -999999, afftedRowErr
	}

	return affectedRow, nil
}

/*
Delete Single Data

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Affected Rows
*/
func (connect *DataBaseConnector) MrDeleteQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	delResult, delErr := connect.Exec(queryString, arguments...)

	if delErr != nil {
		return -99999, delErr
	}

	defer connect.Close()

	// Insert ID
	affectedRow, afftedRowErr := delResult.RowsAffected()

	if afftedRowErr != nil {
		return -999999, afftedRowErr
	}

	return affectedRow, nil
}

/*
INSERT Multiple Data with DB Transaction

@ queryString: Query String with prepared statement
*/
func (connect *DataBaseConnector) MrInsertMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, txErr
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
			return nil, execErr
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
func (connect *DataBaseConnector) MrUpdateMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, txErr
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
			return nil, execErr
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
func (connect *DataBaseConnector) MrDeleteMultiple(queryList []string) ([]sql.Result, error) {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		return nil, txErr
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
			return nil, execErr
		}

		txResultList = append(txResultList, txResult)
	}

	if commitErr := tx.Commit(); commitErr != nil {
		return nil, fmt.Errorf("commit transaction error: %w", commitErr)
	}

	return txResultList, nil
}
