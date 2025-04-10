package gdct_test

import (
	"testing"

	"github.com/donghquinn/gdct"
)

func TestCheckPostTest(t *testing.T) {
	conn, connErr := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
		Host:     "192.168.0.241",
		Port:     5432,
		UserName: "its",
		Password: "1234",
		Database: "its",
		SslMode:  "disable",
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
	conn, connErr := gdct.InitConnection(gdct.MariaDB, gdct.DBConfig{
		Host:         "192.168.0.241",
		Port:         3306,
		UserName:     "its",
		Password:     "1234",
		Database:     "its",
		MaxLifeTime:  600,
		MaxIdleConns: 10,
		MaxOpenConns: 50,
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
