package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"myapp/models"
)

func main() {
	ctx := context.Background()
	
	if err := models.InitDB("sqlite", "./blog.db"); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	fmt.Println("☄️ Comet Blog Example")
	fmt.Println("=====================")
	
	if err := createSampleData(ctx); err != nil {
		log.Fatal("Failed to create sample data:", err)
	}
	
	if err := queryExamples(ctx); err != nil {
		log.Fatal("Failed to run queries:", err)
	}
	
	fmt.Println("\n✅ Example completed successfully!")
}

func createSampleData(ctx context.Context) error {
	fmt.Println("\n📝 Creating sample data...")
	
	user := &models.User{
		Email:    "john@comet.dev",
		Name:     "John Doe",
		Age:      28,
		IsActive: true,
		Bio:      "Software developer passionate about Go and databases",
	}
	
	if err := user.Save(ctx); err != nil {
		return fmt.Errorf("failed to create user: %v", err)
	}
	
	fmt.Printf("Created user: %s (ID: %d)\n", user.Name, user.ID)
	
	category := &models.Category{
		Name: "Technology",
		Slug: "technology",
	}
	
	if err := category.Save(ctx); err != nil {
		return fmt.Errorf("failed to create category: %v", err)
	}
	
	fmt.Printf("Created category: %s (ID: %d)\n", category.Name, category.ID)
	
	posts := []*models.Post{
		{
			Title:      "Getting Started with Comet ORM",
			Content:    "Comet is a blazing-fast, schema-first ORM for Go that makes database operations simple and type-safe.",
			Published:  true,
			AuthorID:   user.ID,
			CategoryID: &category.ID,
		},
		{
			Title:      "Building Modern Go Applications",
			Content:    "Learn how to build scalable Go applications with clean architecture and modern tooling.",
			Published:  false,
			AuthorID:   user.ID,
			CategoryID: &category.ID,
		},
		{
			Title:      "Database Design Best Practices",
			Content:    "Explore the fundamentals of good database design and how to implement them effectively.",
			Published:  true,
			AuthorID:   user.ID,
			CategoryID: &category.ID,
		},
	}
	
	for _, post := range posts {
		if err := post.Save(ctx); err != nil {
			return fmt.Errorf("failed to create post: %v", err)
		}
		fmt.Printf("Created post: %s (ID: %d)\n", post.Title, post.ID)
	}
	
	tags := []*models.Tag{
		{Name: "Go", Color: "#00ADD8"},
		{Name: "Database", Color: "#336791"},
		{Name: "ORM", Color: "#FF6B6B"},
		{Name: "Tutorial", Color: "#4ECDC4"},
	}
	
	for _, tag := range tags {
		if err := tag.Save(ctx); err != nil {
			return fmt.Errorf("failed to create tag: %v", err)
		}
		fmt.Printf("Created tag: %s (ID: %d)\n", tag.Name, tag.ID)
	}
	
	profile := &models.Profile{
		UserID:  user.ID,
		Avatar:  "https://avatar.example.com/john.jpg",
		Website: "https://johndoe.dev",
		Github:  "johndoe",
		Twitter: "@johndoe",
	}
	
	if err := profile.Save(ctx); err != nil {
		return fmt.Errorf("failed to create profile: %v", err)
	}
	
	fmt.Printf("Created profile for user: %s (ID: %d)\n", user.Name, profile.ID)
	
	return nil
}

func queryExamples(ctx context.Context) error {
	fmt.Println("\n🔍 Running query examples...")
	
	fmt.Println("\n1. Find all users:")
	users, err := models.UserQuery.Find().All(ctx)
	if err != nil {
		return err
	}
	
	for _, userInterface := range users {
		user := userInterface.(*models.User)
		fmt.Printf("  - %s (%s) - Age: %d\n", user.Name, user.Email, user.Age)
	}
	
	fmt.Println("\n2. Find user by ID:")
	user, err := models.UserQuery.FindById(ctx, 1)
	if err != nil {
		return err
	}
	fmt.Printf("  Found: %s (%s)\n", user.Name, user.Email)
	
	fmt.Println("\n3. Find published posts:")
	publishedPosts, err := models.PostQuery.Find().
		Where("published", "=", true).
		OrderBy("created_at", "DESC").
		All(ctx)
	if err != nil {
		return err
	}
	
	fmt.Printf("  Found %d published posts:\n", len(publishedPosts))
	for _, postInterface := range publishedPosts {
		post := postInterface.(*models.Post)
		fmt.Printf("    - %s\n", post.Title)
	}
	
	fmt.Println("\n4. Count total posts:")
	totalPosts, err := models.PostQuery.Find().Count(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  Total posts: %d\n", totalPosts)
	
	fmt.Println("\n5. Find posts with pagination:")
	pagedPosts, err := models.PostQuery.Find().
		OrderBy("created_at", "DESC").
		Limit(2).
		Offset(0).
		All(ctx)
	if err != nil {
		return err
	}
	
	fmt.Printf("  Page 1 (2 posts):\n")
	for _, postInterface := range pagedPosts {
		post := postInterface.(*models.Post)
		fmt.Printf("    - %s\n", post.Title)
	}
	
	fmt.Println("\n6. Find posts by category:")
	techPosts, err := models.PostQuery.Find().
		Where("category_id", "=", 1).
		All(ctx)
	if err != nil {
		return err
	}
	
	fmt.Printf("  Technology posts: %d\n", len(techPosts))
	
	fmt.Println("\n7. Check if user exists:")
	exists, err := models.UserQuery.Find().
		Where("email", "=", "john@comet.dev").
		Exists(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  User exists: %t\n", exists)
	
	fmt.Println("\n8. Find first post:")
	firstPost, err := models.PostQuery.Find().
		OrderBy("created_at", "ASC").
		First(ctx)
	if err != nil {
		return err
	}
	post := firstPost.(*models.Post)
	fmt.Printf("  First post: %s\n", post.Title)
	
	fmt.Println("\n9. Update user:")
	user.Bio = "Updated bio: Senior Go developer and Comet contributor"
	if err := user.Save(ctx); err != nil {
		return err
	}
	fmt.Printf("  Updated user bio\n")
	
	fmt.Println("\n10. Find active users:")
	activeUsers, err := models.UserQuery.Find().
		Where("is_active", "=", true).
		Where("age", ">=", 18).
		All(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  Active adult users: %d\n", len(activeUsers))
	
	return nil
}

func demonstrateAdvancedFeatures(ctx context.Context) error {
	fmt.Println("\n✨ Advanced features demo...")
	
	fmt.Println("\n1. Complex WHERE conditions:")
	posts, err := models.PostQuery.Find().
		Where("published", "=", true).
		Where("created_at", ">", time.Now().AddDate(0, -1, 0)).
		OrderBy("created_at", "DESC").
		Limit(5).
		All(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  Recent published posts: %d\n", len(posts))
	
	fmt.Println("\n2. Select specific fields:")
	titleOnlyPosts, err := models.PostQuery.Find().
		Select("title", "created_at").
		Where("published", "=", true).
		All(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  Posts with selected fields: %d\n", len(titleOnlyPosts))
	
	fmt.Println("\n3. Raw SQL query:")
	rawPosts, err := models.PostQuery.Raw(
		"SELECT * FROM posts WHERE title LIKE ? ORDER BY created_at DESC",
		"%Comet%",
	).All(ctx)
	if err != nil {
		return err
	}
	fmt.Printf("  Posts containing 'Comet': %d\n", len(rawPosts))
	
	return nil
}

func cleanupData(ctx context.Context) error {
	fmt.Println("\n🧹 Cleaning up sample data...")
	
	posts, err := models.PostQuery.Find().All(ctx)
	if err != nil {
		return err
	}
	
	for _, postInterface := range posts {
		post := postInterface.(*models.Post)
		if err := post.Delete(ctx); err != nil {
			return err
		}
	}
	
	fmt.Println("  Deleted all posts")
	return nil
}
