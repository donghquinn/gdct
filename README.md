# Go-Database-Client

## Notice Since v1.0.0
* Added Query Builder
    * Query Build for dynamic queries
* If you only need query builder, try [go-query-builder](https://github.com/donghquinn/gqbd)!

## Introduction
* It's Database Client Package
* I've combined [go-query-builder](https://github.com/donghquinn/gqbd) with database client package
    * This package will allow to building queries and connect database pool
* In order word, it's combined tool for querying database includeing query strings and connections.

## Dependencies
* It depends on postgres and mysql driver

### Postgres
```zsh
go get -u github.com/lib/pq
```

### Mariadb / Mysql
```zsh
go get -u github.com/go-sql-driver/mysql
```

---

## Installation

```zsh
go get github.com/donghquinn/gdct
```

---

## Usage

* Every  Single Method will close connection after transaction commited.
* So you have to open connection again for every time.
* Postgres start with Pg and Mariadb start with Mr
* (2025-04-10 Added) QueryBuilderOneRow() and QueryBuilderRows() is the mothods for builded query strings
    * QueryBuilderOneRow will query single row
    * QueryBuilderRows will query mutliple rows

### Mariadb / mysql

* Select Rows with query builder
    * Use QueryBuilderRows for multiple rows and QueryBuilderOneRow for single rows

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    var (
        maxLifeTime = 600
        maxIdelConns = 50
        maxOpenConns = 10
    )

    conn, _ := gdct.InitConnect(gdct.MariaDB, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        MaxLifeTime: &maxLifeTime,
        MaxIdleConns: &maxIdelConns,
        MaxOpenConns: &maxOpenConns
    })

    // ...

    qb := gdct.BuildSelect(gdct.MariaDB, "table_name", "col1").
        Where("col1 = ?", 100).
        OrderBy("col1", "ASC", nil).
        Limit(10).
        Offset(5)

	queryString, args, err := qb.Build()

    queryResult, queryErr := conn.QueryBuilderRows(queryString, args)
}
```


* select one row query
    * You can use MrSelect... method for query string, not using query builder

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.MariaDB, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
    })

	queryResult, queryErr := conn.MrSelectSingle("SELECT COUNT(example_id) FROM example_table WHERE example_id = ? AND example_status = ?", "1234", "1")

    // ...
}
```

* Insert Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.MariaDB, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "INSERT INTO example_table (column1, column2, column3) VALUES ($1, $2, $3)",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.MrInsertMultiple(queryList)

    // ...
}

```

* Update Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.MariaDB, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "UPDATE example_table SET column1 = $1, column2 = $2, column = $3",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.MrUpdateMultiple(queryList)

    // ...
}

```

* DELETE Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.MariaDB, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "DELETE example_table WHERE column1 = $1, column2 = $2, column = $3",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.MrDeleteMultiple(queryList)

    // ...
}

```

### Postgres
* All the methods are started with 'pg'
    * pgSelectSingle
    * pgSelectMultiple

* Select Rows with query builder
    * Use QueryBuilderRows for multiple rows and QueryBuilderOneRow for single rows


```go
package main

import "github.com/donghquinn/gdct"

func main() {
    var (
        sslMode = "disable"
        maxLifeTime = 600
        maxIdelConns = 50
        maxOpenConns = 10
    )

    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode,
        MaxLifeTime: &maxLifeTime,
        MaxIdleConns: &maxIdelConns,
        MaxOpenConns: &maxOpenConns
    })

    // ...

    // Query Building
    qb := gqbd.BuildSelect(gqbd.PostgreSQL, "example_table e", "e.id", "e.name", "u.user").
    LeftJoin("user_table u", "u.user_id = e.id")

	if userName != "" {
		qb = qb.Where("u.user_name LIKE ?", "%"+userName+"%")
	}

	// title이 비어있지 않은 경우에만 조건 추가
	if title != "" {
		qb = qb.Where("e.name LIKE ?", "%"+title+"%")
	}
	// 상태 조건은 항상 추가
	qb = qb.Where("e.example_status = ?", "1")

	// 정렬, 오프셋, 제한 설정
	qb = qb.OrderBy(orderByColumn, "DESC", nil).
		Offset(offset).
		Limit(limit)

	queryString, args, err := qb.Build()


    // Send Query and get returns
    queryResult, queryErr := conn.QueryBuilderRows(queryString, args)
}
```

* Select without query builder


```go
package main

import "github.com/donghquinn/gdct"

func main() {
    sslMode := "disable"

    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode
    })

	queryResult, queryErr := conn.PgSelectSingle("SELECT COUNT(example_id) FROM example_table WHERE example_id = $1 AND example_status = $2", "1234", "1")

    // ...
}
```

* Insert
    * Can get returning values from INSERT queries

```go
    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode
    })

    var exampleId int64

    insertErr := conn.PgInsertQuery(
        "INSERT example (exam_name, exam_status) VALUES ($1, $2) RETURNING exam_id", 
        []interface{}{&exampleSeqe}, // Returning Value
        "Example Name", "10" // Arguments
    )

```

* With Query Building

```go
	data := map[string]interface{}{
		"col1": 200,
		"col2": "test",
	}

	qb := gqbd.BuildInsert(gqbd.PostgreSQL, "table_name").
		Values(data)

	query, args, err := qb.Build()

    // .... 

    insertResult, insertErr := conn.QueryBuilderInsert(query, args)
    // ...
```


* Update
    * With Query Building

```go
	data := map[string]interface{}{
		"col1": 200,
		"col2": "test",
	}

	qb := gqbd.BuildUpdate(gqbd.PostgreSQL, "table_name").
		Set(data).
        Where("col1 = ?", 100)

	query, args, err := qb.Build()

    // .... 

    insertResult, insertErr := conn.QueryBuilderUpdate(query, args)
    // ...
```

* Insert Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "INSERT INTO example_table (column1, column2, column3) VALUES ($1, $2, $3)",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.PgInsertMultiple(queryList)

    // ...
}

```

* Update Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode,
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "UPDATE example_table SET column1 = $1, column2 = $2, column = $3",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.PgUpdateMultiple(queryList)

    // ...
}

```


* DELETE Multiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.PostgreSQL, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: &sslMode,
    })

    queryList := make([]gdct.PreparedQuery, len(dataList))

    for _ data := range dataList {
        queryData := gdct.PreparedQuery{
            Query: "DELETE example_table WHERE column1 = $1, column2 = $2, column = $3",
            Params: []interface{}{
                data.exampleItem,
                data.exampleItem2,
                data.exampleItem3,
            }
        }

        queryList = queryList.append(queryList, queryData)
    }

	insertResultList, queryErr := conn.PgUpdateMultiple(queryList)

    // ...
}

```