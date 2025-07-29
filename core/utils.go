package core

import (
	"reflect"
	"strings"
	"unicode"
)

func ToSnakeCase(str string) string {
	var result strings.Builder
	
	for i, r := range str {
		if i > 0 && unicode.IsUpper(r) {
			result.WriteRune('_')
		}
		result.WriteRune(unicode.ToLower(r))
	}
	
	return result.String()
}

func ToPascalCase(str string) string {
	parts := strings.Split(str, "_")
	var result strings.Builder
	
	for _, part := range parts {
		if len(part) > 0 {
			result.WriteRune(unicode.ToUpper(rune(part[0])))
			result.WriteString(part[1:])
		}
	}
	
	return result.String()
}

func ToCamelCase(str string) string {
	pascal := ToPascalCase(str)
	if len(pascal) > 0 {
		return strings.ToLower(string(pascal[0])) + pascal[1:]
	}
	return pascal
}

func ToPlural(str string) string {
	if strings.HasSuffix(str, "y") {
		return str[:len(str)-1] + "ies"
	}
	if strings.HasSuffix(str, "s") || strings.HasSuffix(str, "x") || 
	   strings.HasSuffix(str, "z") || strings.HasSuffix(str, "ch") || 
	   strings.HasSuffix(str, "sh") {
		return str + "es"
	}
	return str + "s"
}

func GetTableName(modelName string) string {
	snake := ToSnakeCase(modelName)
	return ToPlural(snake)
}

func IsZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.String:
		return rv.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return rv.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return rv.Float() == 0
	case reflect.Bool:
		return !rv.Bool()
	case reflect.Slice, reflect.Map, reflect.Array:
		return rv.Len() == 0
	case reflect.Ptr, reflect.Interface:
		return rv.IsNil()
	}
	
	return reflect.DeepEqual(v, reflect.Zero(rv.Type()).Interface())
}

func GetSQLType(goType string, driver string) string {
	baseType := strings.TrimSuffix(goType, "?")
	
	switch driver {
	case "postgres":
		return getPostgresType(baseType)
	case "mysql":
		return getMySQLType(baseType)
	case "sqlite":
		return getSQLiteType(baseType)
	default:
		return getSQLiteType(baseType)
	}
}

func getPostgresType(goType string) string {
	switch goType {
	case "int", "Int":
		return "INTEGER"
	case "int64":
		return "BIGINT"
	case "string", "String":
		return "VARCHAR(255)"
	case "bool", "Boolean":
		return "BOOLEAN"
	case "float64", "Float":
		return "DOUBLE PRECISION"
	case "time.Time", "DateTime":
		return "TIMESTAMP"
	default:
		return "TEXT"
	}
}

func getMySQLType(goType string) string {
	switch goType {
	case "int", "Int":
		return "INT"
	case "int64":
		return "BIGINT"
	case "string", "String":
		return "VARCHAR(255)"
	case "bool", "Boolean":
		return "BOOLEAN"
	case "float64", "Float":
		return "DOUBLE"
	case "time.Time", "DateTime":
		return "TIMESTAMP"
	default:
		return "TEXT"
	}
}

func getSQLiteType(goType string) string {
	switch goType {
	case "int", "Int", "int64":
		return "INTEGER"
	case "string", "String":
		return "TEXT"
	case "bool", "Boolean":
		return "INTEGER"
	case "float64", "Float":
		return "REAL"
	case "time.Time", "DateTime":
		return "DATETIME"
	default:
		return "TEXT"
	}
}

func EscapeIdentifier(identifier string) string {
	return "`" + strings.ReplaceAll(identifier, "`", "``") + "`"
}

func BuildPlaceholders(count int) string {
	if count <= 0 {
		return ""
	}
	
	placeholders := make([]string, count)
	for i := range placeholders {
		placeholders[i] = "?"
	}
	
	return strings.Join(placeholders, ", ")
}
