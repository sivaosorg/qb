# QueryBuilder - qb

---

QueryBuilder is a powerful and flexible Go library for building SQL queries with ease. It simplifies the process of constructing SQL statements, making it effortless to interact with your database.

## Features

- SELECT Queries: Construct complex SELECT queries with dynamic columns, conditions, and JOIN operations.
- UPDATE Queries: Build UPDATE statements to modify database records with ease.
- DELETE Queries: Easily create DELETE statements to remove data from your tables.
- JSON/JSONB Support: Seamlessly handle JSON and JSONB data types in your queries.
- JOIN Operations: Perform INNER JOIN, LEFT JOIN, RIGHT JOIN, and FULL JOIN operations with JSON/JSONB columns.
- Parameterized Queries: Safely construct parameterized queries to prevent SQL injection.
<!-- - Alias Support: Assign aliases to tables for more readable and maintainable queries. -->

## Installation

```bash
go get -u github.com/sivaosorg/qb
```

## Table of Contents

- [QueryBuilder - qb](#querybuilder---qb)
  - [Features](#features)
  - [Installation](#installation)
  - [Table of Contents](#table-of-contents)
    - [Selects, Ordering, Limit and Offset](#selects-ordering-limit-and-offset)
    - [GroupBy / Having](#groupby--having)
    - [Where, AndWhere and OrWhere clauses](#where-andwhere-and-orwhere-clauses)
    - [WhereIn and WhereNotIn clauses](#wherein-and-wherenotin-clauses)
    - [WhereNull and WhereNotNull clauses](#wherenull-and-wherenotnull-clauses)
    - [WhereExists and WhereNotExists clauses](#whereexists-and-wherenotexists-clauses)
    - [WhereBetween and WhereNotBetween clauses](#wherebetween-and-wherenotbetween-clauses)
    - [Left / Right / Inner / Left Outer Joins](#left--right--inner--left-outer-joins)
    - [Insert](#insert)
    - [Update](#update)
    - [Drop, Truncate and Rename](#drop-truncate-and-rename)
    - [Increment and Decrement](#increment-and-decrement)
    - [Union and Union All](#union-and-union-all)
    - [Transaction Mode](#transaction-mode)
    - [Aggregates](#aggregates)
    - [Create Table](#create-table)
    - [Add / Modify / Drop columns](#add--modify--drop-columns)
    - [Chunk](#chunk)
  - [Ref](#ref)
  - [Contribution](#contribution)

### Selects, Ordering, Limit and Offset

You might not always need to retrieve all columns from a database table. With the select method, you have the flexibility to define a custom select clause for your query:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=database_username dbname=database_name password='password' sslmode=disable"))

func main() {
	test01()
}

func test01() {
	query := db.
		Table("or_user").
		Select("or_user.user_id, or_user.user_name, or_user.phone , or_user.modifydate")

	result, err := query.OrderBy("or_user.user_id", "DESC").Limit(5).Offset(0).Get()
	if err != nil {
		panic(err)
	}
	query.PrintQuery()
	for k, v := range result {
		for field, value := range v {
			log.Println(fmt.Sprintf("key = %v => field = %v, value = %v", k, field, value))
		}
	}
}
```

### GroupBy / Having

The GroupBy and Having methods may be used to group the query results.

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test02()
}

func test02() {
	query := db.
		Table("or_user").
		Select("or_user.user_id, or_user.user_name, or_user.phone , or_user.modifydate")

	result, err := query.GroupBy("or_user.user_id").Having("or_user.user_id", ">", 1).Get()
	if err != nil {
		panic(err)
	}
	query.PrintQuery()
	for k, v := range result {
		for field, value := range v {
			log.Println(fmt.Sprintf("key = %v => field = %v, value = %v", k, field, value))
		}
	}
}
```

### Where, AndWhere and OrWhere clauses

The simplest form of the "where" function necessitates three arguments. The initial argument designates the column name, followed by the second argument, which specifies an operator drawn from the set of supported operators within the database. Last but not least, the third argument entails the value against which the column is to be evaluated.

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test03()
}

func test03() {
	query := db.
		Table("or_user").
		Select("or_user.user_id, or_user.user_name, or_user.phone , or_user.modifydate")

	query = query.Where("or_user.user_id", "=", 611235).
		AndWhere("or_user_role.user_id", "=", 611235)

	result, err := query.OrderBy("or_user.user_id", "DESC").Limit(5).Offset(0).Get()
	if err != nil {
		panic(err)
	}
	query.PrintQuery()
	for k, v := range result {
		for field, value := range v {
			log.Println(fmt.Sprintf("key = %v => field = %v, value = %v", k, field, value))
		}
	}
}
```

### WhereIn and WhereNotIn clauses

```go
result, err := db.Table("or_user").WhereIn("id", []int64{1, 2, 3}).OrWhereIn("name", []string{"Aris", "Jake"}).Get()
```

### WhereNull and WhereNotNull clauses

```go
result, err := db.Table("or_user").WhereNull("user_id").OrWhereNotNull("username").Get()
```

### WhereExists and WhereNotExists clauses

The whereExists method empowers you to craft SQL clauses for evaluating existence. When using the whereExists method, you can provide a \*DB argument, which will be assigned a query builder instance. This instance enables you to specify the query that should be enclosed within the "exists" clause:

```go
result, er := db.Table("or_user").Select("username").WhereExists(
    db.Table("or_user").Select("username").Where("user_id", ">=", int64(12345)),
).First()
```

### WhereBetween and WhereNotBetween clauses

```go
result, err := db.Table("or_user").Select("username").WhereBetween("user_id", 1111, 4444).Get()
```

```go
result, err := db.Table("or_user").Select("username").WhereNotBetween("user_id", 3333, 5555).Get()
```

### Left / Right / Inner / Left Outer Joins

```go
query := db.
		Table("or_user").
		Select("or_user.user_id, or_user.user_name, or_user.phone , or_user.modifydate, or_role.*").
		LeftJoin("or_user_role", "or_user.user_id", "=", "or_user_role.user_id").
		LeftJoin("or_role", "or_role.role_id", "=", "or_user_role.role_id")
```

### Insert

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test04()
}

func test04() {
	err := db.Table("table1").Insert(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})

	// insert returning id
	id, err := db.Table("table1").InsertGetId(map[string]interface{}{"foo": "foo foo foo", "bar": "bar bar bar", "baz": int64(123)})

	// batch insert
	err = db.Table("table1").InsertBatch([]map[string]interface{}{
		0: {"foo": "foo foo foo", "bar": "bar bar bar", "baz": 123},
		1: {"foo": "foo foo foo foo", "bar": "bar bar bar bar", "baz": 1234},
		2: {"foo": "foo foo foo foo foo", "bar": "bar bar bar bar bar", "baz": 12345},
	})
}
```

### Update

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test05()
}

func test05() {
	rows, err := db.Table("or_user").Where("user_id", ">", 3).Update(map[string]interface{}{"username": "paul"})
}
```

### Drop, Truncate and Rename

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test06()
}

func test06() {
	db.Drop("table_name")
	db.DropIfExists("table_name")
	db.Truncate("table_name")
	db.Rename("table_name1", "table_name2")
}
```

### Increment and Decrement

The query builder offers convenient functions for increasing or decreasing the value of a specified column. This feature serves as a shortcut, presenting a more expressive and concise interface than manually crafting the update statement. Both of these functions require two arguments: the target column and a second parameter to determine the amount by which the column should be incremented or decremented.

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test07()
}

func test07() {
	db.Table("or_user").Increase("hit", 3)
	db.Table("or_user").Decrease("hit", 1)
}
```

### Union and Union All

The query builder also offers a streamlined method to combine two queries using a "union" operation. For instance, you can initiate an initial query and employ the union method to unite it with a second query:

```go
package main

import (
	"fmt"
	"log"

	"github.com/sivaosorg/qb/qb"
)

var db = qb.NewQbDb(qb.NewQbConn("postgres", "host=127.0.0.1 port=5432 user=db_username dbname=db_name password='password' sslmode=disable"))

func main() {
	test08()
}

func test08() {
	union := db.Table("or_user").Select("user_id", "user_name").Union()
	// or if UNION ALL is of need
	// union := db.Table("or_user").Select("user_id", "user_name").UnionAll()
	res, err := union.Table("or_user_bk").Select("user_id", "user_name").Get()
}
```

### Transaction Mode

In transaction mode, you have the flexibility to execute arbitrary queries alongside your code, allowing for error detection and automatic rollback in case of issues, or successful commit when everything proceeds smoothly:

```go
err := db.InTransaction(func() (interface{}, error) {
    return db.Table("or_user").Select("user_name", "user_id", "phone").Get()
})
```

### Aggregates

Furthermore, the query builder offers a range of aggregate functions, including Count, Max, Min, Avg, and Sum. You can invoke any of these functions once you've constructed your query:

```go
cnt, err := db.Table(UsersTable).WHere("user_id", ">=", 11111).Count()

avg, err := db.Table(UsersTable).Avg("user_id")

mx, err := db.Table(UsersTable).Max("user_id")

mn, err := db.Table(UsersTable).Min("user_id")

sum, err := db.Table(UsersTable).Sum("user_id")
```

### Create Table

For generating a new database table, employ the CreateTable method. The Schema method necessitates two arguments. The initial argument denotes the table's name, while the second argument is an anonymous function or closure that receives a Table struct, enabling you to define the characteristics of the new table:

```go
result, err := db.Schema("public", func(table *Table) error {
    table.Increments("id")
    table.String("title", 128).Default("The quick brown fox jumped over the lazy dog").Unique("idx_ttl")
    table.SmallInt("cnt").Default(1)
    table.Integer("points").NotNull()
    table.BigInt("likes").Index("idx_likes")
    table.Text("comment").Comment("user comment").Collation("de_DE")
    table.DblPrecision("likes_to_points").Default(0.0)
    table.Char("tag", 10)
    table.DateTime("created_at", true)
    table.DateTimeTz("updated_at", true)
    table.Decimal("tax", 2, 2)
    table.TsVector("body")
    table.TsQuery("body_query")
    table.Jsonb("settings")
    table.Point("pt")
    table.Polygon("poly")
    table.TableComment("big table for big data")
	return nil
})

// to make a foreign key constraint from another table
_, err = db.Schema("tbl_ref", func(table *Table) error {
    table.Increments("id")
    table.Integer("big_tbl_id").ForeignKey("fk_idx_big_tbl_id", "big_tbl", "id").Concurrently().IfNotExists()
    table.Char("tag", 10).Index("idx_tag").Include("likes", "created_at")  // to add index on existing column just repeat stmt + index e.g.:
    table.Rename("settings", "options")
    return nil
})
```

### Add / Modify / Drop columns

The Table structure provided in the second argument of Schema can also be utilized to update existing tables, much like how it's created initially. The Change method enables you to alter existing column types or modify the attributes of columns within the table.

```go
result, err := db.Schema("or_user", func(table *Table) error {
    table.String("user_id", 128).Change()
    return nil
})
```

```go
result, err := db.Schema("or_user", func(table *Table) error {
    table.DropColumn("expired").IfExists()
    // To drop an index on the column
    table.DropIndex("idx_user_id")
    return nil
})
```

### Chunk

When dealing with a substantial number of database records, it's advisable to leverage the chunk method. This method retrieves a small portion of the results in each iteration and passes each chunk to a closure for processing.

```go
err = db.Table("or_user").Select("username").Where("user_id", "=", id).Chunk(300, func(users []map[string]interface{}) bool {
    for _, m := range users {
        if val, ok := m["user_id"];ok {
            calc += diffFormula(val.(int64))
        }
        // or you can return false here to stop running chunks
    }
    return true
})
```

## Ref

- [PostgreSQL](https://popsql.com/learn-sql/postgresql)
- [Quick tricks in PostgreSQL](https://www.freecodecamp.org/news/postgresql-tricks/)

## Contribution

Contributions, bug reports, and feature requests are welcome! Please feel free to open issues and submit pull requests.
