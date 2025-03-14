package database

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
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

type PostgresInstance struct {
	conn *sql.DB
}

// DB 연결 인스턴스
func InitPostgresConnection(cfg MariadbConfig) (*PostgresInstance, error) {
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.UserName,
		cfg.Password,
		cfg.Host,
		cfg.Password,
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

	connect := &PostgresInstance{db}

	return connect, nil
}

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

// 테이블 생성
func (connect *PostgresInstance) CheckPostgresConnection() error {
	log.Printf("Waiting for Database Connection,,,")
	time.Sleep(time.Second * 10)

	pingErr := connect.conn.Ping()

	if pingErr != nil {
		log.Printf("[DATABASE] Database Ping Error: %v", pingErr)
		return pingErr
	}

	defer connect.conn.Close()

	return nil
}

func (connect *PostgresInstance) CreateTable(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.conn.Begin()

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

// 쿼리
func (connect *PostgresInstance) QueryRows(queryString string, args ...string) (*sql.Rows, error) {
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

// 쿼리
func (connect *PostgresInstance) QueryOne(queryString string, args ...string) (*sql.Row, error) {
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

// 인서트 쿼리
func (connect *PostgresInstance) InsertQuery(queryString string, returns []interface{}, args ...string) error {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	queryResult := connect.conn.QueryRow(queryString, arguments...)

	defer connect.conn.Close()

	if returns != nil {
		// Insert ID
		if scanErr := queryResult.Scan(returns...); scanErr != nil {
			log.Printf("[INSERT] Get Insert Result Scan Error: %v", scanErr)
			return scanErr
		}
	}

	return nil
}

// 인서트 쿼리
func (connect *PostgresInstance) UpdateQuery(queryString string, returns []interface{}, args ...string) error {
	var arguments []interface{}

	for _, arg := range args {
		arguments = append(arguments, arg)
	}

	_, queryErr := connect.conn.Exec(queryString, arguments...)

	defer connect.conn.Close()

	if queryErr != nil {
		// Insert ID

		log.Printf("[INSERT] Get Update Error: %v", queryErr)
		return queryErr
	}

	return nil
}

func (connect *PostgresInstance) InsertMultiple(queryList []string) error {
	ctx := context.Background()

	tx, txErr := connect.conn.Begin()

	if txErr != nil {
		log.Printf("[INSERT_MULTIPLE] Begin Transaction Error: %v", txErr)
		return txErr
	}

	defer tx.Rollback()

	for _, queryString := range queryList {
		_, execErr := tx.ExecContext(ctx, queryString)

		if execErr != nil {
			tx.Rollback()
			log.Printf("[INSERT_MULTIPLE] Insert Querystring Transaction Exec Error: %v", execErr)
			return execErr
		}
	}

	commitErr := tx.Commit()

	if commitErr != nil {
		log.Printf("[INSERT_MULTIPLE] Commit Transaction Error: %v", commitErr)
		return commitErr
	}

	return nil
}
