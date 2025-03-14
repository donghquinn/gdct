package mariadb

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type MariadbConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         string
	Database     string
	MaxLifeTime  time.Duration // time.Duration 타입을 권장 (예: 60 * time.Second)
	MaxIdleConns int
	MaxOpenConns int
}

type MariaDbInstance struct {
	conn *sql.DB
}

// Initiate Mariadb Connection
func InitMariadbConnection(cfg MariadbConfig) (*MariaDbInstance, error) {
	dbUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Password,
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

	connect := &MariaDbInstance{db}
	return connect, nil
}

/*
Default Values

Max Life Time: 60
Max Idle Connections: 50
Max Open Connections: 100
*/
func decideDefaultConfigs(cfg MariadbConfig) MariadbConfig {
	if cfg.MaxLifeTime == 0 {
		cfg.MaxLifeTime = 60 * time.Second
	}
	if cfg.MaxIdleConns == 0 {
		cfg.MaxIdleConns = 50
	}
	if cfg.MaxOpenConns == 0 {
		cfg.MaxOpenConns = 100
	}
	return cfg
}

func (connect *MariaDbInstance) CreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.conn.Begin()

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
func (connect *MariaDbInstance) CheckConnection() error {
	// log.Printf("Waiting for Database Connection,,,")
	// time.Sleep(time.Second * 10)

	pingErr := connect.conn.Ping()

	if pingErr != nil {
		log.Printf("[CHECK] Database Ping Error: %v", pingErr)
		return pingErr
	}

	defer connect.conn.Close()

	return nil
}

/*
Query Multiple Rows

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Multiple Row Result
*/
func (connect *MariaDbInstance) SelectMultiple(queryString string, args ...string) (*sql.Rows, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result, err := connect.conn.Query(queryString, arguments...)

	if err != nil {
		log.Printf("[QUERY] Query Error: %v\n", err)

		return nil, err
	}

	defer connect.conn.Close()

	return result, nil
}

/*
Query Single Row

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Single Row Result
*/
func (connect *MariaDbInstance) SelectSingle(queryString string, args ...string) (*sql.Row, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	result := connect.conn.QueryRow(queryString, arguments...)

	if result.Err() != nil {
		log.Printf("[QUERY] Query Error: %v\n", result.Err())

		return nil, result.Err()
	}

	defer connect.conn.Close()

	return result, nil
}

/*
Insert Single Data

@queryString: Query String with prepared statement
@args: Query Parameters
@Return: Insert ID
*/
func (connect *MariaDbInstance) InsertQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	insertResult, insertErr := connect.conn.Exec(queryString, arguments...)

	if insertErr != nil {
		log.Printf("[INSERT] Insert Query Err: %v", insertErr)

		return -99999, insertErr
	}

	defer connect.conn.Close()

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
func (connect *MariaDbInstance) UpdateQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	updateResult, updateErr := connect.conn.Exec(queryString, arguments...)

	if updateErr != nil {
		log.Printf("[UPDATE] Update Query Err: %v", updateErr)

		return -99999, updateErr
	}

	defer connect.conn.Close()

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
func (connect *MariaDbInstance) DeleteQuery(queryString string, args ...string) (int64, error) {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	delResult, delErr := connect.conn.Exec(queryString, arguments...)

	if delErr != nil {
		log.Printf("[DELETE] Delete Query Err: %v", delErr)

		return -99999, delErr
	}

	defer connect.conn.Close()

	// Insert ID
	affectedRow, afftedRowErr := delResult.RowsAffected()

	if afftedRowErr != nil {
		log.Printf("[DELETE] Get Affected Numbers Error: %v", afftedRowErr)

		return -999999, afftedRowErr
	}

	return affectedRow, nil
}
