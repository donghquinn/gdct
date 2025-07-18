package gdct_test

import (
	"fmt"
	"log"

	"github.com/donghquinn/gdct"
	_ "github.com/go-sql-driver/mysql" // MySQL/MariaDB
	_ "github.com/lib/pq"              // PostgreSQL
	_ "github.com/mattn/go-sqlite3"    // SQLite
)

func MultiDatabaseExample() {
	// Example showing the same query across different databases
	demonstrateMultiDatabaseSupport()
}

func demonstrateMultiDatabaseSupport() {
	fmt.Println("=== Multi-Database Support Example ===")

	// Same query logic for different databases
	databases := []struct {
		name   string
		dbType gdct.DBType
		config gdct.DBConfig
	}{
		{
			name:   "PostgreSQL",
			dbType: gdct.PostgreSQL,
			config: gdct.DBConfig{
				Host:     "localhost",
				Port:     5432,
				Database: "testdb",
				UserName: "user",
				Password: "password",
			},
		},
		{
			name:   "MySQL",
			dbType: gdct.MariaDB,
			config: gdct.DBConfig{
				Host:     "localhost",
				Port:     3306,
				Database: "testdb",
				UserName: "user",
				Password: "password",
			},
		},
		{
			name:   "SQLite",
			dbType: gdct.Sqlite,
			config: gdct.DBConfig{
				Database: "./test.db",
			},
		},
	}

	for _, db := range databases {
		fmt.Printf("\n--- %s Example ---\n", db.name)
		demonstrateDatabase(db.dbType, db.config)
	}
}

func demonstrateDatabase(dbType gdct.DBType, config gdct.DBConfig) {
	// Build identical query logic for different databases
	// GDCT handles the dialect differences automatically

	// SELECT query
	selectQuery, selectArgs, err := gdct.BuildSelect(dbType, "products p").
		Select("p.id", "p.name", "p.price", "c.name as category_name").
		LeftJoin("categories c", "c.id = p.category_id").
		Where("p.price > ?", 10.00).
		Where("p.in_stock = ?", true).
		OrderBy("p.created_at", "DESC", nil).
		Limit(20).
		Build()

	if err != nil {
		log.Printf("SELECT query build failed: %v", err)
		return
	}

	fmt.Printf("SELECT Query: %s\n", selectQuery)
	fmt.Printf("SELECT Args: %v\n", selectArgs)

	// INSERT query
	productData := map[string]interface{}{
		"name":        "New Product",
		"price":       29.99,
		"category_id": 1,
		"in_stock":    true,
	}

	var insertQuery string
	var insertArgs []interface{}

	if dbType == gdct.PostgreSQL {
		// PostgreSQL supports RETURNING
		insertQuery, insertArgs, err = gdct.BuildInsert(dbType, "products").
			Values(productData).
			Returning("id").
			Build()
	} else {
		// MySQL/SQLite don't use RETURNING
		insertQuery, insertArgs, err = gdct.BuildInsert(dbType, "products").
			Values(productData).
			Build()
	}

	if err != nil {
		log.Printf("INSERT query build failed: %v", err)
		return
	}

	fmt.Printf("INSERT Query: %s\n", insertQuery)
	fmt.Printf("INSERT Args: %v\n", insertArgs)

	// UPDATE query
	updateData := map[string]interface{}{
		"price":    39.99,
		"in_stock": false,
	}

	updateQuery, updateArgs, err := gdct.BuildUpdate(dbType, "products").
		Set(updateData).
		Where("id = ?", 123).
		Build()

	if err != nil {
		log.Printf("UPDATE query build failed: %v", err)
		return
	}

	fmt.Printf("UPDATE Query: %s\n", updateQuery)
	fmt.Printf("UPDATE Args: %v\n", updateArgs)

	// Connection example (commented out as it requires actual databases)
	/*
		db, err := gdct.InitConnection(dbType, config)
		if err != nil {
			log.Printf("Connection to %s failed: %v", dbType, err)
			return
		}
		defer db.Close()

		// Test connection
		switch dbType {
		case gdct.PostgreSQL:
			err = db.PgCheckConnection()
		case gdct.MariaDB:
			err = db.MrCheckConnection()
		case gdct.Sqlite:
			err = db.SqCheckConnection()
			// SQLite specific optimizations
			db.SqEnableWAL()
			db.SqEnableForeignKeys()
		}

		if err != nil {
			log.Printf("Connection check failed: %v", err)
			return
		}

		fmt.Printf("✅ Successfully connected to %s\n", dbType)
	*/
}

// Example of database-specific features
func sqliteSpecificFeatures() {
	fmt.Println("\n=== SQLite Specific Features ===")

	// SQLite connection with optimizations
	db, err := gdct.InitConnection(gdct.Sqlite, gdct.DBConfig{
		Database: "./example.db",
	})
	if err != nil {
		log.Printf("SQLite connection failed: %v", err)
		return
	}
	defer db.Close()

	// Enable SQLite optimizations
	if err := db.SqEnableWAL(); err != nil {
		log.Printf("Failed to enable WAL: %v", err)
	} else {
		fmt.Println("✅ WAL mode enabled for better concurrency")
	}

	if err := db.SqEnableForeignKeys(); err != nil {
		log.Printf("Failed to enable foreign keys: %v", err)
	} else {
		fmt.Println("✅ Foreign key constraints enabled")
	}

	// Get SQLite version
	version, err := db.SqGetVersion()
	if err != nil {
		log.Printf("Failed to get SQLite version: %v", err)
	} else {
		fmt.Printf("SQLite version: %s\n", version)
	}

	// Maintenance operations
	fmt.Println("Performing maintenance operations...")

	if err := db.SqVacuum(); err != nil {
		log.Printf("VACUUM failed: %v", err)
	} else {
		fmt.Println("✅ VACUUM completed")
	}

	if err := db.SqAnalyze(); err != nil {
		log.Printf("ANALYZE failed: %v", err)
	} else {
		fmt.Println("✅ ANALYZE completed")
	}
}
