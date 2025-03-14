package gdct_test

import (
	"testing"

	"github.com/donghquinn/gdct"
)

func TestCheckPostTest(t *testing.T) {
	conn, connErr := gdct.InitConnection("postgres", gdct.DBConfig{
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

func TestCheckMariaTest(t *testing.T) {
	conn, connErr := gdct.InitConnection("mariadb", gdct.DBConfig{
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
	conn, connErr := gdct.InitConnection("mysql", gdct.DBConfig{
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
