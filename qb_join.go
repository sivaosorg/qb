package qb

import "fmt"

// InnerJoin joins tables by getting elements if found in both
func (q *QbDB) InnerJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinInner, table, left+operator+right)
}

// InnerJoinIf joins tables by getting elements if found in both
func (q *QbDB) InnerJoinIf(fnc func() bool, table, left, operator, right string) *QbDB {
	if !fnc() {
		return q
	}
	return q.InnerJoin(table, left, operator, right)
}

// LeftJoin joins tables by getting elements from left without those that null on the right
func (q *QbDB) LeftJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinLeft, table, left+operator+right)
}

// LeftJoinIf joins tables by getting elements from left without those that null on the right
func (q *QbDB) LeftJoinIf(fnc func() bool, table, left, operator, right string) *QbDB {
	if !fnc() {
		return q
	}
	return q.LeftJoin(table, left, operator, right)
}

// RightJoin joins tables by getting elements from right without those that null on the left
func (q *QbDB) RightJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinRight, table, left+operator+right)
}

// RightJoinIf joins tables by getting elements from right without those that null on the left
func (q *QbDB) RightJoinIf(fnc func() bool, table, left, operator, right string) *QbDB {
	if !fnc() {
		return q
	}
	return q.RightJoin(table, left, operator, right)
}

// FullJoin joins tables by getting all elements of both sets
func (q *QbDB) FullJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinFull, table, left+operator+right)
}

// FullJoinIf joins tables by getting all elements of both sets
func (q *QbDB) FullJoinIf(fnc func() bool, table, left, operator, right string) *QbDB {
	if !fnc() {
		return q
	}
	return q.FullJoin(table, left, operator, right)
}

// FullOuterJoin joins tables by getting an outer sets
func (q *QbDB) FullOuterJoin(table, left, operator, right string) *QbDB {
	return q.buildJoin(JoinFullOuter, table, left+operator+right)
}

// FullOuterJoinIf joins tables by getting an outer sets
func (q *QbDB) FullOuterJoinIf(fnc func() bool, table, left, operator, right string) *QbDB {
	if !fnc() {
		return q
	}
	return q.FullOuterJoin(table, left, operator, right)
}

func (q *QbDB) buildJoin(joinType, table, on string) *QbDB {
	q.Builder.join = append(q.Builder.join, fmt.Sprintf("%s%s%s%s%s%s%s", " ", joinType, " JOIN ", table, " ON ", on, " "))
	return q
}

// Todo
// CrossJoin joins tables by getting intersection of sets
// main reason: MySQL/PostgreSQL versions are different here impl their difference
//func (r *DB) CrossJoin(table string, left string, operator string, right string) *DB {
//	return r.buildJoin(JoinCross, table, left+operator+right)
//}
