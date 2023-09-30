package qb

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"strings"
)

func IsStringEmpty(s string) bool {
	return len(s) == 0 ||
		s == "" ||
		strings.TrimSpace(s) == "" ||
		len(strings.TrimSpace(s)) == 0
}

func IsStringNotEmpty(s string) bool {
	return !IsStringEmpty(s)
}

func JsonString(data interface{}) string {
	s, ok := data.(string)
	if ok {
		return s
	}
	result, err := json.Marshal(data)
	if err != nil {
		log.Printf(err.Error())
		return ""
	}
	return string(result)
}

func Interface2Slice(slice interface{}) ([]interface{}, error) {
	var err error
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		err = errors.New("interfaceToSlice() given a non-slice type")
	}
	v := make([]interface{}, s.Len())
	for i := 0; i < s.Len(); i++ {
		v[i] = s.Index(i).Interface()
	}
	return v, err
}

func prepareValues(values []map[string]any) []any {
	var result []any
	for _, v := range values {
		_, val, _ := prepareBindings(v)
		result = append(result, val...)
	}
	return result
}

func prepareValue(value any) []any {
	var values []any
	switch v := value.(type) {
	case string:
		//if where { // todo: left comments for further exploration, probably incorrect behavior for pg driver
		//	values = append(values, "'"+v+"'")
		//} else {
		values = append(values, v)
		//}
	case int:
		values = append(values, strconv.FormatInt(int64(v), 10))
	case float64:
		values = append(values, fmt.Sprintf("%g", v))
	case int64:
		values = append(values, strconv.FormatInt(v, 10))
	case uint64:
		values = append(values, strconv.FormatUint(v, 10))
	case []any:
		for _, vi := range v {
			values = append(values, prepareValue(vi)...)
		}
	case nil:
		values = append(values, nil)
	}
	return values
}

// prepareBindings prepares slices to split in favor of INSERT sql statement
func prepareBindings(data map[string]any) (columns []string, values []any, bindings []string) {
	i := 1
	for column, value := range data {
		if strings.Contains(column, SqlOperatorIs) || strings.Contains(column, SqlOperatorBetween) {
			continue
		}
		columns = append(columns, column)
		pValues := prepareValue(value)
		if len(pValues) > 0 {
			values = append(values, pValues...)
			for range pValues {
				bindings = append(bindings, fmt.Sprintf("%s%s", "$", strconv.FormatInt(int64(i), 10)))
				i++
			}
		}
	}
	return
}

// prepareInsertBatch prepares slices to split in favor of INSERT sql statement
func prepareInsertBatch(data []map[string]any) (columns []string, values [][]any) {
	values = make([][]any, len(data))
	columnToIdx := make(map[string]int)
	i := 0
	for k, v := range data {
		values[k] = make([]any, len(v))
		for column, value := range v {
			if k == 0 {
				columns = append(columns, column)
				// todo: don't know yet how to match them explicitly (it is bad idea, but it works well now)
				columnToIdx[column] = i
				i++
			}
			switch casted := value.(type) {
			case string:
				values[k][columnToIdx[column]] = casted
			case int:
				values[k][columnToIdx[column]] = strconv.FormatInt(int64(casted), 10)
			case float64:
				values[k][columnToIdx[column]] = fmt.Sprintf("%g", casted)
			case int64:
				values[k][columnToIdx[column]] = strconv.FormatInt(casted, 10)
			case uint64:
				values[k][columnToIdx[column]] = strconv.FormatUint(casted, 10)
			}
		}
	}
	return
}

func transform2String(value any) string {
	switch v := value.(type) {
	case string:
		return "'" + v + "'"
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case uint64:
		return strconv.FormatUint(v, 10)
	case float64:
		return fmt.Sprintf("%g", v)
	}
	return ""
}

// prepares slice for Where bindings, IN/NOT IN etc
func prepareSlice(in []any) (out []string) {
	for _, value := range in {
		switch v := value.(type) {
		case string:
			out = append(out, v)
		case int:
			out = append(out, strconv.FormatInt(int64(v), 10))
		case float64:
			out = append(out, fmt.Sprintf("%g", v))
		case int64:
			out = append(out, strconv.FormatInt(v, 10))
		case uint64:
			out = append(out, strconv.FormatUint(v, 10))
		}
	}

	return
}

// composes WHERE clause string for particular query stmt
func composeWhere(whereBindings []map[string]any, startedAt int) string {
	where := " WHERE 1=1 " // where any level tables, combine with any condition
	i := startedAt
	for _, m := range whereBindings {
		for k, v := range m {
			switch vi := v.(type) {
			case []any:
				placeholders := make([]string, 0, len(vi))
				for range vi {
					placeholders = append(placeholders, "$"+strconv.Itoa(i))
					i++
				}
				where += k + " (" + strings.Join(placeholders, ", ") + ")"
			default:
				if strings.Contains(k, SqlOperatorIs) || strings.Contains(k, SqlOperatorBetween) {
					where += k + " " + vi.(string)
					break
				}
				where += k + " $" + strconv.Itoa(i)
				i++
			}
		}
	}
	return where
}

// composers ORDER BY clause string for particular query stmt
func composeOrderBy(orderBy []map[string]string, orderByRaw *string) string {
	if len(orderBy) > 0 {
		orderVal := ""
		for _, m := range orderBy {
			for field, direct := range m {
				if IsStringEmpty(orderVal) {
					orderVal = " ORDER BY " + field + " " + direct
				} else {
					orderVal += ", " + field + " " + direct
				}
			}
		}
		return orderVal
	} else if orderByRaw != nil {
		// return " ORDER BY " + *orderByRaw
		return fmt.Sprintf("%s%s", " ORDER BY ", *orderByRaw)
	}
	return ""
}

// build any date/time type with defaults preset
func buildDateTime(column, colType, defType string, isDefault bool) *qbColumn {
	col := &qbColumn{Name: column, ColumnType: qbColType(colType)}
	if isDefault {
		col.Default = &defType
	}
	return col
}

// builds column definition
func composeColumn(column *qbColumn) string {
	return column.Name + " " + string(column.ColumnType) + buildColumnOptions(column)
}

// builds column definition
func composeAddColumn(tableName string, column *qbColumn) string {
	return columnDef(tableName, column, Add)
}

// builds column definition
func composeModifyColumn(tableName string, column *qbColumn) string {
	return columnDef(tableName, column, column.Operator)
}

// builds column definition
func composeDrop(tableName string, column *qbColumn) string {
	if column.IsIndex {
		return dropIndexDef(column)
	}
	return columnDef(tableName, column, Drop)
}

// concat all definition in 1 string expression
func columnDef(tableName string, column *qbColumn, operator string) (colDef string) {
	colDef = AlterTable + tableName + operator + "COLUMN " + applyExistence(column.IfExists) + column.Name

	if operator == Rename {
		return colDef + " TO " + *column.RenameTo
	}
	if operator == Modify {
		colDef += " TYPE "
	}
	if operator != Drop {
		colDef += " " + string(column.ColumnType) + buildColumnOptions(column)
	}
	return
}

func applyExistence(ifExists uint) string {
	if ifExists == IfExistsUndeclared {
		return ""
	}
	if ifExists == IfExists {
		return IfExistsExp
	}
	return IfNotExistsExp
}

func dropIndexDef(column *qbColumn) string {
	return fmt.Sprintf("%s%s%s", "DROP INDEX ", applyExistence(column.IfExists), column.IdxName)
}

func buildColumnOptions(column *qbColumn) (colSchema string) {
	if column.IsPrimaryKey {
		colSchema += " PRIMARY KEY"
	}
	if column.IsNotNull {
		colSchema += " NOT NULL"
	}
	if column.Default != nil {
		colSchema += " DEFAULT " + *column.Default
	}
	if column.Collation != nil {
		colSchema += " COLLATE \"" + *column.Collation + "\""
	}
	return
}

// build index for table on particular column depending on an index type
func composeIndex(tableName string, column *qbColumn) string {
	if column.IsIndex {
		return "CREATE INDEX " + applyIdxConcurrency(column.IsIdxConcurrent) + applyExistence(column.IfExists) +
			column.IdxName + " ON " + tableName + " (" + column.Name + ")" + applyIncludes(column.Includes)
	}
	if column.IsUnique {
		return "CREATE UNIQUE INDEX " + applyIdxConcurrency(column.IsIdxConcurrent) + applyExistence(column.IfExists) +
			column.IdxName + " ON " + tableName + " (" + column.Name + ")" + applyIncludes(column.Includes)
	}
	if column.ForeignKey != nil {
		if column.IsIdxConcurrent {
			concurrentFk := ""
			words := strings.Fields(*column.ForeignKey)
			for _, word := range words {
				seq := " " + word + " "
				if word == Constraint {
					seq += " " + Concurrently + " "
				}
				concurrentFk += seq
			}
			return concurrentFk
		}
		return *column.ForeignKey
	}
	return ""
}

func applyIdxConcurrency(isIdxConcurrent bool) string {
	if isIdxConcurrent {
		return Concurrently
	}
	return ""
}

func applyIncludes(includes []string) string {
	if len(includes) > 0 {
		incFields := ""
		l := len(includes)
		for i, include := range includes {
			incFields += include
			if i < l-1 {
				incFields += ", "
			}
		}
		return fmt.Sprintf(" INCLUDE(%s)", incFields)
	}
	return ""
}

func composeComment(tableName string, column *qbColumn) string {
	if column.Comment != nil {
		return "COMMENT ON COLUMN " + tableName + "." + column.Name + " IS '" + *column.Comment + "'"
	}
	return ""
}

// SetCacheExecuteStmt to storage query SQL native executed
func setCacheExecuteStmt(query string) {
	cacheExecuteStmt = query
}

// GetCacheExecuteStmt return the query SQL have just executed
func getCacheExecuteStmt() string {
	return cacheExecuteStmt
}
