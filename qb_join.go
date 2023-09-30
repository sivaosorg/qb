package qb

import "fmt"

func (q *QbDB) buildJoin(joinType, table, on string) *QbDB {
	// q.Builder.join = append(q.Builder.join, " "+joinType+" JOIN "+table+" ON "+on+" ")
	q.Builder.join = append(q.Builder.join, fmt.Sprintf("%s%s%s%s%s%s%s", " ", joinType, " JOIN ", table, " ON ", on, " "))
	return q
}

// InnerJoin joins tables by getting elements if found in both
func (q *QbDB) InnerJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinInner, table, left+operator+right)
}

// LeftJoin joins tables by getting elements from left without those that null on the right
func (q *QbDB) LeftJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinLeft, table, left+operator+right)
}

// RightJoin joins tables by getting elements from right without those that null on the left
func (q *QbDB) RightJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinRight, table, left+operator+right)
}

// FullJoin joins tables by getting all elements of both sets
func (q *QbDB) FullJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinFull, table, left+operator+right)
}

// FullOuterJoin joins tables by getting an outer sets
func (q *QbDB) FullOuterJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinFullOuter, table, left+operator+right)
}

// Todo
// CrossJoin joins tables by getting intersection of sets
// main reason: MySQL/PostgreSQL versions are different here impl their difference
//func (r *DB) CrossJoin(table string, left string, operator string, right string) *DB {
//	return r.buildJoin(JoinCross, table, left+operator+right)
//}
