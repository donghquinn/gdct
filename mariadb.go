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

func (connect *DataBaseConnector) MrCreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.Begin()

	if txErr != nil {
		log.Printf("[CREATE_TABLE] Begin Transaction Error: %v", txErr)
		return txErr
	}

	defer tx.Rollback()

	for _, queryString := range queryList {
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[CREATE_TABLE] Create Table Querystring Transaction Exec Error: %v", execErr)
			return execErr
		}
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[CREATE_TABLE] Commit Transaction Error: %v", commitErr)
		return commitErr
	}

	return nil
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
func (connect *DataBaseConnector) MrSelectSingle(queryString string, args ...string) (*sql.Row, error) {
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
		log.Printf("[INSERT] Insert Query Err: %v", insertErr)

		return -99999, insertErr
	}

	defer connect.Close()

	// Insert ID
	insertId, insertIdErr := insertResult.LastInsertId()

	if insertIdErr != nil {
		log.Printf("[INSERT] Get Insert ID Error: %v", insertIdErr)

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
		log.Printf("[UPDATE] Update Query Err: %v", updateErr)

		return -99999, updateErr
	}

	defer connect.Close()

	affectedRow, afftedRowErr := updateResult.RowsAffected()

	if afftedRowErr != nil {
		log.Printf("[UPDATE] Get Affected Rows Error: %v", afftedRowErr)

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
		log.Printf("[DELETE] Delete Query Err: %v", delErr)

		return -99999, delErr
	}

	defer connect.Close()

	// Insert ID
	affectedRow, afftedRowErr := delResult.RowsAffected()

	if afftedRowErr != nil {
		log.Printf("[DELETE] Get Affected Numbers Error: %v", afftedRowErr)

		return -999999, afftedRowErr
	}

	return affectedRow, nil
}
