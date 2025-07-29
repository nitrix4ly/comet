package core

import (
	"context"
	"database/sql"
	"time"
)

type Model interface {
	TableName() string
	Save(ctx context.Context) error
	Delete(ctx context.Context) error
	IsNew() bool
}

type QueryBuilder interface {
	Where(field, operator string, value interface{}) QueryBuilder
	WhereIn(field string, values []interface{}) QueryBuilder
	WhereNot(field, operator string, value interface{}) QueryBuilder
	OrderBy(field, direction string) QueryBuilder
	Limit(limit int) QueryBuilder
	Offset(offset int) QueryBuilder
	Select(fields ...string) QueryBuilder
	Include(relations ...string) QueryBuilder
	
	All(ctx context.Context) ([]interface{}, error)
	First(ctx context.Context) (interface{}, error)
	Last(ctx context.Context) (interface{}, error)
	Count(ctx context.Context) (int64, error)
	Exists(ctx context.Context) (bool, error)
}

type Driver interface {
	Connect(dsn string) (*sql.DB, error)
	Migrate(schema *Schema) error
	BuildQuery(query *Query) (string, []interface{})
	GetDialect() string
}

type Schema struct {
	Models []ModelSchema `json:"models"`
}

type ModelSchema struct {
	Name      string        `json:"name"`
	TableName string        `json:"table_name"`
	Fields    []FieldSchema `json:"fields"`
	Relations []Relation    `json:"relations"`
}

type FieldSchema struct {
	Name         string      `json:"name"`
	Type         string      `json:"type"`
	Optional     bool        `json:"optional"`
	Unique       bool        `json:"unique"`
	Primary      bool        `json:"primary"`
	AutoGen      bool        `json:"auto_gen"`
	Default      interface{} `json:"default"`
	DatabaseType string      `json:"database_type"`
}

type Relation struct {
	Name      string   `json:"name"`
	Type      string   `json:"type"`
	Model     string   `json:"model"`
	Fields    []string `json:"fields"`
	References []string `json:"references"`
}

type Query struct {
	Table     string
	Fields    []string
	Wheres    []WhereClause
	Orders    []OrderClause
	LimitVal  *int
	OffsetVal *int
	Includes  []string
}

type WhereClause struct {
	Field    string
	Operator string
	Value    interface{}
	Not      bool
}

type OrderClause struct {
	Field     string
	Direction string
}

type BaseModel struct {
	ID        int       `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
	isNew     bool      `json:"-"`
}

func (b *BaseModel) IsNew() bool {
	return b.isNew || b.ID == 0
}

func (b *BaseModel) SetNew(isNew bool) {
	b.isNew = isNew
}

func (b *BaseModel) Touch() {
	now := time.Now()
	if b.IsNew() {
		b.CreatedAt = now
	}
	b.UpdatedAt = now
}

type DB struct {
	conn   *sql.DB
	driver Driver
}

func NewDB(driver Driver, dsn string) (*DB, error) {
	conn, err := driver.Connect(dsn)
	if err != nil {
		return nil, err
	}
	
	return &DB{
		conn:   conn,
		driver: driver,
	}, nil
}

func (db *DB) Query(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return db.conn.QueryContext(ctx, query, args...)
}

func (db *DB) QueryRow(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return db.conn.QueryRowContext(ctx, query, args...)
}

func (db *DB) Exec(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return db.conn.ExecContext(ctx, query, args...)
}

func (db *DB) Close() error {
	return db.conn.Close()
}

var GlobalDB *DB

func SetDB(db *DB) {
	GlobalDB = db
}

func GetDB() *DB {
	return GlobalDB
}
