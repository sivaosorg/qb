package qb

import "fmt"

// list all join types
const (
	JoinInner     = "INNER"
	JoinLeft      = "LEFT"
	JoinRight     = "RIGHT"
	JoinFull      = "FULL"
	JoinFullOuter = "FULL OUTER"
	Where         = " WHERE "
	And           = " AND "
	Or            = " OR "
)

// list all sql operators
const (
	SqlOperatorBetween    = "BETWEEN"
	SqlOperatorNotBetween = "NOT BETWEEN"
	SqlOperatorIs         = "IS"
	SqlOperatorAnd        = "AND"
	SqlOperatorOr         = "OR"
)

// list all invalid types
const (
	SqlSpecificValueNull    = "NULL"
	SqlSpecificValueNotNull = "NOT NULL"
)

// math operation
const (
	plusSign  = "+"
	minusSign = "-"
)

const (
	IfExistsUndeclared = iota
	IfExists
	IfNotExists
)

// column types
const (
	TypeSerial       = "SERIAL"
	TypeBigSerial    = "BIGSERIAL"
	TypeSmallInt     = "SMALLINT"
	TypeInt          = "INTEGER"
	TypeBigInt       = "BIGINT"
	TypeBoolean      = "BOOLEAN"
	TypeText         = "TEXT"
	TypeVarCharacter = "VARCHAR"
	TypeChar         = "CHAR"
	TypeDate         = "DATE"
	TypeTime         = "TIME"
	TypeDateTime     = "TIMESTAMP"
	TypeDateTimeTz   = "TIMESTAMPTZ"
	CurrentDate      = "CURRENT_DATE"
	CurrentTime      = "CURRENT_TIME"
	CurrentDateTime  = "NOW()"
	TypeDblPrecision = "DOUBLE PRECISION"
	TypeNumeric      = "NUMERIC"
	TypeTsVector     = "TSVECTOR"
	TypeTsQuery      = "TSQUERY"
	TypeJson         = "JSON"
	TypeJsonb        = "JSONB"
	TypePoint        = "POINT"
	TypePolygon      = "POLYGON"
)

// specific for PostgreSQL driver and SQL std
const (
	DefaultSchema  = "public"
	SemiColon      = ";"
	AlterTable     = "ALTER TABLE "
	Add            = " ADD "
	Modify         = " ALTER "
	Drop           = " DROP "
	Rename         = " RENAME "
	IfExistsExp    = " IF EXISTS "
	IfNotExistsExp = " IF NOT EXISTS "
	Concurrently   = " CONCURRENTLY "
	Constraint     = " CONSTRAINT "
)

var (
	errTableCallBeforeOp        = fmt.Errorf("sql: there was no Table() call with table name set")
	errTransactionModeWithoutTx = fmt.Errorf("sql: there was no *sql.Tx object set properly")
)

var (
	// Cache Execute Stmt will be stored value for the query SQL has been called before
	cacheExecuteStmt string = ""
)
