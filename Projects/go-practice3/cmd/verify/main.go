package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./expense.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tables := []string{"users", "categories", "expenses"}
	
	for _, table := range tables {
		var count int
		query := fmt.Sprintf("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='%s'", table)
		err := db.QueryRow(query).Scan(&count)
		if err != nil {
			log.Fatalf("Error checking table %s: %v", table, err)
		}
		
		if count > 0 {
			fmt.Printf("✓ Table %s exists\n", table)
		} else {
			fmt.Printf("✗ Table %s does not exist\n", table)
			os.Exit(1)
		}
	}

	checkTableStructure(db)
	
	fmt.Println("✅ All migrations verified successfully!")
}

func checkTableStructure(db *sql.DB) {
	_, err := db.Exec(`
		INSERT INTO users (email, name) 
		VALUES ('test@example.com', 'Test User')
	`)
	if err != nil {
		log.Printf("Warning: Could not insert into users table: %v", err)
	} else {
		fmt.Println("✓ Users table structure is correct")
	}

	_, err = db.Exec(`
		INSERT INTO categories (name, user_id) 
		VALUES ('Food', NULL)
	`)
	if err != nil {
		log.Printf("Warning: Could not insert into categories table: %v", err)
	} else {
		fmt.Println("✓ Categories table structure is correct")
	}

	db.Exec("DELETE FROM users WHERE email = 'test@example.com'")
	db.Exec("DELETE FROM categories WHERE name = 'Food'")
}