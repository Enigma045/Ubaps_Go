package main

import (
	"context"
	"fmt"
	"log"
	"ubaps/Db"
)

func main() {
	Db.ConnectDB()
	defer Db.DB.Close()

	ctx := context.Background()

	// Check if column exists first
	var exists bool
	queryCheck := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name='applications' AND column_name='comments'
		);
	`
	err := Db.DB.QueryRow(ctx, queryCheck).Scan(&exists)
	if err != nil {
		log.Fatal("Failed to check column existence:", err)
	}

	if !exists {
		fmt.Println("Adding 'comments' column to 'applications' table...")
		queryAdd := `ALTER TABLE applications ADD COLUMN comments JSONB DEFAULT '[]'::jsonb;`
		_, err = Db.DB.Exec(ctx, queryAdd)
		if err != nil {
			log.Fatal("Failed to add column:", err)
		}
		fmt.Println("✅ Column 'comments' added successfully.")
	} else {
		fmt.Println("Column 'comments' already exists.")
	}
}
