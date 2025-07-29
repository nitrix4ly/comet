package gen

import (
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/nitrix4ly/comet/core"
)

type Generator struct {
	parser *Parser
}

func NewGenerator() *Generator {
	return &Generator{
		parser: NewParser(),
	}
}

func (g *Generator) GenerateFromFile(schemaFile, outputDir string) error {
	schema, err := g.parser.ParseFile(schemaFile)
	if err != nil {
		return err
	}

	for _, model := range schema.Models {
		if err := g.generateModel(model, outputDir); err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) GenerateHelpers(outputDir string) error {
	return g.generateBaseFiles(outputDir)
}

func (g *Generator) generateModel(model core.ModelSchema, outputDir string) error {
	filename := filepath.Join(outputDir, strings.ToLower(model.Name)+".go")
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("model").Parse(modelTemplate))
	
	data := struct {
		Model          core.ModelSchema
		PackageName    string
		GoType         func(string) string
		DatabaseType   func(string) string
		IsOptional     func(core.FieldSchema) bool
		HasTimestamps  func() bool
	}{
		Model:       model,
		PackageName: "models",
		GoType:      g.getGoType,
		DatabaseType: func(t string) string {
			return core.GetSQLType(t, "postgres")
		},
		IsOptional: func(f core.FieldSchema) bool {
			return f.Optional
		},
		HasTimestamps: func() bool {
			return true
		},
	}

	return tmpl.Execute(file, data)
}

func (g *Generator) generateBaseFiles(outputDir string) error {
	if err := g.generateDBFile(outputDir); err != nil {
		return err
	}
	
	return g.generateConfigFile(outputDir)
}

func (g *Generator) generateDBFile(outputDir string) error {
	filename := filepath.Join(outputDir, "db.go")
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("db").Parse(dbTemplate))
	
	data := struct {
		PackageName string
	}{
		PackageName: "models",
	}

	return tmpl.Execute(file, data)
}

func (g *Generator) generateConfigFile(outputDir string) error {
	filename := filepath.Join(outputDir, "config.go")
	
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	tmpl := template.Must(template.New("config").Parse(configTemplate))
	
	data := struct {
		PackageName string
	}{
		PackageName: "models",
	}

	return tmpl.Execute(file, data)
}

func (g *Generator) getGoType(fieldType string) string {
	switch fieldType {
	case "Int":
		return "int"
	case "String":
		return "string"
	case "Boolean":
		return "bool"
	case "Float":
		return "float64"
	case "DateTime":
		return "time.Time"
	default:
		return "string"
	}
}

const modelTemplate = `package {{.PackageName}}

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/nitrix4ly/comet/core"
)

type {{.Model.Name}} struct {
{{- range .Model.Fields}}
	{{.Name}} {{if .Optional}}*{{end}}{{$.GoType .Type}} ` + "`json:\"{{.Name | ToSnakeCase}}\" db:\"{{.Name | ToSnakeCase}}\"`" + `
{{- end}}
{{- if .HasTimestamps}}
	CreatedAt time.Time ` + "`json:\"created_at\" db:\"created_at\"`" + `
	UpdatedAt time.Time ` + "`json:\"updated_at\" db:\"updated_at\"`" + `
{{- end}}
	isNew bool ` + "`json:\"-\"`" + `
}

func (m *{{.Model.Name}}) TableName() string {
	return "{{.Model.TableName}}"
}

func (m *{{.Model.Name}}) IsNew() bool {
	return m.isNew{{range .Model.Fields}}{{if .Primary}} || m.{{.Name}} == 0{{end}}{{end}}
}

func (m *{{.Model.Name}}) Save(ctx context.Context) error {
	db := core.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	now := time.Now()
	if m.IsNew() {
{{- if .HasTimestamps}}
		m.CreatedAt = now
{{- end}}
		return m.insert(ctx, db)
	}
	
{{- if .HasTimestamps}}
	m.UpdatedAt = now
{{- end}}
	return m.update(ctx, db)
}

func (m *{{.Model.Name}}) Delete(ctx context.Context) error {
	db := core.GetDB()
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	query := "DELETE FROM {{.Model.TableName}} WHERE {{range .Model.Fields}}{{if .Primary}}{{.Name | ToSnakeCase}} = ?{{end}}{{end}}"
	_, err := db.Exec(ctx, query{{range .Model.Fields}}{{if .Primary}}, m.{{.Name}}{{end}}{{end}})
	return err
}

func (m *{{.Model.Name}}) insert(ctx context.Context, db *core.DB) error {
	query := ` + "`INSERT INTO {{.Model.TableName}} (" +
		`{{range $i, $field := .Model.Fields}}{{if not .Primary}}{{if $i}}, {{end}}{{.Name | ToSnakeCase}}{{end}}{{end}}` +
		`{{if .HasTimestamps}}, created_at, updated_at{{end}}) VALUES (` +
		`{{range $i, $field := .Model.Fields}}{{if not .Primary}}{{if $i}}?, {{else}}?{{end}}{{end}}{{end}}` +
		`{{if .HasTimestamps}}, ?, ?{{end}})`+"`"+`
	
	result, err := db.Exec(ctx, query{{range .Model.Fields}}{{if not .Primary}}, m.{{.Name}}{{end}}{{end}}{{if .HasTimestamps}}, m.CreatedAt, m.UpdatedAt{{end}})
	if err != nil {
		return err
	}

{{range .Model.Fields}}{{if .Primary}}{{if .AutoGen}}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	m.{{.Name}} = {{$.GoType .Type}}(id)
{{end}}{{end}}{{end}}
	m.isNew = false
	return nil
}

func (m *{{.Model.Name}}) update(ctx context.Context, db *core.DB) error {
	query := ` + "`UPDATE {{.Model.TableName}} SET " +
		`{{range $i, $field := .Model.Fields}}{{if not .Primary}}{{if $i}}, {{end}}{{.Name | ToSnakeCase}} = ?{{end}}{{end}}` +
		`{{if .HasTimestamps}}, updated_at = ?{{end}} WHERE ` +
		`{{range .Model.Fields}}{{if .Primary}}{{.Name | ToSnakeCase}} = ?{{end}}{{end}}`+"`"+`
	
	_, err := db.Exec(ctx, query{{range .Model.Fields}}{{if not .Primary}}, m.{{.Name}}{{end}}{{end}}{{if .HasTimestamps}}, m.UpdatedAt{{end}}{{range .Model.Fields}}{{if .Primary}}, m.{{.Name}}{{end}}{{end}})
	return err
}

var {{.Model.Name}}Query = &{{.Model.Name}}QueryBuilder{}

type {{.Model.Name}}QueryBuilder struct{}

func (q *{{.Model.Name}}QueryBuilder) Find() core.QueryBuilder {
	return core.NewQueryExecutor("{{.Model.TableName}}", "{{.Model.Name}}", scan{{.Model.Name}})
}

func (q *{{.Model.Name}}QueryBuilder) FindById(ctx context.Context, id {{range .Model.Fields}}{{if .Primary}}{{$.GoType .Type}}{{end}}{{end}}) (*{{.Model.Name}}, error) {
	result, err := q.Find().Where("{{range .Model.Fields}}{{if .Primary}}{{.Name | ToSnakeCase}}{{end}}{{end}}", "=", id).First(ctx)
	if err != nil {
		return nil, err
	}
	return result.(*{{.Model.Name}}), nil
}

func (q *{{.Model.Name}}QueryBuilder) Raw(query string, args ...interface{}) core.QueryBuilder {
	return core.NewQueryExecutor("{{.Model.TableName}}", "{{.Model.Name}}", scan{{.Model.Name}})
}

func scan{{.Model.Name}}(rows *sql.Rows) (interface{}, error) {
	var m {{.Model.Name}}
	err := rows.Scan(
{{- range .Model.Fields}}
		&m.{{.Name}},
{{- end}}
{{- if .HasTimestamps}}
		&m.CreatedAt,
		&m.UpdatedAt,
{{- end}}
	)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
`

const dbTemplate = `package {{.PackageName}}

import (
	"github.com/nitrix4ly/comet/core"
	"github.com/nitrix4ly/comet/drivers"
)

func InitDB(driverName, dsn string) error {
	var driver core.Driver
	
	switch driverName {
	case "postgres":
		driver = &drivers.PostgresDriver{}
	case "mysql":
		driver = &drivers.MySQLDriver{}
	case "sqlite":
		driver = &drivers.SQLiteDriver{}
	default:
		driver = &drivers.SQLiteDriver{}
	}
	
	db, err := core.NewDB(driver, dsn)
	if err != nil {
		return err
	}
	
	core.SetDB(db)
	return nil
}
`

const configTemplate = `package {{.PackageName}}

import (
	"os"
)

type Config struct {
	DatabaseURL      string
	DatabaseProvider string
}

func LoadConfig() *Config {
	return &Config{
		DatabaseURL:      getEnv("COMET_DATABASE_URL", "sqlite://./comet.db"),
		DatabaseProvider: getEnv("COMET_DATABASE_PROVIDER", "sqlite"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
`
