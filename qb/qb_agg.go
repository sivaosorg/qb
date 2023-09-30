package qb

// Count counts resulting rows based on clause
func (q *QbDB) Count() (countRows int64, err error) {
	builder := q.Builder
	builder.columns = []string{"COUNT(*)"}
	query := builder.buildSelect()
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&countRows)
	return
}

// Avg calculates average for specified column
func (q *QbDB) Avg(column string) (avg float64, err error) {
	builder := q.Builder
	builder.columns = []string{"AVG(" + column + ")"}
	query := builder.buildSelect()
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&avg)
	return
}

// Min calculates minimum for specified column
func (q *QbDB) Min(column string) (min float64, err error) {
	builder := q.Builder
	builder.columns = []string{"MIN(" + column + ")"}
	query := builder.buildSelect()
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&min)
	return
}

// Max calculates maximum for specified column
func (q *QbDB) Max(column string) (max float64, err error) {
	builder := q.Builder
	builder.columns = []string{"MAX(" + column + ")"}
	query := builder.buildSelect()
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&max)
	return
}

// Sum calculates sum for specified column
func (q *QbDB) Sum(column string) (max float64, err error) {
	builder := q.Builder
	builder.columns = []string{"SUM(" + column + ")"}
	query := builder.buildSelect()
	err = q.Sql().QueryRow(query, prepareValues(q.Builder.whereBindings)...).Scan(&max)
	return
}
