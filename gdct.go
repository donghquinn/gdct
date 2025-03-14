package gdct

import (
	"database/sql"
	"fmt"
	"time"
)

type DBType string

const (
	PostgreSQL DBType = "postgres"
	MariaDB    DBType = "mariadb"
	MysqlDb    DBType = "mysql"
)

type DBConfig struct {
	UserName     string
	Password     string
	Host         string
	Port         int
	Database     string
	SslMode      string
	MaxLifeTime  time.Duration // time.Duration 타입을 권장 (예: 60 * time.Second)
	MaxIdleConns int
	MaxOpenConns int
}

type DataBaseConnector struct {
	*sql.DB
}

// InitConnection: DBConfig에 따라 알맞은 connection pool 생성
func InitConnection(dbType DBType, cfg DBConfig) (*DataBaseConnector, error) {
	switch dbType {
	case MariaDB:
		return InitMariadbConnection(cfg)
	case MysqlDb:
		return InitMariadbConnection(cfg)
	case PostgreSQL:
		return InitPostgresConnection(cfg)
	default:
		return nil, fmt.Errorf("unsupported DB type: %s", dbType)
	}
}
