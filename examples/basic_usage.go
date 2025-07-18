package examples

import (
	"fmt"
	"log"
	"time"

	"github.com/donghquinn/gdct"
	_ "github.com/lib/pq"
)

func BasicUsageExample() {
	// Example 1: Basic connection and query
	basicExample()

	// Example 2: Dynamic query building
	dynamicQueryExample()

	// Example 3: Insert, Update, Delete operations
	crudOperationsExample()

	// Example 4: Transactions
	transactionExample()
}

func basicExample() {
	fmt.Println("=== Basic Usage Example ===")

	// Connect to PostgreSQL
	db, err := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
		Host:     "localhost",
		Port:     5432,
		Database: "testdb",
		UserName: "user",
		Password: "password",
	})
	if err != nil {
		log.Printf("Connection failed: %v", err)
		return
	}
	defer db.Close()

	// Simple SELECT query
	query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users").
		Where("age > ?", 18).
		Where("status = ?", "active").
		OrderBy("created_at", "DESC", nil).
		Limit(10).
		Build()

	if err != nil {
		log.Printf("Query build failed: %v", err)
		return
	}

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Args: %v\n", args)

	// Execute query (uncomment when you have a real database)
	// rows, err := db.QueryBuilderRows(query, args)
	// if err != nil {
	//     log.Printf("Query execution failed: %v", err)
	//     return
	// }
	// defer rows.Close()
	//
	// // Process results...
}

func dynamicQueryExample() {
	fmt.Println("\n=== Dynamic Query Building ===")

	// Simulate user search parameters
	searchName := "John"
	minAge := 25
	sortBy := "name"
	allowedSortColumns := map[string]bool{
		"name":       true,
		"age":        true,
		"created_at": true,
	}

	// Build query dynamically
	qb := gdct.BuildSelect(gdct.PostgreSQL, "users", "id", "name", "email", "age")

	// Add conditions only if parameters are provided
	if searchName != "" {
		qb = qb.Where("name ILIKE ?", "%"+searchName+"%")
	}
	if minAge > 0 {
		qb = qb.Where("age >= ?", minAge)
	}
	if sortBy != "" {
		qb = qb.OrderBy(sortBy, "ASC", allowedSortColumns)
	}

	query, args, err := qb.Build()
	if err != nil {
		log.Printf("Dynamic query build failed: %v", err)
		return
	}

	fmt.Printf("Dynamic Query: %s\n", query)
	fmt.Printf("Args: %v\n", args)
}

func crudOperationsExample() {
	fmt.Println("\n=== CRUD Operations ===")

	// INSERT example
	userData := map[string]interface{}{
		"name":       "John Doe",
		"email":      "john@example.com",
		"age":        30,
		"created_at": time.Now(),
	}

	insertQuery, insertArgs, err := gdct.BuildInsert(gdct.PostgreSQL, "users").
		Values(userData).
		Returning("id").
		Build()

	if err != nil {
		log.Printf("Insert query build failed: %v", err)
		return
	}

	fmt.Printf("Insert Query: %s\n", insertQuery)
	fmt.Printf("Insert Args: %v\n", insertArgs)

	// UPDATE example
	updateData := map[string]interface{}{
		"last_login":  time.Now(),
		"login_count": "login_count + 1", // Raw SQL expression
		"updated_at":  time.Now(),
	}

	updateQuery, updateArgs, err := gdct.BuildUpdate(gdct.PostgreSQL, "users").
		Set(updateData).
		Where("id = ?", 123).
		Build()

	if err != nil {
		log.Printf("Update query build failed: %v", err)
		return
	}

	fmt.Printf("Update Query: %s\n", updateQuery)
	fmt.Printf("Update Args: %v\n", updateArgs)

	// DELETE example
	deleteQuery, deleteArgs, err := gdct.BuildDelete(gdct.PostgreSQL, "users").
		Where("last_login < ?", time.Now().AddDate(0, -6, 0)).
		Where("status = ?", "inactive").
		Build()

	if err != nil {
		log.Printf("Delete query build failed: %v", err)
		return
	}

	fmt.Printf("Delete Query: %s\n", deleteQuery)
	fmt.Printf("Delete Args: %v\n", deleteArgs)
}

func transactionExample() {
	fmt.Println("\n=== Transaction Example ===")

	// Example: Create an order and update user's order count
	userID := 123
	orderTotal := 99.99

	// Prepare multiple queries for transaction
	queries := []gdct.PreparedQuery{
		{
			Query:  "INSERT INTO orders (user_id, total, created_at) VALUES ($1, $2, $3)",
			Params: []interface{}{userID, orderTotal, time.Now()},
		},
		{
			Query:  "UPDATE users SET orders_count = orders_count + 1, updated_at = $2 WHERE id = $1",
			Params: []interface{}{userID, time.Now()},
		},
	}

	fmt.Printf("Transaction queries prepared: %d\n", len(queries))
	for i, q := range queries {
		fmt.Printf("  Query %d: %s\n", i+1, q.Query)
		fmt.Printf("  Params %d: %v\n", i+1, q.Params)
	}

	// Execute transaction (uncomment when you have a real database)
	// db, err := gdct.InitConnection(gdct.PostgreSQL, config)
	// if err != nil {
	//     log.Printf("Connection failed: %v", err)
	//     return
	// }
	// defer db.Close()
	//
	// results, err := db.PgInsertMultiple(queries)
	// if err != nil {
	//     log.Printf("Transaction failed: %v", err)
	//     return
	// }
	//
	// fmt.Printf("Transaction completed successfully. Results: %d\n", len(results))
}
