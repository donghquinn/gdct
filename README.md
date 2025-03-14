# Go-Database-Client

## Introduction
* It's Database Client Package
* It provides creating connection pool, queries, and graceful shutdown

## Dependencies

### Postgres
```zsh
go get -u github.com/lib/pq
```

### Mariadb / Mysql
```zsh
go get -u github.com/go-sql-driver/mysql
```

## Installation

```zsh
go get github.com/donghquinn/gdct
```


## Usage

### Mariadb / mysql

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect("mariadb", &gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

    pingErr := conn.MrCreateTable()

    // ...
}
```

### Postgres
* All the methods are started with 'pg'
    * pgSelectSingle
    * pgSelectMultiple

```go
package main

import "github.com/donghquinn/gdct"

func main() {
    conn, _ := gdct.InitConnect("postgres", &gdct.DBConfig{
        UserName: "test",
        Password: "1234",
        Host: "192.168.0.101",
        Port: 123,
        Database: "test_db",
        MaxLifeTime: 600,
        MaxIdleConns: 50,
        MaxOpenConns: 10
    })

    pingErr := conn.PgCheckConnection()

    // ...
}
```