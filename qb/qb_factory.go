package qb

import (
	"fmt"
)

// Get builds all sql statements chained before and executes query collecting data to the slice
func (q *QbDB) Get() ([]map[string]any, error) {
	builder := q.Builder
	if IsStringEmpty(builder.table) {
		return nil, errTableCallBeforeOp
	}
	query := ""
	if len(builder.union) > 0 { // got union - need different logic to glue
		for _, subUnion := range builder.union {
			if IsStringEmpty(subUnion) {
				continue
			}
			// query += subUnion + " UNION "
			query += fmt.Sprintf("%s%s", subUnion, " UNION ")
			if builder.isUnionAll {
				// query += "ALL "
				query += fmt.Sprintf("%s", "ALL ")
			}
		}
		query += builder.buildSelect()
		// clean union (all) after ensuring selects are built
		q.Builder.union = []string{}
		q.Builder.isUnionAll = false
	} else {
		query = builder.buildSelect()
	}
	rows, err := q.Sql().Query(query, prepareValues(q.Builder.whereBindings)...)
	if err != nil {
		return nil, err
	}
	columns, _ := rows.Columns()
	count := len(columns)
	values := make([]any, count)
	valuesCount := make([]any, count)
	// collecting data from struct with fields
	var response []map[string]any

	for rows.Next() {
		collect := make(map[string]any, count)
		for i := range columns {
			valuesCount[i] = &values[i]
		}
		err := rows.Scan(valuesCount...)
		if err != nil {
			return nil, err
		}
		for i, col := range columns {
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				collect[col] = string(b)
			} else {
				collect[col] = val
			}
		}
		response = append(response, collect)
	}
	return response, nil
}

// First getting the 1st row of query
func (q *QbDB) First() (map[string]interface{}, error) {
	result, err := q.Get()
	if err != nil {
		return nil, err
	}
	if len(result) > 0 {
		return result[0], nil
	}
	return nil, fmt.Errorf("no records were produced by query: %s", q.GetQuery())
}

// Value gets the value of column in first query resulting row
func (q *QbDB) Value(column string) (value interface{}, err error) {
	q.Select(column)
	result, err := q.First()
	if err != nil {
		return
	}
	if val, ok := result[column]; ok {
		return val, err
	}
	return
}

// Find retrieves a single row by it's id column value
func (q *QbDB) Find(id uint64) (map[string]interface{}, error) {
	return q.Where("id", "=", id).First()
}

// Pull getting values of a particular column and place them into slice
func (q *QbDB) Pull(column string) (value []interface{}, err error) {
	result, err := q.Get()
	if err != nil {
		return nil, err
	}
	value = make([]interface{}, len(result))
	for k, m := range result {
		value[k] = m[column]
	}
	return
}

// PullParticulars getting values of a particular key/value columns and place them into map
func (q *QbDB) PullParticulars(columnKey, columnValue string) (value []map[interface{}]interface{}, err error) {
	result, err := q.Get()
	if err != nil {
		return nil, err
	}
	value = make([]map[interface{}]interface{}, len(result))
	for k, m := range result {
		value[k] = make(map[interface{}]interface{})
		value[k][m[columnKey]] = m[columnValue]
	}
	return
}
