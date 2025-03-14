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
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
	)

	db, err := sql.Open("postgres", dbUrl)

	if err != nil {
		log.Printf("[DATABASE] Start Database Connection Error: %v", err)

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
func (connect *DataBaseConnector) PgCheckConnection() error {
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

func (connect *DataBaseConnector) PgCreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		log.Printf("[DATABASE] Begin Transaction Error: %v", txErr)
		return txErr
	}

	defer tx.Rollback()

	for _, queryString := range queryList {
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[DATABASE] Create Table Querystring Transaction Exec Error: %v", execErr)
			return execErr
		}
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[DATABASE] Commit Transaction Error: %v", commitErr)
		return commitErr
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
		log.Printf("[QUERY] Query Error: %v\n", err)

		return nil, err
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

	if result.Err() != nil {
		log.Printf("[QUERY] Query Error: %v\n", result.Err())

		return nil, result.Err()
	}

	defer connect.Close()

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
			log.Printf("[INSERT] Get Insert Result Scan Error: %v", scanErr)
			return scanErr
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
		log.Printf("[UPDATE] Update Query Error: %v", queryErr)
		return -999, queryErr
	}

	// Insert ID
	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		log.Printf("[UPDATE] Get Affected Rows Error: %v", afftedRowErr)

		return -999, afftedRowErr
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
		log.Printf("[DELETE] Delete Query Error: %v", queryErr)
		return -999, queryErr
	}

	// Insert ID
	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		log.Printf("[DELETE] Get Affected Rows Error: %v", afftedRowErr)

		return -999, afftedRowErr
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
		log.Printf("[INSERT_MULTIPLE] Begin Transaction Error: %v", txErr)
		return []sql.Result{}, txErr
	}

	defer tx.Rollback()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[INSERT_MULTIPLE] Insert Querystring Transaction Exec Error: %v", execErr)
			return []sql.Result{}, execErr
		}

		txResultList = append(txResultList, txResult)
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[INSERT_MULTIPLE] Commit Transaction Error: %v", commitErr)
		return []sql.Result{}, commitErr
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
		log.Printf("[UPDATE_MULTIPLE] Begin Transaction Error: %v", txErr)
		return []sql.Result{}, txErr
	}

	defer tx.Rollback()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[UPDATE_MULTIPLE] Update Querystring Transaction Exec Error: %v", execErr)
			return []sql.Result{}, execErr
		}

		txResultList = append(txResultList, txResult)
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[UPDATE_MULTIPLE] Commit Transaction Error: %v", commitErr)
		return []sql.Result{}, commitErr
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
		log.Printf("[DELETE_MULTIPLE] Begin Transaction Error: %v", txErr)
		return []sql.Result{}, txErr
	}

	defer tx.Rollback()

	var txResultList []sql.Result

	for _, queryString := range queryList {
		txResult, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[DELETE_MULTIPLE] Delete Querystring Transaction Exec Error: %v", execErr)
			return []sql.Result{}, execErr
		}

		txResultList = append(txResultList, txResult)
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[DELETE_MULTIPLE] Commit Transaction Error: %v", commitErr)
		return []sql.Result{}, commitErr
	}

	return txResultList, nil
}
