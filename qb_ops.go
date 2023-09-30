package qb

import (
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

// Insert inserts one row with param bindings
func (q *QbDB) Insert(data map[string]any) error {
	if q.Txn != nil {
		return q.Txn.Insert(data)
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `)`
	setCacheExecuteStmt(query)
	_, err := q.Sql().Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts one row with param bindings
func (q *QbTxn) Insert(data map[string]any) error {
	if q.Tx == nil {
		return errTransactionModeWithoutTx
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `)`
	setCacheExecuteStmt(query)
	_, err := q.Tx.Exec(query, values...)
	if err != nil {
		return err
	}
	return nil
}

// InsertGetId inserts one row with param bindings and returning id
func (q *QbDB) InsertGetId(data map[string]any) (uint64, error) {
	if q.Txn != nil {
		return q.Txn.InsertGetId(data)
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `) RETURNING id`
	setCacheExecuteStmt(query)
	var id uint64
	err := q.Sql().QueryRow(query, values...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// InsertGetId inserts one row with param bindings and returning id
func (q *QbTxn) InsertGetId(data map[string]any) (uint64, error) {
	if q.Tx == nil {
		return 0, errTransactionModeWithoutTx
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `) RETURNING id`
	setCacheExecuteStmt(query)
	var id uint64
	err := q.Tx.QueryRow(query, values...).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

// InsertBatch inserts multiple rows based on transaction
func (q *QbDB) InsertBatch(data []map[string]any) error {
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return errTableCallBeforeOp
	}
	txn, err := q.Sql().Begin()
	if err != nil {
		log.Fatal(err)
	}
	columns, values := prepareInsertBatch(data)
	stmt, err := txn.Prepare(pq.CopyIn(builder.table, columns...))
	if err != nil {
		return err
	}
	for _, value := range values {
		_, err = stmt.Exec(value...)
		if err != nil {
			return err
		}
	}
	_, err = stmt.Exec()
	if err != nil {
		return err
	}
	err = stmt.Close()
	if err != nil {
		return err
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}

// Update builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbDB) Update(data map[string]any) (int64, error) {
	if q.Txn != nil {
		return q.Txn.Update(data)
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	setVal := ""
	l := len(columns)
	for k, col := range columns {
		setVal += fmt.Sprintf("%s%s%s", col, " = ", bindings[k])
		if k < l-1 {
			setVal += ", "
		}
	}
	query := `UPDATE "` + q.Builder.table + `" SET ` + setVal
	if IsStringNotEmpty(q.Builder.from) {
		query += fmt.Sprintf("%s%s", " FROM ", q.Builder.from)
	}
	q.Builder.startBindingsAt = l + 1
	query += q.Builder.buildClauses()
	values = append(values, prepareValues(q.Builder.whereBindings)...)
	setCacheExecuteStmt(query)
	result, err := q.Sql().Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Update builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbTxn) Update(data map[string]any) (int64, error) {
	if q.Tx == nil {
		return 0, errTransactionModeWithoutTx
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	setVal := ""
	l := len(columns)
	for k, col := range columns {
		setVal += fmt.Sprintf("%s%s%s", col, " = ", bindings[k])
		if k < l-1 {
			setVal += ", "
		}
	}
	query := `UPDATE "` + q.Builder.table + `" SET ` + setVal
	if IsStringNotEmpty(q.Builder.from) {
		query += fmt.Sprintf("%s%s", " FROM ", q.Builder.from)
	}
	q.Builder.startBindingsAt = l + 1
	query += q.Builder.buildClauses()
	values = append(values, prepareValues(q.Builder.whereBindings)...)
	setCacheExecuteStmt(query)
	result, err := q.Tx.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete builds a DELETE stmt with corresponding where clause if stated
// returning affected rows
func (q *QbDB) Delete() (int64, error) {
	if q.Txn != nil {
		return q.Txn.Delete()
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	query := `DELETE FROM "` + q.Builder.table + `"`
	query += q.Builder.buildClauses()
	setCacheExecuteStmt(query)
	result, err := q.Sql().Exec(query, prepareValues(q.Builder.whereBindings)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Delete builds a DELETE stmt with corresponding where clause if stated
// returning affected rows
func (q *QbTxn) Delete() (int64, error) {
	if q.Tx == nil {
		return 0, errTransactionModeWithoutTx
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	query := `DELETE FROM "` + q.Builder.table + `"`
	query += q.Builder.buildClauses()
	setCacheExecuteStmt(query)
	result, err := q.Tx.Exec(query, prepareValues(q.Builder.whereBindings)...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Replace inserts data if conflicting row hasn't been found, else it will update an existing one
func (q *QbDB) Replace(data map[string]any, conflict string) (int64, error) {
	if q.Txn != nil {
		return q.Txn.Replace(data, conflict)
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `) ON CONFLICT(` + conflict + `) DO UPDATE SET `
	for i, v := range columns {
		// columns[i] = v + " = excluded." + v
		columns[i] = fmt.Sprintf("%s%s%s", v, " = excluded.", v)
	}
	query += strings.Join(columns, ", ")
	setCacheExecuteStmt(query)
	result, err := q.Sql().Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// Replace inserts data if conflicting row hasn't been found, else it will update an existing one
func (q *QbTxn) Replace(data map[string]any, conflict string) (int64, error) {
	if q.Tx == nil {
		return 0, errTransactionModeWithoutTx
	}
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return 0, errTableCallBeforeOp
	}
	columns, values, bindings := prepareBindings(data)
	query := `INSERT INTO "` + builder.table + `" (` + strings.Join(columns, `, `) + `) VALUES(` + strings.Join(bindings, `, `) + `) ON CONFLICT(` + conflict + `) DO UPDATE SET `
	for i, v := range columns {
		// columns[i] = v + " = excluded." + v
		columns[i] = fmt.Sprintf("%s%s%s", v, " = excluded.", v)
	}
	query += strings.Join(columns, ", ")
	setCacheExecuteStmt(query)
	result, err := q.Tx.Exec(query, values...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

// InTransaction executes fn passed as an argument in transaction mode
// if there are no results returned - txn will be rolled back, otherwise committed and returned
func (q *QbDB) InTransaction(fn func() (any, error)) error {
	txn, err := q.Sql().Begin()
	if err != nil {
		return err
	}
	// assign transaction and builder to Txn entity
	q.Txn = &QbTxn{
		Tx:      txn,
		Builder: q.Builder,
	}
	defer func() {
		// clear Txn object after commit
		q.Txn = nil
	}()
	result, err := fn()
	if err != nil {
		errTxn := txn.Rollback()
		if errTxn != nil {
			return errTxn
		}
		return err
	}
	isOk := false
	switch v := result.(type) {
	case int:
		if v > 0 {
			isOk = true
		}
	case int64:
		if v > 0 {
			isOk = true
		}
	case uint64:
		if v > 0 {
			isOk = true
		}
	case []map[string]any:
		if len(v) > 0 {
			isOk = true
		}
	case map[string]any:
		if len(v) > 0 {
			isOk = true
		}
	}

	if !isOk {
		return txn.Rollback()
	}
	err = txn.Commit()
	if err != nil {
		return err
	}
	return nil
}
