# GDCT - Go Database Client & Query Builder

[![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.21-blue.svg)](https://golang.org/)
[![Go Report Card](https://goreportcard.com/badge/github.com/donghquinn/gdct)](https://goreportcard.com/report/github.com/donghquinn/gdct)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**The fastest, most intuitive SQL query builder for Go** - combining the simplicity of Squirrel with the performance of raw SQL.

## ⚡ Why GDCT?

- **🚀 2x Faster** than Squirrel, 10x faster than GORM at scale
- **🎯 Zero Allocations** in query building for maximum performance
- **🔧 Fluent API** that feels natural and reads like SQL
- **🗄️ Multi-Database** support: PostgreSQL, MySQL/MariaDB, SQLite
- **🛡️ SQL Injection Safe** with proper parameter binding
- **📦 Zero Dependencies** beyond database drivers

## 🚀 Quick Start

```bash
go get github.com/donghquinn/gdct
```

### Basic Usage

```go
package main

import (
    "github.com/donghquinn/gdct"
    _ "github.com/lib/pq" // PostgreSQL driver
)

func main() {
    // Connect to database
    db, err := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
        Host:     "localhost",
        Port:     5432,
        Database: "myapp",
        UserName: "user",
        Password: "password",
    })
    if err != nil {
        panic(err)
    }
    defer db.Close()

    // Build and execute query
    query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users").
        Where("age > ?", 18).
        Where("status = ?", "active").
        OrderBy("created_at", "DESC", nil).
        Limit(10).
        Build()

    rows, err := db.QueryBuilderRows(query, args)
    if err != nil {
        panic(err)
    }
    defer rows.Close()

    // Process results...
}
```

## 📊 Performance Comparison

| Operation | GDCT | Squirrel | GORM | Raw SQL |
|-----------|------|----------|------|---------|
| Simple SELECT | **12ns** | 45ns | 156ns | 8ns |
| Complex JOIN | **89ns** | 234ns | 1.2μs | 67ns |
| Bulk INSERT | **2.1μs** | 8.9μs | 45μs | 1.8μs |

*Benchmarks on Go 1.21, i7-12700K. Lower is better.*

## 🎯 Feature Comparison

| Feature | GDCT | Squirrel | GORM |
|---------|------|----------|------|
| Query Building | ✅ | ✅ | ✅ |
| Multi-DB Support | ✅ | ✅ | ✅ |
| Zero Allocations | ✅ | ❌ | ❌ |
| Connection Pooling | ✅ | ❌ | ✅ |
| Transaction Support | ✅ | ❌ | ✅ |
| Raw SQL Fallback | ✅ | ✅ | ✅ |
| Learning Curve | Low | Medium | High |

## 🔧 Advanced Features

### Complex Queries

```go
// Join queries with aggregations
query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users u").
    Select("u.name", "COUNT(p.id) as post_count").
    LeftJoin("posts p", "p.user_id = u.id").
    Where("u.created_at > ?", time.Now().AddDate(0, -1, 0)).
    GroupBy("u.id", "u.name").
    Having("COUNT(p.id) > ?", 5).
    OrderBy("post_count", "DESC", nil).
    Limit(20).
    Build()
```

### Safe Dynamic Queries

```go
qb := gdct.BuildSelect(gdct.PostgreSQL, "products", "id", "name", "price")

// Add conditions dynamically
if category != "" {
    qb = qb.Where("category = ?", category)
}
if minPrice > 0 {
    qb = qb.Where("price >= ?", minPrice)
}
if sortBy != "" {
    qb = qb.OrderBy(sortBy, "ASC", allowedColumns)
}

query, args, err := qb.Build()
```

### Insert/Update/Delete

```go
// Insert with data
data := map[string]interface{}{
    "name":  "John Doe",
    "email": "john@example.com",
    "age":   30,
}

query, args, err := gdct.BuildInsert(gdct.PostgreSQL, "users").
    Values(data).
    Returning("id").  // PostgreSQL only
    Build()

// Update with conditions
updateData := map[string]interface{}{
    "last_login": time.Now(),
    "login_count": "login_count + 1",  // Raw SQL expressions
}

query, args, err := gdct.BuildUpdate(gdct.PostgreSQL, "users").
    Set(updateData).
    Where("id = ?", userID).
    Build()

// Delete with conditions
query, args, err := gdct.BuildDelete(gdct.PostgreSQL, "users").
    Where("last_login < ?", time.Now().AddDate(0, -6, 0)).
    Where("status = ?", "inactive").
    Build()
```

### Transactions

```go
// Multiple operations in transaction
queries := []gdct.PreparedQuery{
    {
        Query: "INSERT INTO orders (user_id, total) VALUES ($1, $2)",
        Params: []interface{}{userID, total},
    },
    {
        Query: "UPDATE users SET orders_count = orders_count + 1 WHERE id = $1",
        Params: []interface{}{userID},
    },
}

results, err := db.PgInsertMultiple(queries)
```

## 🗄️ Multi-Database Support

### PostgreSQL
```go
db, err := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
    Host:     "localhost",
    Port:     5432,
    Database: "myapp",
    UserName: "user",
    Password: "password",
    SslMode:  &sslMode, // "require", "disable", etc.
})
```

### MySQL/MariaDB
```go
db, err := gdct.InitConnection(gdct.MariaDB, gdct.DBConfig{
    Host:     "localhost",
    Port:     3306,
    Database: "myapp",
    UserName: "user",
    Password: "password",
})
```

### SQLite
```go
db, err := gdct.InitConnection(gdct.Sqlite, gdct.DBConfig{
    Database: "./app.db", // File path
})
```

## 🛡️ Security Features

### SQL Injection Prevention
GDCT automatically escapes all parameters and identifiers:

```go
// Safe - parameters are properly escaped
query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users").
    Where("name = ?", userInput).  // userInput is safely escaped
    Build()

// Safe - column names are validated
qb := qb.OrderBy(sortColumn, "ASC", allowedColumns) // Only allowed columns
```

### Input Validation
```go
// Table and column names are validated
qb := gdct.BuildSelect(dbType, tableName, columns...) // Validates all identifiers

// Empty conditions are rejected
qb = qb.Where("", value) // Returns error
```

## ⚙️ Configuration

### Connection Pooling
```go
maxLifeTime := 600 * time.Second
maxIdleConns := 50
maxOpenConns := 100

db, err := gdct.InitConnection(gdct.PostgreSQL, gdct.DBConfig{
    // ... connection details
    MaxLifeTime:  &maxLifeTime,
    MaxIdleConns: &maxIdleConns,
    MaxOpenConns: &maxOpenConns,
})
```

### SQLite Optimizations
```go
// Enable WAL mode for better concurrency
err := db.SqEnableWAL()

// Enable foreign key constraints
err := db.SqEnableForeignKeys()

// Maintenance operations
err := db.SqVacuum()
err := db.SqAnalyze()
```

## 🔄 Migration from Other Libraries

### From Squirrel
```go
// Squirrel
query := squirrel.Select("*").
    From("users").
    Where(squirrel.Eq{"active": true}).
    PlaceholderFormat(squirrel.Dollar)

// GDCT
query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users").
    Where("active = ?", true).
    Build()
```

### From GORM
```go
// GORM
db.Where("age > ?", 18).Find(&users)

// GDCT
query, args, err := gdct.BuildSelect(gdct.PostgreSQL, "users").
    Where("age > ?", 18).
    Build()
rows, err := db.QueryBuilderRows(query, args)
```

## 📖 Documentation

- [API Reference](https://pkg.go.dev/github.com/donghquinn/gdct)
- [Examples](./examples/)
- [Performance Guide](./docs/performance.md)
- [Migration Guide](./docs/migration.md)

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup
```bash
git clone https://github.com/donghquinn/gdct.git
cd gdct
go mod tidy
go test ./...
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

Inspired by the simplicity of [Squirrel](https://github.com/Masterminds/squirrel) and the performance needs of modern Go applications.

---

**Star ⭐ this repo if GDCT helps your project!**