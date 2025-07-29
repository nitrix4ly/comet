# 📚 Comet Documentation

## 💻 Installation

### From Go Install (Recommended)
```bash
go install github.com/nitrix4ly/comet/cli@latest
```

### From Source
```bash
git clone https://github.com/nitrix4ly/comet.git
cd comet
go build -o comet ./cli
sudo mv comet /usr/local/bin/
```

## 🧩 Schema Definition (.cmt)

Create schema files in the `schema/` directory with `.cmt` extension:

```prisma
model User {
  id        Int      @id @auto
  email     String   @unique
  name      String?
  age       Int      @default(0)
  createdAt DateTime @default(now())
  updatedAt DateTime @updatedAt
  
  posts     Post[]   @relation("UserPosts")
}

model Post {
  id       Int    @id @auto
  title    String
  content  String?
  authorId Int
  
  author   User   @relation("UserPosts", fields: [authorId], references: [id])
}
```

### Field Types
- `Int` - Integer
- `String` - Text
- `Boolean` - True/false
- `DateTime` - Timestamp
- `Float` - Decimal number

### Attributes
- `@id` - Primary key
- `@auto` - Auto-increment
- `@unique` - Unique constraint
- `@default(value)` - Default value
- `@updatedAt` - Auto-update timestamp
- `@relation(name)` - Define relationships

### Modifiers
- `?` - Optional field (nullable)
- `[]` - Array/slice

## 🔧 CLI Commands

### Generate Models
```bash
comet gen
```
Generates Go structs and query methods from schema files.

### Run Migrations
```bash
comet migrate
```
Creates and applies database migrations.

### Seed Database
```bash
comet seed
```
Runs seed files from `seeds/` directory.

### Additional Options
```bash
comet gen --output models/     # Custom output directory
comet migrate --dry-run        # Preview migrations
comet seed --file seeds/users.go
```

## 🔁 Development Workflow

1. **Define Schema**: Create or modify `.cmt` files in `schema/`
2. **Generate Code**: Run `comet gen` to create Go models
3. **Create Migration**: Run `comet migrate` to update database
4. **Use Models**: Import generated models in your application

### After Schema Changes
```bash
# Always regenerate after schema changes
comet gen
comet migrate
```

## 🔌 Database Configuration

### Connection Strings

Create a `comet.yaml` config file:

```yaml
database:
  provider: "postgres"  # postgres, mysql, sqlite
  url: "postgres://user:pass@localhost/dbname?sslmode=disable"
  
  # Alternative format
  host: "localhost"
  port: 5432
  user: "myuser"
  password: "mypass"
  database: "mydb"
```

### Environment Variables
```bash
export COMET_DATABASE_URL="postgres://user:pass@localhost/dbname"
export COMET_DATABASE_PROVIDER="postgres"
```

### DSN Examples

**PostgreSQL:**
```
postgres://user:password@localhost:5432/database?sslmode=disable
```

**MySQL:**
```
user:password@tcp(localhost:3306)/database?parseTime=true
```

**SQLite:**
```
file:./database.db?cache=shared&mode=rwc
```

## 🧪 Example Usage

### Basic CRUD Operations

```go
package main

import (
    "context"
    "myapp/models"
)

func main() {
    ctx := context.Background()
    
    // Create
    user := &models.User{
        Email: "john@example.com",
        Name:  "John Doe",
        Age:   30,
    }
    err := user.Save(ctx)
    
    // Find by ID
    user, err = models.User.FindById(ctx, 1)
    
    // Find with conditions
    users, err := models.User.Find().
        Where("age", ">", 18).
        Where("email", "LIKE", "%@gmail.com").
        OrderBy("createdAt", "DESC").
        Limit(10).
        All(ctx)
    
    // Update
    user.Name = "John Smith"
    err = user.Save(ctx)
    
    // Delete
    err = user.Delete(ctx)
}
```

### Advanced Queries

```go
// Count records
count, err := models.User.Find().
    Where("age", ">=", 18).
    Count(ctx)

// First/Last
user, err := models.User.Find().
    OrderBy("createdAt", "DESC").
    First(ctx)

// Exists
exists, err := models.User.Find().
    Where("email", "=", "test@example.com").
    Exists(ctx)

// Raw SQL
users, err := models.User.Raw(`
    SELECT * FROM users 
    WHERE created_at > $1
`, time.Now().AddDate(0, -1, 0)).All(ctx)
```

### Relationships

```go
// Include related data
posts, err := models.Post.Find().
    Include("author").
    All(ctx)

// Access related data
for _, post := range posts {
    fmt.Println(post.Author.Name)
}

// Create with relations
post := &models.Post{
    Title:    "My Post",
    Content:  "Post content",
    AuthorId: user.Id,
}
err = post.Save(ctx)
```

## 🧪 Running Examples

Navigate to the `testz/` directory:

```bash
cd testz/
go run blog.go
```

## 🛠 Common Errors & Troubleshooting

### "models package not found"
**Solution**: Run `comet gen` to generate models first.

### "no such table" error
**Solution**: Run `comet migrate` to create database tables.

### Connection refused
**Solution**: Check database is running and connection string is correct.

### "field not found" after schema changes
**Solution**: Regenerate models with `comet gen` and run `comet migrate`.

### Performance Issues
- Use `Limit()` for large datasets
- Add database indexes for frequently queried fields
- Use `Select()` to fetch only needed columns

## 🔧 Advanced Configuration

### Custom Naming Conventions
```yaml
naming:
  tables: "snake_case"      # user_posts
  columns: "snake_case"     # created_at
  models: "PascalCase"      # UserPost
```

### Connection Pooling
```yaml
database:
  max_open_conns: 25
  max_idle_conns: 5
  conn_max_lifetime: "1h"
```

### Logging
```yaml
logging:
  level: "info"        # debug, info, warn, error
  sql_queries: true    # Log all SQL queries
```

## 🚀 Best Practices

1. **Schema Design**: Keep models focused and avoid deep nesting
2. **Migrations**: Never edit existing migrations, create new ones
3. **Performance**: Use indexes on frequently queried fields
4. **Testing**: Use SQLite for fast tests, production database for integration tests
5. **Error Handling**: Always check errors from database operations

## 📈 Migration to Comet

### From GORM
```go
// GORM
db.Where("age > ?", 18).Find(&users)

// Comet
users, err := models.User.Find().Where("age", ">", 18).All(ctx)
```

### From Raw SQL
```go
// Raw SQL
rows, err := db.Query("SELECT * FROM users WHERE age > $1", 18)

// Comet
users, err := models.User.Find().Where("age", ">", 18).All(ctx)
```

## 🔗 Integration Examples

### With Gin
```go
func GetUsers(c *gin.Context) {
    users, err := models.User.Find().All(c.Request.Context())
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, users)
}
```

### With GraphQL
```go
func (r *queryResolver) Users(ctx context.Context) ([]*models.User, error) {
    return models.User.Find().All(ctx)
}
```

This documentation covers the essential aspects of using Comet. For more advanced features and updates, check the GitHub repository.
