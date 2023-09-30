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

// AndWhere accepts left operand-operator-right operand to apply them to where clause
// with AND logical operator
func (q *QbDB) AndWhere(operand, operator string, value any) *QbDB {
	return q.buildWhere("AND", operand, operator, value)
}

// OrWhere accepts left operand-operator-right operand to apply them to where clause
// with OR logical operator
func (q *QbDB) OrWhere(operand, operator string, value any) *QbDB {
	return q.buildWhere("OR", operand, operator, value)
}

func (q *QbDB) buildWhere(prefix, operand, operator string, value any) *QbDB {
	if IsStringNotEmpty(prefix) {
		prefix = fmt.Sprintf("%s%s%s", " ", prefix, " ")
	}
	q.Builder.whereBindings = append(q.Builder.whereBindings, map[string]any{prefix + operand + " " + operator: value})
	return q
}

// WhereBetween sets the clause BETWEEN 2 values
func (q *QbDB) WhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere("", column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// OrWhereBetween sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorOr, column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// AndWhereBetween sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorAnd, column, SqlOperatorBetween, transform2String(value1)+And+transform2String(value2))
}

// WhereNotBetween sets the clause NOT BETWEEN 2 values
func (q *QbDB) WhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere("", column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// OrWhereNotBetween sets the clause OR BETWEEN 2 values
func (q *QbDB) OrWhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorOr, column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// AndWhereNotBetween sets the clause AND BETWEEN 2 values
func (q *QbDB) AndWhereNotBetween(column string, value1, value2 any) *QbDB {
	return q.buildWhere(SqlOperatorAnd, column, SqlOperatorNotBetween, transform2String(value1)+And+transform2String(value2))
}

// WhereRaw accepts custom string to apply it to where clause
func (q *QbDB) WhereRaw(raw string) *QbDB {
	q.Builder.where = Where + raw
	return q
}

// OrWhereRaw accepts custom string to apply it to where clause with logical OR
func (q *QbDB) OrWhereRaw(raw string) *QbDB {
	q.Builder.where += Or + raw
	return q
}

// AndWhereRaw accepts custom string to apply it to where clause with logical OR
func (q *QbDB) AndWhereRaw(raw string) *QbDB {
	q.Builder.where += And + raw
	return q
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

// WhereNotIn appends NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) WhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("", field, "NOT IN", ins)
	return q
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

// OrWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) OrWhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("OR", field, "NOT IN", ins)
	return q
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

// AndWhereNotIn appends OR NOT IN (val1, val2, val3...) stmt to WHERE clause
func (q *QbDB) AndWhereNotIn(field string, values any) *QbDB {
	ins, err := Interface2Slice(values)
	if err != nil {
		return nil
	}
	q.buildWhere("AND", field, "NOT IN", ins)
	return q
}

// WhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) WhereNull(field string) *QbDB {
	return q.buildWhere("", field, SqlOperatorIs, SqlSpecificValueNull)
}

// WhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) WhereNotNull(field string) *QbDB {
	return q.buildWhere("", field, SqlOperatorIs, SqlSpecificValueNotNull)
}

// OrWhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) OrWhereNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorOr, field, SqlOperatorIs, SqlSpecificValueNull)
}

// OrWhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) OrWhereNotNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorOr, field, SqlOperatorIs, SqlSpecificValueNotNull)
}

// AndWhereNull appends fieldName IS NULL stmt to WHERE clause
func (q *QbDB) AndWhereNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorAnd, field, SqlOperatorIs, SqlSpecificValueNull)
}

// AndWhereNotNull appends fieldName IS NOT NULL stmt to WHERE clause
func (q *QbDB) AndWhereNotNull(field string) *QbDB {
	return q.buildWhere(SqlOperatorAnd, field, SqlOperatorIs, SqlSpecificValueNotNull)
}
