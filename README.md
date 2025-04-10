# Go-Database-Client

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

* Check Connection

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
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

    pingErr := conn.MrCheckConnection()

    // ...

    qb := gdct.BuildSelect(gdct.MariaDB, "table_name", "col1").
    Where("col1 = ?", 100).
    OrderBy("col1", "ASC", nil).
    Limit(10).
    Offset(5)

	queryString, args, err := qb.Build()

    queryResult, queryErr := dbCon.QueryRows(queryString, args)
}
```

* select query

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
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

	queryResult, queryErr := conn.MrSelectSingle("SELECT COUNT(example_id) FROM example_table WHERE example_id = ? AND example_status = ?", "1234", "1")

    // ...
}

```


### Postgres
* All the methods are started with 'pg'
    * pgSelectSingle
    * pgSelectMultiple

* Check Connection 

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.Postgres, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        SslMode: "disable",
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

    pingErr := conn.PgCheckConnection()

    // ...
}
```

* Select


```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect(gdct.Postgres, gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

	queryResult, queryErr := conn.PgSelectSingle("SELECT COUNT(example_id) FROM example_table WHERE example_id = $1 AND example_status = $2", "1234", "1")

    // ...
}
```