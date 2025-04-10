package gdct

import "time"

/*
Default Values

Max Life Time: 60
Max Idle Connections: 50
Max Open Connections: 100
*/
func decideDefaultConfigs(cfg DBConfig, dbType DBType) DBConfig {
	if cfg.MaxLifeTime == nil {
		defaultLifetime := 60 * time.Second
		cfg.MaxLifeTime = &defaultLifetime
	}

	if cfg.MaxIdleConns == nil {
		defaultIdleConns := 50
		cfg.MaxIdleConns = &defaultIdleConns
	}

	if cfg.MaxOpenConns == nil {
		defaultOpenConns := 100
		cfg.MaxOpenConns = &defaultOpenConns
	}

	if dbType == PostgreSQL && cfg.SslMode == nil {
		defaultSslMode := "disable"
		cfg.SslMode = &defaultSslMode
	}

	return cfg
}

// Convert String Slice into Interface Slice
func convertArgs(args []string) []interface{} {
	arguments := make([]interface{}, len(args))
	for i, arg := range args {
		arguments[i] = arg
	}
	return arguments
}
