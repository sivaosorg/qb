package qb

import (
	"database/sql"
	"fmt"
	"strconv"
)

// Schema creates and/or manipulates table structure with an appropriate types/indices/comments/defaults/nulls etc
func (q *QbDB) Schema(tableName string, fn func(table *QbTable) error) (result sql.Result, err error) {
	tbl := &QbTable{tableName: tableName}
	err = fn(tbl) // run fn with Table struct passed to collect columns to []*column slice
	if err != nil {
		return nil, err
	}
	l := len(tbl.columns)
	if l > 0 {
		tableExistsOk, err := q.HasTable(DefaultSchema, tableName)
		if err != nil {
			return nil, err
		}
		if tableExistsOk { // modify tbl by adding/modifying/deleting columns/indices
			return q.modifyTable(tbl)
		}
		// create table with relative columns/indices
		return q.createTable(tbl)
	}
	return
}

// SchemaIfNotExists creates table structure if not exists with an appropriate types/indices/comments/defaults/nulls etc
func (q *QbDB) SchemaIfNotExists(tableName string, fn func(table *QbTable) error) (result sql.Result, err error) {
	tbl := &QbTable{tableName: tableName}
	err = fn(tbl) // run fn with Table struct passed to collect columns to []*column slice
	if err != nil {
		return nil, err
	}
	l := len(tbl.columns)
	if l > 0 {
		// create table with relative columns/indices
		tbl.ifExists = IfNotExists
		return q.createTable(tbl)
	}
	return
}

func (q *QbDB) createIndices(indices []string) (result sql.Result, err error) {
	if len(indices) == 0 {
		return nil, fmt.Errorf("Indices is empty")
	}
	for _, idx := range indices {
		if idx != "" {
			result, err = q.Sql().Exec(idx)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

func (q *QbDB) createComments(comments []string) (result sql.Result, err error) {
	for _, comment := range comments {
		if comment != "" {
			result, err = q.Sql().Exec(comment)
			if err != nil {
				return nil, err
			}
		}
	}
	return
}

func (q *QbTable) composeTableComment() string {
	if q.comment != nil {
		return "COMMENT ON TABLE " + q.tableName + " IS '" + *q.comment + "'"
	}
	return ""
}

// Increments creates auto incremented primary key integer column
func (q *QbTable) Increments(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeSerial, IsPrimaryKey: true})
	return q
}

// BigIncrements creates auto incremented primary key big integer column
func (q *QbTable) BigIncrements(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeBigSerial, IsPrimaryKey: true})
	return q
}

// SmallInt creates small integer column
func (q *QbTable) SmallInt(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeSmallInt})
	return q
}

// Integer creates an integer column
func (q *QbTable) Integer(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeInt})
	return q
}

// BigInt creates big integer column
func (q *QbTable) BigInt(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeBigInt})
	return q
}

// String creates varchar(len) column
func (q *QbTable) String(column string, len uint64) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: qbColType(TypeVarCharacter + "(" + strconv.FormatUint(len, 10) + ")")})
	return q
}

// Char creates char(len) column
func (q *QbTable) Char(column string, len uint64) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: qbColType(TypeChar + "(" + strconv.FormatUint(len, 10) + ")")})
	return q
}

// Boolean creates boolean type column
func (q *QbTable) Boolean(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeBoolean})
	return q
}

// Text	creates text type column
func (q *QbTable) Text(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeText})
	return q
}

// DblPrecision	creates dbl precision type column
func (q *QbTable) DblPrecision(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeDblPrecision})
	return q
}

// Numeric creates exact, user-specified precision number
func (q *QbTable) Numeric(column string, precision, scale uint64) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: qbColType(TypeNumeric + "(" + strconv.FormatUint(precision, 10) + ", " + strconv.FormatUint(scale, 10) + ")")})
	return q
}

// Decimal alias for Numeric as for PostgreSQL they are the same
func (q *QbTable) Decimal(column string, precision, scale uint64) *QbTable {
	return q.Numeric(column, precision, scale)
}

// NotNull sets the last column to not null
func (q *QbTable) NotNull() *QbTable {
	q.columns[len(q.columns)-1].IsNotNull = true
	return q
}

// Collation sets the last column to specified collation
func (q *QbTable) Collation(collation string) *QbTable {
	q.columns[len(q.columns)-1].Collation = &collation
	return q
}

// Default sets the default column value
func (q *QbTable) Default(value interface{}) *QbTable {
	v := transform2String(value)
	q.columns[len(q.columns)-1].Default = &v
	return q
}

// Comment sets the column comment
func (q *QbTable) Comment(comment string) *QbTable {
	q.columns[len(q.columns)-1].Comment = &comment
	return q
}

// TableComment sets the comment for table
func (q *QbTable) TableComment(comment string) {
	q.comment = &comment
}

// Index sets the last column to btree index
func (q *QbTable) Index(indexName string) *QbTable {
	q.columns[len(q.columns)-1].IdxName = indexName
	q.columns[len(q.columns)-1].IsIndex = true
	return q
}

// Unique sets the last column to unique index
func (q *QbTable) Unique(indexName string) *QbTable {
	q.columns[len(q.columns)-1].IdxName = indexName
	q.columns[len(q.columns)-1].IsUnique = true
	return q
}

// ForeignKey sets the last column to reference rfcTbl on onCol with idxName foreign key index
func (q *QbTable) ForeignKey(indexName, referTable, onColumn string) *QbTable {
	key := AlterTable + q.tableName + " ADD CONSTRAINT " + indexName + " FOREIGN KEY (" + q.columns[len(q.columns)-1].Name + ") REFERENCES " + referTable + " (" + onColumn + ")"
	q.columns[len(q.columns)-1].ForeignKey = &key
	return q
}

func (q *QbTable) Concurrently() *QbTable {
	q.columns[len(q.columns)-1].IsIdxConcurrent = true
	return q
}

func (q *QbTable) Include(columns ...string) *QbTable {
	q.columns[len(q.columns)-1].Includes = columns
	return q
}

// Date	creates date column with an ability to set current_date as default value
func (q *QbTable) Date(column string, isDefault bool) *QbTable {
	q.columns = append(q.columns, buildDateTime(column, TypeDate, CurrentDate, isDefault))
	return q
}

// Time creates time column with an ability to set current_time as default value
func (q *QbTable) Time(column string, isDefault bool) *QbTable {
	q.columns = append(q.columns, buildDateTime(column, TypeTime, CurrentTime, isDefault))
	return q
}

// DateTime creates date.time column with an ability to set NOW() as default value
func (q *QbTable) DateTime(column string, isDefault bool) *QbTable {
	q.columns = append(q.columns, buildDateTime(column, TypeDateTime, CurrentDateTime, isDefault))
	return q
}

// DateTimeTz creates date.time column with an ability to set NOW() as default value + time zone support
func (q *QbTable) DateTimeTz(column string, isDefault bool) *QbTable {
	q.columns = append(q.columns, buildDateTime(column, TypeDateTimeTz, CurrentDateTime, isDefault))
	return q
}

// TsVector creates tsvector typed column
func (q *QbTable) TsVector(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeTsVector})
	return q
}

// TsQuery creates tsquery typed column
func (q *QbTable) TsQuery(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeTsQuery})
	return q
}

// Json creates json text typed column
func (q *QbTable) Json(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeJson})
	return q
}

// Jsonb creates jsonb typed column
func (q *QbTable) Jsonb(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypeJsonb})
	return q
}

// Point creates point geometry typed column
func (q *QbTable) Point(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypePoint})
	return q
}

// Polygon creates point geometry typed column
func (q *QbTable) Polygon(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, ColumnType: TypePolygon})
	return q
}

// Change the column type/length/nullable etc options
func (q *QbTable) Change() {
	q.columns[len(q.columns)-1].IsModify = true
}

// IfNotExists add column/index if not exists
func (q *QbTable) IfNotExists() *QbTable {
	q.columns[len(q.columns)-1].IfExists = IfNotExists
	return q
}

// IfExists drop column/index if exists
func (q *QbTable) IfExists() *QbTable {
	q.columns[len(q.columns)-1].IfExists = IfExists
	return q
}

// Rename the column "from" to the "to"
func (q *QbTable) Rename(from, to string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: from, RenameTo: &to, IsModify: true})
	return q
}

// DropColumn the column named colNm in this table context
func (q *QbTable) DropColumn(column string) *QbTable {
	q.columns = append(q.columns, &qbColumn{Name: column, IsDrop: true})
	return q
}

// DropIndex the column named idxNm in this table context
func (q *QbTable) DropIndex(indexName string) *QbTable {
	q.columns = append(q.columns, &qbColumn{IdxName: indexName, IsDrop: true, IsIndex: true})
	return q
}

// createTable create table with relative columns/indices
func (q *QbDB) createTable(t *QbTable) (result sql.Result, err error) {
	l := len(t.columns)
	var indices []string
	var comments []string
	query := "CREATE TABLE " + applyExistence(t.ifExists) + t.tableName + "("
	for k, col := range t.columns {
		query += composeColumn(col)
		if k < l-1 {
			query += ","
		}
		indices = append(indices, composeIndex(t.tableName, col))
		comments = append(comments, composeComment(t.tableName, col))
	}
	query += ")"

	result, err = q.Sql().Exec(query)
	if err != nil {
		return nil, err
	}
	// create indices
	_, err = q.createIndices(indices)
	if err != nil {
		return nil, err
	}
	// create comments
	comments = append(comments, t.composeTableComment())
	_, err = q.createComments(comments)
	if err != nil {
		return nil, err
	}
	return
}

// adds, modifies or deletes column
func (q *QbDB) modifyTable(t *QbTable) (result sql.Result, err error) {
	l := len(t.columns)
	var indices []string
	var comments []string
	query := ""
	for key, column := range t.columns {
		if column.IsModify {
			column.Operator = Modify
			if column.RenameTo != nil {
				column.Operator = Rename
			}
			query += composeModifyColumn(t.tableName, column)
		} else if column.IsDrop {
			query += composeDrop(t.tableName, column)
		} else {
			isCol, _ := q.HasColumns(DefaultSchema, t.tableName, column.Name) // create new column/comment/index or just add comments indices
			if !isCol {
				query += composeAddColumn(t.tableName, column)
			}
			indices = append(indices, composeIndex(t.tableName, column))
			comments = append(comments, composeComment(t.tableName, column))
		}
		if key < l-1 {
			query += SemiColon
		}
	}
	result, err = q.Sql().Exec(query)
	if err != nil {
		return nil, err
	}
	// create indices
	_, err = q.createIndices(indices)
	if err != nil {
		return nil, err
	}
	// create comments
	_, err = q.createComments(comments)
	if err != nil {
		return nil, err
	}
	return
}
