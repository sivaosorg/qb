package qb

import (
	"fmt"
	"log"
	"strings"

	"github.com/lib/pq"
)

func NewQbOps() *QbOps {
	q := &QbOps{}
	q.SetArgs(map[string]interface{}{})
	return q
}

func (q *QbOps) SetArgs(value map[string]interface{}) *QbOps {
	q.args = value
	return q
}

func (q *QbOps) SetField(fnc func() bool, field string, value interface{}) *QbOps {
	if q.args == nil {
		q.SetArgs(map[string]interface{}{})
	}
	if !fnc() {
		return q
	}
	q.args[field] = value
	return q
}

func (q *QbOps) AppendField(field string, value interface{}) *QbOps {
	q.SetField(func() bool {
		return true
	}, field, value)
	return q
}

func (q *QbOps) GetArgs() map[string]interface{} {
	return q.args
}

func (q *QbOps) Json() string {
	return JsonString(q.args)
}

// Insert inserts one row with param bindings
func (q *QbDB) Insert(data map[string]any) error {
	if len(data) == 0 {
		return fmt.Errorf("Data is empty")
	}
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

// InsertIf inserts one row with param bindings
func (q *QbDB) InsertIf(ops *QbOps) error {
	return q.Insert(ops.GetArgs())
}

// Insert inserts one row with param bindings
func (q *QbTxn) Insert(data map[string]any) error {
	if len(data) == 0 {
		return fmt.Errorf("Data is empty")
	}
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

// InsertIf inserts one row with param bindings
func (q *QbTxn) InsertIf(ops *QbOps) error {
	return q.Insert(ops.GetArgs())
}

// InsertGetId inserts one row with param bindings and returning id
func (q *QbDB) InsertGetId(data map[string]any) (uint64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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

// InsertGetIdIf inserts one row with param bindings and returning id
func (q *QbDB) InsertGetIdIf(ops *QbOps) (uint64, error) {
	return q.InsertGetId(ops.GetArgs())
}

// InsertGetId inserts one row with param bindings and returning id
func (q *QbTxn) InsertGetId(data map[string]any) (uint64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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

// InsertGetIdIf inserts one row with param bindings and returning id
func (q *QbTxn) InsertGetIdIf(ops *QbOps) (uint64, error) {
	return q.InsertGetId(ops.GetArgs())
}

// InsertBatch inserts multiple rows based on transaction
func (q *QbDB) InsertBatch(data []map[string]any) error {
	if len(data) == 0 {
		return fmt.Errorf("Data is empty")
	}
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

// InsertBatchIf inserts multiple rows based on transaction
func (q *QbDB) InsertBatchIf(ops ...*QbOps) error {
	data := []map[string]interface{}{}
	for _, v := range ops {
		data = append(data, v.GetArgs())
	}
	return q.InsertBatch(data)
}

// Update builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbDB) Update(data map[string]any) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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

// UpdateIf builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbDB) UpdateIf(ops *QbOps) (int64, error) {
	return q.Update(ops.GetArgs())
}

// Update builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbTxn) Update(data map[string]any) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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

// UpdateIf builds an UPDATE sql stmt with corresponding where/from clauses if stated
// returning affected rows
func (q *QbTxn) UpdateIf(ops *QbOps) (int64, error) {
	return q.Update(ops.GetArgs())
}

// Replace inserts data if conflicting row hasn't been found, else it will update an existing one
func (q *QbDB) Replace(data map[string]any, conflict string) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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
func (q *QbDB) ReplaceIf(ops *QbOps, conflict string) (int64, error) {
	return q.Replace(ops.GetArgs(), conflict)
}

// Replace inserts data if conflicting row hasn't been found, else it will update an existing one
func (q *QbTxn) Replace(data map[string]any, conflict string) (int64, error) {
	if len(data) == 0 {
		return 0, fmt.Errorf("Data is empty")
	}
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

func (q *QbTxn) ReplaceIf(ops *QbOps, conflict string) (int64, error) {
	return q.Replace(ops.GetArgs(), conflict)
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
