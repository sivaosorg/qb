package qb

import "database/sql"

type qbColType string

type QbConn struct {
	db *sql.DB `json:"-"`
}

type QbDB struct {
	Builder *qbBuilder `json:"-"`
	Conn    *QbConn    `json:"-"`
	Txn     *QbTxn     `json:"-"`
}

type QbTxn struct {
	Tx      *sql.Tx    `json:"-"`
	Builder *qbBuilder `json:"-"`
}

// QbTable is the type for operations on table schema
type QbTable struct {
	ifExists  uint        `json:"-"`
	columns   []*qbColumn `json:"-"`
	tableName string      `json:"-"`
	comment   *string     `json:"-"`
}

type qbBuilder struct {
	whereBindings   []map[string]any
	startBindingsAt int
	where           string
	table           string
	from            string
	join            []string
	orderBy         []map[string]string
	orderByRaw      *string
	groupBy         string
	having          string
	columns         []string
	union           []string
	isUnionAll      bool
	offset          int64
	limit           int64
	lockForUpdate   *string
	whereExists     string
}

type qbColumn struct {
	IsNotNull       bool
	IsPrimaryKey    bool
	IsIndex         bool
	IsIdxConcurrent bool
	IsUnique        bool
	IsDrop          bool
	IsModify        bool
	IfExists        uint
	Includes        []string
	Name            string
	RenameTo        *string
	ColumnType      qbColType
	Default         *string
	ForeignKey      *string
	IdxName         string
	Comment         *string
	Collation       *string
	Operator        string
}
