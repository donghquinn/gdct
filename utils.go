package gdct

import "time"

/*
Default Values

Max Life Time: 60
Max Idle Connections: 50
Max Open Connections: 100
*/
func decideDefaultConfigs(cfg DBConfig) DBConfig {
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

// Convert String Slice into Interface Slice
func convertArgs(args []string) []interface{} {
	arguments := make([]interface{}, len(args))
	for i, arg := range args {
		arguments[i] = arg
	}
	return arguments
}
