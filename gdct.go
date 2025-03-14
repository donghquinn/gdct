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
	MaxLifeTime  time.Duration
	MaxIdleConns int
	MaxOpenConns int
}

type DataBaseConnector struct {
	*sql.DB
}

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
