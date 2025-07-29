# ☄️ Comet Documentation

<div align="center">

![☄️ Comet Documentation](https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=800&h=300&fit=crop&crop=center&auto=format&q=80)

*Complete guide to using ☄️ Comet ORM*

</div>

## Installation

<div align="center">

![Installation](https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Schema Definition (.cmt)

<div align="center">

![Schema Design](https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=600&h=250&fit=crop&crop=center&auto=format&q=80)

</div>

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

## CLI Commands

<div align="center">

![CLI Commands](https://images.unsplash.com/photo-1629654297299-c8506221ca97?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Development Workflow

<div align="center">

![Development Workflow](https://images.unsplash.com/photo-1551288049-bebda4e38f71?w=600&h=250&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Database Configuration

<div align="center">

![Database Configuration](https://images.unsplash.com/photo-1558618666-fcd25c85cd64?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Example Usage

<div align="center">

![Code Examples](https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=600&h=250&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Running Examples

<div align="center">

![Testing](https://images.unsplash.com/photo-1551650975-87deedd944c3?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

Navigate to the `testz/` directory:

```bash
cd testz/
go run blog.go
```

## Common Errors & Troubleshooting

<div align="center">

![Troubleshooting](https://images.unsplash.com/photo-1504639725590-34d0984388bd?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Advanced Configuration

<div align="center">

![Advanced Configuration](https://images.unsplash.com/photo-1518186285589-2f7649de83e0?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Best Practices

<div align="center">

![Best Practices](https://images.unsplash.com/photo-1556075798-4825dfaaf498?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

1. **Schema Design**: Keep models focused and avoid deep nesting
2. **Migrations**: Never edit existing migrations, create new ones
3. **Performance**: Use indexes on frequently queried fields
4. **Testing**: Use SQLite for fast tests, production database for integration tests
5. **Error Handling**: Always check errors from database operations

## Migration to ☄️ Comet

<div align="center">

![Migration](https://images.unsplash.com/photo-1504639725590-34d0984388bd?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

## Integration Examples

<div align="center">

![Integration](https://images.unsplash.com/photo-1460925895917-afdab827c52f?w=600&h=200&fit=crop&crop=center&auto=format&q=80)

</div>

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

This documentation covers the essential aspects of using ☄️ Comet. For more advanced features and updates, check the GitHub repository.

---

<div align="center">

![☄️ Comet Footer](https://images.unsplash.com/photo-1419242902214-272b3f66ee7a?w=800&h=200&fit=crop&crop=center&auto=format&q=80)

**Built with ❤️ by [Nitrix](https://github.com/nitrix4ly)**

</div>
