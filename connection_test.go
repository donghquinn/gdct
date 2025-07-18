package gdct

import (
	"testing"
	"time"
)

func TestCheckPostTest(t *testing.T) {
	sslMode := "disable" // Only Postgres

	conn, connErr := InitConnection(PostgreSQL, DBConfig{
		Host:     "192.168.0.241",
		Port:     5432,
		UserName: "its",
		Password: "1234",
		Database: "its",
		SslMode:  &sslMode,
	})

	if connErr != nil {
		t.Fatalf("[POST_CHECK] Create Connection Test Error: %v", connErr)
	}

	pingErr := conn.PgCheckConnection()
	if pingErr != nil {
		t.Fatalf("[POST_CHECK] Connection Test Error: %v", pingErr)
	}
}

func TestCheckSqlTest(t *testing.T) {
	conn, connErr := InitConnection(Sqlite, DBConfig{
		Database: "./db.sqlite",
	})

	if connErr != nil {
		t.Fatalf("[SQLITE_CHECK] Create Connection Test Error: %v", connErr)
	}

	pingErr := conn.MrCheckConnection()
	if pingErr != nil {
		t.Fatalf("[SQLITE_CHECK] Connection Test Error: %v", pingErr)
	}
}

func TestCheckMariaTest(t *testing.T) {
	lifetime := time.Duration(600) * time.Second
	idleConns := 10
	openConns := 50

	conn, connErr := InitConnection(MariaDB, DBConfig{
		Host:         "192.168.0.241",
		Port:         3306,
		UserName:     "its",
		Password:     "1234",
		Database:     "its",
		MaxLifeTime:  &lifetime,
		MaxIdleConns: &idleConns,
		MaxOpenConns: &openConns,
	})

	if connErr != nil {
		t.Fatalf("[MARIA_CHECK] Create Connection Test Error: %v", connErr)
	}

	pingErr := conn.MrCheckConnection()
	if pingErr != nil {
		t.Fatalf("[MARIA_CHECK] Connection Test Error: %v", pingErr)
	}
}

func TestCheckMysqlTest(t *testing.T) {
	conn, connErr := InitConnection(Mysql, DBConfig{
		Host:     "192.168.0.241",
		Port:     3306,
		UserName: "its",
		Password: "1234",
		Database: "its",
	})

	if connErr != nil {
		t.Fatalf("[MYSQL_CHECK] Create Connection Test Error: %v", connErr)
	}

	pingErr := conn.MrCheckConnection()
	if pingErr != nil {
		t.Fatalf("[MYSQL_CHECK] Connection Test Error: %v", pingErr)
	}
}
