package core

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type QueryExecutor struct {
	query     *Query
	modelType string
	scanner   func(*sql.Rows) (interface{}, error)
}

func NewQueryExecutor(table, modelType string, scanner func(*sql.Rows) (interface{}, error)) *QueryExecutor {
	return &QueryExecutor{
		query: &Query{
			Table:  table,
			Fields: []string{"*"},
		},
		modelType: modelType,
		scanner:   scanner,
	}
}

func (qe *QueryExecutor) Where(field, operator string, value interface{}) QueryBuilder {
	qe.query.Wheres = append(qe.query.Wheres, WhereClause{
		Field:    field,
		Operator: operator,
		Value:    value,
		Not:      false,
	})
	return qe
}

func (qe *QueryExecutor) WhereIn(field string, values []interface{}) QueryBuilder {
	placeholders := make([]string, len(values))
	for i := range placeholders {
		placeholders[i] = "?"
	}
	
	qe.query.Wheres = append(qe.query.Wheres, WhereClause{
		Field:    field,
		Operator: "IN",
		Value:    fmt.Sprintf("(%s)", strings.Join(placeholders, ",")),
	})
	return qe
}

func (qe *QueryExecutor) WhereNot(field, operator string, value interface{}) QueryBuilder {
	qe.query.Wheres = append(qe.query.Wheres, WhereClause{
		Field:    field,
		Operator: operator,
		Value:    value,
		Not:      true,
	})
	return qe
}

func (qe *QueryExecutor) OrderBy(field, direction string) QueryBuilder {
	qe.query.Orders = append(qe.query.Orders, OrderClause{
		Field:     field,
		Direction: strings.ToUpper(direction),
	})
	return qe
}

func (qe *QueryExecutor) Limit(limit int) QueryBuilder {
	qe.query.LimitVal = &limit
	return qe
}

func (qe *QueryExecutor) Offset(offset int) QueryBuilder {
	qe.query.OffsetVal = &offset
	return qe
}

func (qe *QueryExecutor) Select(fields ...string) QueryBuilder {
	qe.query.Fields = fields
	return qe
}

func (qe *QueryExecutor) Include(relations ...string) QueryBuilder {
	qe.query.Includes = append(qe.query.Includes, relations...)
	return qe
}

func (qe *QueryExecutor) All(ctx context.Context) ([]interface{}, error) {
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	query, args := qe.buildSelectQuery()
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var results []interface{}
	for rows.Next() {
		item, err := qe.scanner(rows)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	
	return results, rows.Err()
}

func (qe *QueryExecutor) First(ctx context.Context) (interface{}, error) {
	qe.query.LimitVal = intPtr(1)
	
	db := GetDB()
	if db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	
	query, args := qe.buildSelectQuery()
	rows, err := db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	
	return qe.scanner(rows)
}

func (qe *QueryExecutor) Last(ctx context.Context) (interface{}, error) {
	if len(qe.query.Orders) == 0 {
		qe.query.Orders = append(qe.query.Orders, OrderClause{
			Field:     "id",
			Direction: "DESC",
		})
	}
	
	return qe.First(ctx)
}

func (qe *QueryExecutor) Count(ctx context.Context) (int64, error) {
	db := GetDB()
	if db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	
	countQuery := &Query{
		Table:     qe.query.Table,
		Fields:    []string{"COUNT(*)"},
		Wheres:    qe.query.Wheres,
		Orders:    nil,
		LimitVal:  nil,
		OffsetVal: nil,
	}
	
	query, args := qe.buildSelectQueryFromQuery(countQuery)
	
	var count int64
	err := db.QueryRow(ctx, query, args...).Scan(&count)
	return count, err
}

func (qe *QueryExecutor) Exists(ctx context.Context) (bool, error) {
	count, err := qe.Count(ctx)
	return count > 0, err
}

func (qe *QueryExecutor) buildSelectQuery() (string, []interface{}) {
	return qe.buildSelectQueryFromQuery(qe.query)
}

func (qe *QueryExecutor) buildSelectQueryFromQuery(q *Query) (string, []interface{}) {
	var parts []string
	var args []interface{}
	
	fields := strings.Join(q.Fields, ", ")
	parts = append(parts, fmt.Sprintf("SELECT %s FROM %s", fields, q.Table))
	
	if len(q.Wheres) > 0 {
		var whereParts []string
		for _, where := range q.Wheres {
			operator := where.Operator
			if where.Not {
				operator = "NOT " + operator
			}
			
			if where.Operator == "IN" {
				whereParts = append(whereParts, fmt.Sprintf("%s %s %v", where.Field, operator, where.Value))
			} else {
				whereParts = append(whereParts, fmt.Sprintf("%s %s ?", where.Field, operator))
				args = append(args, where.Value)
			}
		}
		parts = append(parts, "WHERE "+strings.Join(whereParts, " AND "))
	}
	
	if len(q.Orders) > 0 {
		var orderParts []string
		for _, order := range q.Orders {
			orderParts = append(orderParts, fmt.Sprintf("%s %s", order.Field, order.Direction))
		}
		parts = append(parts, "ORDER BY "+strings.Join(orderParts, ", "))
	}
	
	if q.LimitVal != nil {
		parts = append(parts, fmt.Sprintf("LIMIT %d", *q.LimitVal))
	}
	
	if q.OffsetVal != nil {
		parts = append(parts, fmt.Sprintf("OFFSET %d", *q.OffsetVal))
	}
	
	return strings.Join(parts, " "), args
}

func intPtr(i int) *int {
	return &i
}
