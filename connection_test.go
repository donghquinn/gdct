package gdct_test

import (
	"testing"
	"time"

	"github.com/donghquinn/gdct"
)

func TestCheckPostTest(t *testing.T) {
	sslMode := "disable" // Only Postgres

	conn, connErr := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
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

// func TestQueryTest(t *testing.T) {
// 	conn, connErr := gdct.InitConnection("postgres", gdct.DBConfig{
// 		Host:     "192.168.0.241",
// 		Port:     5432,
// 		UserName: "its",
// 		Password: "1234",
// 		Database: "its",
// 		SslMode:  "disable",
// 	})

// 	if connErr != nil {
// 		t.Fatalf("[POST_SELECT_SINGLE] Create Connection Test Error: %v", connErr)
// 	}

// 	queryResult, queryErr := conn.PgSelectSingle("SELECT COUNT(example_id) FROM example_table WHERE example_id = $1", "1234")
// 	if queryErr != nil {
// 		t.Fatalf("[POST_SELECT_SINGLE] Connection Test Error: %v", queryErr)
// 	}

// 	var totalCount int64

// 	if scanErr := queryResult.Scan(&totalCount); scanErr != nil {
// 		t.Fatalf("[POST_SELECT_SINGLE] Scan Error: %v", scanErr)
// 	}

// 	t.Logf("[POST_SELECT_SINGLE] Total Count :%d", totalCount)
// }

func TestCheckMariaTest(t *testing.T) {
	lifetime := time.Duration(600) * time.Second
	idleConns := 10
	openConns := 50

	conn, connErr := gdct.InitConnection(gdct.MariaDB, gdct.DBConfig{
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
	conn, connErr := gdct.InitConnection(gdct.Mysql, gdct.DBConfig{
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
