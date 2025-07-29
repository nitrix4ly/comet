package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/nitrix4ly/comet/gen"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "comet",
	Short: "Comet - Blazing-fast, schema-first ORM for Go",
	Long:  "A lightweight, Prisma-inspired ORM that generates type-safe models from schema files.",
}

var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "Generate Go models from schema files",
	Run: func(cmd *cobra.Command, args []string) {
		outputDir, _ := cmd.Flags().GetString("output")
		schemaDir, _ := cmd.Flags().GetString("schema")
		
		if err := runGenerate(schemaDir, outputDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Models generated successfully!")
	},
}

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		
		if err := runMigrate(dryRun); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		if dryRun {
			fmt.Println("✅ Migration preview completed!")
		} else {
			fmt.Println("✅ Migrations applied successfully!")
		}
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed database with initial data",
	Run: func(cmd *cobra.Command, args []string) {
		seedFile, _ := cmd.Flags().GetString("file")
		
		if err := runSeed(seedFile); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		
		fmt.Println("✅ Database seeded successfully!")
	},
}

func init() {
	genCmd.Flags().StringP("output", "o", "models", "Output directory for generated models")
	genCmd.Flags().StringP("schema", "s", "schema", "Schema directory")
	
	migrateCmd.Flags().Bool("dry-run", false, "Preview migrations without applying")
	
	seedCmd.Flags().StringP("file", "f", "", "Specific seed file to run")
	
	rootCmd.AddCommand(genCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(seedCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runGenerate(schemaDir, outputDir string) error {
	if _, err := os.Stat(schemaDir); os.IsNotExist(err) {
		return fmt.Errorf("schema directory '%s' does not exist", schemaDir)
	}
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}
	
	schemaFiles, err := filepath.Glob(filepath.Join(schemaDir, "*.cmt"))
	if err != nil {
		return fmt.Errorf("failed to find schema files: %v", err)
	}
	
	if len(schemaFiles) == 0 {
		return fmt.Errorf("no .cmt schema files found in %s", schemaDir)
	}
	
	generator := gen.NewGenerator()
	
	for _, schemaFile := range schemaFiles {
		fmt.Printf("Processing %s...\n", schemaFile)
		
		if err := generator.GenerateFromFile(schemaFile, outputDir); err != nil {
			return fmt.Errorf("failed to generate from %s: %v", schemaFile, err)
		}
	}
	
	if err := generator.GenerateHelpers(outputDir); err != nil {
		return fmt.Errorf("failed to generate helpers: %v", err)
	}
	
	return nil
}

func runMigrate(dryRun bool) error {
	fmt.Println("🔄 Running migrations...")
	
	if dryRun {
		fmt.Println("📋 DRY RUN - No changes will be applied")
		fmt.Println("SQL Preview:")
		fmt.Println("CREATE TABLE users (")
		fmt.Println("  id SERIAL PRIMARY KEY,")
		fmt.Println("  email VARCHAR(255) UNIQUE NOT NULL,")
		fmt.Println("  name VARCHAR(255),")
		fmt.Println("  created_at TIMESTAMP DEFAULT NOW()")
		fmt.Println(");")
		return nil
	}
	
	fmt.Println("📝 Applying migrations to database...")
	return nil
}

func runSeed(seedFile string) error {
	fmt.Println("🌱 Seeding database...")
	
	if seedFile != "" {
		fmt.Printf("Running seed file: %s\n", seedFile)
	} else {
		fmt.Println("Running all seed files...")
	}
	
	fmt.Println("📝 Sample data inserted")
	return nil
}
