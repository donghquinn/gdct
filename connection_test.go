package gdct_test

import (
	"testing"

	"github.com/donghquinn/gdct"
)

func CheckPostTest(t *testing.T) {
	conn, _ := gdct.InitConnection("postgres", gdct.DBConfig{
		Host:     "192.168.0.241",
		Port:     5432,
		UserName: "its",
		Password: "1234",
		Database: "its",
	})

	pingErr := conn.PgCheckConnection()
	if pingErr != nil {
		t.Fatalf("[POST_CHECK] Connection Test Error: %v", pingErr)
	}
}

func CheckMariaTest(t *testing.T) {
	conn, _ := gdct.InitConnection("mariadb", gdct.DBConfig{
		Host:     "192.168.0.241",
		Port:     3306,
		UserName: "its",
		Password: "1234",
		Database: "its",
	})

	pingErr := conn.MrCheckConnection()
	if pingErr != nil {
		t.Fatalf("[MARIA_CHECK] Connection Test Error: %v", pingErr)
	}
}

func CheckMysqlTest(t *testing.T) {
	conn, _ := gdct.InitConnection("mysql", gdct.DBConfig{
		Host:     "192.168.0.241",
		Port:     3306,
		UserName: "its",
		Password: "1234",
		Database: "its",
	})

	pingErr := conn.MrCheckConnection()
	if pingErr != nil {
		t.Fatalf("[MYSQL_CHECK] Connection Test Error: %v", pingErr)
	}
}
