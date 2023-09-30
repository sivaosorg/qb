package qb

import "fmt"

// WhereExists constructs one builder from another to implement WHERE EXISTS sql/dml clause
func (q *QbDB) WhereExists(db *QbDB) *QbDB {
	q.Builder.whereExists = " WHERE EXISTS(" + db.Builder.buildSelect() + ")"
	return q
}

// WhereNotExists constructs one builder from another to implement WHERE NOT EXISTS sql/dml clause
func (q *QbDB) WhereNotExists(db *QbDB) *QbDB {
	q.Builder.whereExists = " WHERE NOT EXISTS(" + db.Builder.buildSelect() + ")"
	return q
}

// Where accepts left operand-operator-right operand to apply them to where clause
func (q *QbDB) Where(operand, operator string, value any) *QbDB {
	return q.buildWhere("", operand, operator, value)
}

// WhereEqIf accepts left operand-operator-right operand to apply them to where clause
func (q *QbDB) WhereEqIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.Where(operand, "=", value)
}

// WhereEqThanIf accepts left operand-operator-right operand to apply them to where clause
func (q *QbDB) WhereEqThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.Where(operand, ">=", value)
}

// WhereLessIf accepts left operand-operator-right operand to apply them to where clause
func (q *QbDB) WhereLessIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.Where(operand, "<", value)
}

// WhereLessThanIf accepts left operand-operator-right operand to apply them to where clause
func (q *QbDB) WhereLessThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.Where(operand, "<=", value)
}

// AndWhere accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhere(operand, operator string, value any) *QbDB {
	return q.buildWhere("AND", operand, operator, value)
}

// AndWhereEqIf accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhereEqIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhere(operand, "=", value)
}

// AndWhereEqThanIf accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhereEqThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhere(operand, ">=", value)
}

// AndWhereLessIf accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhereLessIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhere(operand, "<", value)
}

// AndWhereLessThanIf accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhereLessThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhere(operand, "<=", value)
}

// OrWhere accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhere(operand, operator string, value any) *QbDB {
	return q.buildWhere("OR", operand, operator, value)
}

// OrWhereEqIf accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhereEqIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhere(operand, "=", value)
}

// OrWhereEqThanIf accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhereEqThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhere(operand, ">=", value)
}

// OrWhereLessIf accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhereLessIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhere(operand, "<", value)
}

// OrWhereLessThanIf accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhereLessThanIf(fnc func() bool, operand string, value any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhere(operand, "<=", value)
}

// WhereBetween sets the clause BETWEEN 2 values
func (q *QbDB) WhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere("", column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// WhereBetweenIf sets the clause BETWEEN 2 values
func (q *QbDB) WhereBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereBetween(column, value1, value2)
}

// OrWhereBetween sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorOr, column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// OrWhereBetweenIf sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereBetween(column, value1, value2)
}

// AndWhereBetween sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorAnd, column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// AndWhereBetweenIf sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereBetween(column, value1, value2)
}

// WhereNotBetween sets the clause NOT BETWEEN 2 values
func (q *QbDB) WhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere("", column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// WhereNotBetweenIf sets the clause NOT BETWEEN 2 values
func (q *QbDB) WhereNotBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereNotBetween(column, value1, value2)
}

// OrWhereNotBetween sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorOr, column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// OrWhereNotBetweenIf sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereNotBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereNotBetween(column, value1, value2)
}

// AndWhereNotBetween sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorAnd, column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// AndWhereNotBetweenIf sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereNotBetweenIf(fnc func() bool, column string, value1, value2 any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereNotBetween(column, value1, value2)
}

// WhereRaw accepts custom string to apply it to where clause
func (q *QbDB) WhereRaw(raw string) *QbDB {
	q.Builder.where = Where + raw
	return q
}

// WhereRawIf accepts custom string to apply it to where clause
func (q *QbDB) WhereRawIf(fnc func() bool, raw string) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereRaw(raw)
}

// OrWhereRaw accepts custom string to apply it to where clause with logical OR
func (q *QbDB) OrWhereRaw(raw string) *QbDB {
	q.Builder.where += Or + raw
	return q
}

// OrWhereRawIf accepts custom string to apply it to where clause with logical OR
func (q *QbDB) OrWhereRawIf(fnc func() bool, raw string) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereRaw(raw)
}

// AndWhereRaw accepts custom string to apply it to where clause with logical OR
func (q *QbDB) AndWhereRaw(raw string) *QbDB {
	q.Builder.where += And + raw
	return q
}

// AndWhereRawIf accepts custom string to apply it to where clause with logical OR
func (q *QbDB) AndWhereRawIf(fnc func() bool, raw string) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereRaw(raw)
}

// WhereIn appends IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) WhereIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("", field, "IN", ins)
	return q
}

// WhereInIf appends IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) WhereInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereIn(field, values)
}

// WhereNotIn appends NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) WhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("", field, "NOT IN", ins)
	return q
}

// WhereNotInIf appends NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) WhereNotInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereNotIn(field, values)
}

// OrWhereIn appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) OrWhereIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("OR", field, "IN", ins)
	return q
}

// OrWhereInIf appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) OrWhereInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereIn(field, values)
}

// OrWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) OrWhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("OR", field, "NOT IN", ins)
	return q
}

// OrWhereNotInIf appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) OrWhereNotInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereNotIn(field, values)
}

// AndWhereIn appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) AndWhereIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("AND", field, "IN", ins)
	// r.buildWhere("AND", field, "IN", prepareSlice(ins))
	return q
}

// AndWhereInIf appends OR IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) AndWhereInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereIn(field, values)
}

// AndWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) AndWhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("AND", field, "NOT IN", ins)
	return q
}

// AndWhereNotInIf appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) AndWhereNotInIf(fnc func() bool, field string, values any) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereNotIn(field, values)
}

// WhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) WhereNull(field string) *QbDB {
	return q.buildWhere("", field, SqlOperatorIs, SqlSpecificValueNull)
}

// WhereNullIf appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) WhereNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereNull(field)
}

// WhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) WhereNotNull(field string) *QbDB {
	return q.buildWhere("", field, SqlOperatorIs, SqlSpecificValueNotNull)
}

// WhereNotNullIf appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) WhereNotNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.WhereNotNull(field)
}

// OrWhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) OrWhereNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorOr, field, SqlOperatorIs, SqlSpecificValueNull)
}

// OrWhereNullIf appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) OrWhereNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereNull(field)
}

// OrWhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) OrWhereNotNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorOr, field, SqlOperatorIs, SqlSpecificValueNotNull)
}

// OrWhereNotNullIf appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) OrWhereNotNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.OrWhereNotNull(field)
}

// AndWhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) AndWhereNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorAnd, field, SqlOperatorIs, SqlSpecificValueNull)
}

// AndWhereNullIf appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) AndWhereNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereNull(field)
}

// AndWhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) AndWhereNotNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorAnd, field, SqlOperatorIs, SqlSpecificValueNotNull)
}

// AndWhereNotNullIf appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) AndWhereNotNullIf(fnc func() bool, field string) *QbDB {
	if !fnc() {
		return q
	}
	return q.AndWhereNotNull(field)
}

func (q *QbDB) buildWhere(prefix, operand, operator string, value any) *QbDB {
	if IsStringNotEmpty(prefix) {
		prefix = fmt.Sprintf("%s%s%s", " ", prefix, " ")
	}
	q.Builder.whereBindings = append(q.Builder.whereBindings, map[string]any{prefix + operand + " " + operator: value})
	return q
}
