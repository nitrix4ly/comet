package drivers

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/nitrix4ly/comet/core"
	_ "github.com/lib/pq"
)

type PostgresDriver struct{}

func (d *PostgresDriver) Connect(dsn string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}
	
	return db, nil
}

func (d *PostgresDriver) Migrate(schema *core.Schema) error {
	return fmt.Errorf("migrations not implemented yet")
}

func (d *PostgresDriver) BuildQuery(query *core.Query) (string, []interface{}) {
	var parts []string
	var args []interface{}
	
	fields := strings.Join(query.Fields, ", ")
	parts = append(parts, fmt.Sprintf("SELECT %s FROM %s", fields, query.Table))
	
	if len(query.Wheres) > 0 {
		var whereParts []string
		argIndex := 1
		
		for _, where := range query.Wheres {
			operator := where.Operator
			if where.Not {
				operator = "NOT " + operator
			}
			
			if where.Operator == "IN" {
				whereParts = append(whereParts, fmt.Sprintf("%s %s %v", where.Field, operator, where.Value))
			} else {
				whereParts = append(whereParts, fmt.Sprintf("%s %s $%d", where.Field, operator, argIndex))
				args = append(args, where.Value)
				argIndex++
			}
		}
		parts = append(parts, "WHERE "+strings.Join(whereParts, " AND "))
	}
	
	if len(query.Orders) > 0 {
		var orderParts []string
		for _, order := range query.Orders {
			orderParts = append(orderParts, fmt.Sprintf("%s %s", order.Field, order.Direction))
		}
		parts = append(parts, "ORDER BY "+strings.Join(orderParts, ", "))
	}
	
	if query.LimitVal != nil {
		parts = append(parts, fmt.Sprintf("LIMIT %d", *query.LimitVal))
	}
	
	if query.OffsetVal != nil {
		parts = append(parts, fmt.Sprintf("OFFSET %d", *query.OffsetVal))
	}
	
	return strings.Join(parts, " "), args
}

func (d *PostgresDriver) GetDialect() string {
	return "postgres"
}

func (d *PostgresDriver) CreateTable(model core.ModelSchema) string {
	var columns []string
	
	for _, field := range model.Fields {
		column := d.buildColumnDefinition(field)
		columns = append(columns, column)
	}
	
	sql := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (\n  %s\n)",
		model.TableName,
		strings.Join(columns, ",\n  "))
	
	return sql
}

func (d *PostgresDriver) buildColumnDefinition(field core.FieldSchema) string {
	var parts []string
	
	parts = append(parts, field.Name)
	
	sqlType := core.GetSQLType(field.Type, "postgres")
	if field.Primary && field.AutoGen {
		sqlType = "SERIAL"
	}
	parts = append(parts, sqlType)
	
	if field.Primary {
		parts = append(parts, "PRIMARY KEY")
	}
	
	if field.Unique && !field.Primary {
		parts = append(parts, "UNIQUE")
	}
	
	if !field.Optional && !field.Primary {
		parts = append(parts, "NOT NULL")
	}
	
	if field.Default != nil {
		switch v := field.Default.(type) {
		case string:
			if v == "CURRENT_TIMESTAMP" {
				parts = append(parts, "DEFAULT CURRENT_TIMESTAMP")
			} else {
				parts = append(parts, fmt.Sprintf("DEFAULT '%s'", v))
			}
		case bool:
			parts = append(parts, fmt.Sprintf("DEFAULT %t", v))
		default:
			parts = append(parts, fmt.Sprintf("DEFAULT %v", v))
		}
	}
	
	return strings.Join(parts, " ")
}
