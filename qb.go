package qb

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	_ "github.com/lib/pq" // for PostgreSQL driver
)

func NewQbConn(driverName, dataSourceName string) *QbConn {
	db, err := sql.Open(driverName, dataSourceName)
	if err != nil {
		log.Fatalln(err)
	}
	return &QbConn{db: db}
}

func NewQbConnWith(db *sql.DB) *QbConn {
	return &QbConn{db: db}
}

func newBuilder() *qbBuilder {
	return &qbBuilder{
		columns: []string{"*"},
	}
}

func (q *QbDB) Sql() *sql.DB {
	return q.Conn.db
}

func NewQbDb(q *QbConn) *QbDB {
	b := newBuilder()
	return &QbDB{Builder: b, Conn: q}
}

// Table appends table name to sql query
func (q *QbDB) Table(table string) *QbDB {
	q.reset() // reset before constructing again
	q.Builder.table = table
	return q
}

// resets all builder elements to prepare them for next round
func (q *QbDB) reset() {
	q.Builder.table = ""
	q.Builder.columns = []string{"*"}
	q.Builder.where = ""
	q.Builder.whereBindings = make([]map[string]any, 0)
	q.Builder.groupBy = ""
	q.Builder.having = ""
	q.Builder.orderBy = make([]map[string]string, 0)
	q.Builder.offset = 0
	q.Builder.limit = 0
	q.Builder.join = []string{}
	q.Builder.from = ""
	q.Builder.lockForUpdate = nil
	q.Builder.whereExists = ""
	q.Builder.orderByRaw = nil
	q.Builder.startBindingsAt = 1
	if len(q.Builder.union) == 0 {
		q.Builder.union = []string{}
	}
}

// Select accepts columns to select from a table
func (q *QbDB) Select(args ...string) *QbDB {
	q.Builder.columns = []string{}
	q.Builder.columns = append(q.Builder.columns, args...)
	return q
}

// OrderBy adds ORDER BY expression to SQL stmt
func (q *QbDB) OrderBy(column string, direction string) *QbDB {
	q.Builder.orderBy = append(q.Builder.orderBy, map[string]string{column: direction})
	return q
}

// OrderByRaw adds ORDER BY raw expression to SQL stmt
func (q *QbDB) OrderByRaw(exp string) *QbDB {
	q.Builder.orderByRaw = &exp
	return q
}

// Important
// Don't use for amount big table
// InRandomOrder add ORDER BY random() - note be cautious on big data-tables it can lead to slowing down perf
func (q *QbDB) InRandomOrder() *QbDB {
	q.OrderByRaw("random()")
	return q
}

// GroupBy adds GROUP BY expression to SQL stmt
func (q *QbDB) GroupBy(expr string) *QbDB {
	q.Builder.groupBy = expr
	return q
}

// Having similar to Where but used with GroupBy to apply over the grouped results
func (q *QbDB) Having(operand, operator string, value any) *QbDB {
	q.Builder.having = operand + " " + operator + " " + transform2String(value)
	return q
}

// HavingRaw accepts custom string to apply it to having clause
func (q *QbDB) HavingRaw(raw string) *QbDB {
	q.Builder.having = raw
	return q
}

// OrHavingRaw accepts custom string to apply it to having clause with logical OR
func (q *QbDB) OrHavingRaw(raw string) *QbDB {
	q.Builder.having += Or + raw
	return q
}

// AndHavingRaw accepts custom string to apply it to having clause with logical OR
func (q *QbDB) AndHavingRaw(raw string) *QbDB {
	q.Builder.having += And + raw
	return q
}

// AddSelect accepts additional columns to select from a table
func (q *QbDB) AddSelect(args ...string) *QbDB {
	q.Builder.columns = append(q.Builder.columns, args...)
	return q
}

// SelectRaw accepts custom string to select from a table
func (q *QbDB) SelectRaw(raw string) *QbDB {
	q.Builder.columns = []string{raw}
	return q
}

// Union joins multiple queries omitting duplicate records
func (q *QbDB) Union() *QbDB {
	q.Builder.union = append(q.Builder.union, q.Builder.buildSelect())
	return q
}

// UnionAll joins multiple queries to select all rows from both tables with duplicate
func (q *QbDB) UnionAll() *QbDB {
	q.Union()
	q.Builder.isUnionAll = true
	return q
}

// Offset accepts offset to start slicing results from
func (q *QbDB) Offset(value int64) *QbDB {
	q.Builder.offset = value
	return q
}

// Limit accepts limit to end slicing results to
func (q *QbDB) Limit(value int64) *QbDB {
	q.Builder.limit = value
	return q
}

// Page for pagination
func (q *QbDB) Page(value int64) *QbDB {
	if value < 0 {
		log.Fatalf("Invalid page: %v", value)
	}
	if value > 0 {
		value = value - 1
	}
	q.Builder.page = value
	return q
}

// Size for pagination
func (q *QbDB) Size(value int64) *QbDB {
	if value < 0 {
		log.Fatalf("Invalid size(limit): %v", value)
	}
	q.Builder.size = value
	q.Limit(value)
	q.Offset(q.Builder.page * value)
	return q
}

// Drop drops >=1 tables
func (q *QbDB) Drop(tables string) (sql.Result, error) {
	query := fmt.Sprintf("%s%s", "DROP TABLE ", tables)
	setCacheExecuteStmt(query)
	return q.Sql().Exec(query)
}

// Truncate clears >=1 tables
func (q *QbDB) Truncate(tables string) (sql.Result, error) {
	query := fmt.Sprintf("%s%s", "TRUNCATE ", tables)
	setCacheExecuteStmt(query)
	return q.Sql().Exec(query)
}

// DropIfExists drops >=1 tables if they are existent
func (q *QbDB) DropIfExists(tables ...string) (result sql.Result, err error) {
	for _, table := range tables {
		result, err = q.Sql().Exec(fmt.Sprintf("%s%s%s", "DROP TABLE", IfExistsExp, table))
	}
	return result, err
}

// Rename renames from - to new table name
func (q *QbDB) Rename(from, to string) (sql.Result, error) {
	query := fmt.Sprintf("%s%s%s%s", "ALTER TABLE ", from, " RENAME TO ", to)
	setCacheExecuteStmt(query)
	return q.Sql().Exec(query)
}

// From prepares sql stmt to set data from another table, ex.:
// UPDATE employees SET sales_count = sales_count + 1 FROM accounts
func (q *QbDB) From(fromTable string) *QbDB {
	q.Builder.from = fromTable
	return q
}

// LockForUpdate locks table/row
func (q *QbDB) LockForUpdate() *QbDB {
	str := " FOR UPDATE"
	q.Builder.lockForUpdate = &str
	return q
}

// PrintQuery prints raw sql to stdout
func (q *QbDB) PrintQuery() {
	// log.SetOutput(os.Stdout)
	log.Println(q.Builder.buildSelect())
}

// GetQuery return raw sql to stdout
func (q *QbDB) GetQuery() string {
	return q.Builder.buildSelect()
}

// PrintQueryWithExit prints raw sql to stdout and exit
func (q *QbDB) PrintQueryWithExit() {
	q.PrintQuery()
	os.Exit(0)
}

// HasTable determines whether table exists in particular schema
func (q *QbDB) HasTable(schema, table string) (tblExists bool, err error) {
	query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM pg_tables WHERE  schemaname = '%s' AND tablename = '%s')", schema, table)
	setCacheExecuteStmt(query)
	err = q.Sql().QueryRow(query).Scan(&tblExists)
	return
}

// HasColumns checks whether those cols exists in a particular schema/table
func (q *QbDB) HasColumns(schema, table string, columns ...string) (colsExists bool, err error) {
	andColumns := ""
	for _, v := range columns { // todo: find a way to check columns in 1 query
		andColumns = " AND column_name = '" + v + "'"
		query := fmt.Sprintf("SELECT EXISTS (SELECT 1 FROM information_schema.columns WHERE table_schema='%s' AND table_name='%s'"+andColumns+")", schema, table)
		setCacheExecuteStmt(query)
		err = q.Sql().QueryRow(query).Scan(&colsExists)
		if !colsExists { // if at least once col doesn't exist - return false, nil
			return
		}
	}
	return
}

// Exists checks whether conditional rows are existing (returns true) or not (returns false)
func (q *QbDB) Exists() (ok bool, err error) {
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return false, errTableCallBeforeOp
	}
	query := `SELECT EXISTS(SELECT 1 FROM "` + builder.table + `" ` + builder.buildClauses() + `)`
	setCacheExecuteStmt(query)
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&ok)
	return
}

// DoesNotExists an inverse of Exists
func (q *QbDB) DoesNotExists() (bool, error) {
	ok, err := q.Exists()
	if err != nil {
		return false, err
	}
	return !ok, nil
}

// Increase column on passed value
func (q *QbDB) Increase(column string, value uint64) (int64, error) {
	return q.increaseAndDecrease(column, plusSign, value)
}

// Decrease column on passed value
func (q *QbDB) Decrease(column string, value uint64) (int64, error) {
	return q.increaseAndDecrease(column, minusSign, value)
}

// Chunk run queries by chinks by passing user-land function with an ability to stop execution when needed
// by returning false and proceed to execute queries when return true
func (q *QbDB) Chunk(amount int64, fn func(rows []map[string]interface{}) bool) error {
	columns := q.Builder.columns
	count, err := q.Count()
	if err != nil {
		return err
	}
	q.Builder.columns = columns
	if amount <= 0 {
		return fmt.Errorf("chunk can't be <= 0, your chunk is: %d", amount)
	}
	if count < amount {
		result, err := q.Get()
		if err != nil {
			return err
		}
		fn(result) // execute all resulting records
		return nil
	}
	// executing chunks amount < cnt
	c := int64(math.Ceil(float64(count / amount)))
	var i int64
	for i = 0; i < c; i++ {
		rows, err := q.Offset(i * amount).Limit(amount).Get() // by 100 rows from 100 x n
		if err != nil {
			return err
		}
		result := fn(rows)
		if !result { // stop an execution when false returned by user
			break
		}
	}
	return nil
}

// buildSelect constructs a query for select statement
func (q *qbBuilder) buildSelect() string {
	query := `SELECT ` + strings.Join(q.columns, `, `) + ` FROM ` + q.table
	v := fmt.Sprintf("%s%s", query, q.buildClauses())
	setCacheExecuteStmt(v)
	return v
}

// builds query string clauses
func (q *qbBuilder) buildClauses() string {
	clauses := ""
	for _, j := range q.join {
		clauses += j
	}
	// build where clause
	if len(q.whereBindings) > 0 {
		clauses += composeWhere(q.whereBindings, q.startBindingsAt)
	} else {
		clauses += q.where // std without bindings todo: change all to bindings
	}
	if IsStringNotEmpty(q.groupBy) {
		// clauses += " GROUP BY " + r.groupBy
		clauses += fmt.Sprintf("%s%s", " GROUP BY ", q.groupBy)
	}
	if IsStringNotEmpty(q.having) {
		// clauses += " HAVING " + r.having
		clauses += fmt.Sprintf("%s%s", " HAVING ", q.having)
	}
	clauses += composeOrderBy(q.orderBy, q.orderByRaw)
	if q.limit > 0 {
		// clauses += " LIMIT " + strconv.FormatInt(r.limit, 10)
		clauses += fmt.Sprintf("%s%s", " LIMIT ", strconv.FormatInt(q.limit, 10))
	}
	if q.offset > 0 {
		// clauses += " OFFSET " + strconv.FormatInt(r.offset, 10)
		clauses += fmt.Sprintf("%s%s", " OFFSET ", strconv.FormatInt(q.offset, 10))
	}
	if q.lockForUpdate != nil {
		clauses += *q.lockForUpdate
	}
	return clauses
}

// increments or decrements depending on sign
func (q *QbDB) increaseAndDecrease(column, sign string, on uint64) (int64, error) {
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	query := `UPDATE "` + q.Builder.table + `" SET ` + column + ` = ` + column + sign + strconv.FormatUint(on, 10)
	setCacheExecuteStmt(query)
	result, err := q.Sql().Exec(query)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (q *QbDB) GetRawSQL() string {
	return getCacheExecuteStmt()
}
