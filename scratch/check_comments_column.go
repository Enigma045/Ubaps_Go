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

	var exists bool
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM information_schema.columns 
			WHERE table_name = 'applications' AND column_name = 'comments'
		);
	`
	err := Db.DB.QueryRow(ctx, query).Scan(&exists)
	if err != nil {
		log.Fatal(err)
	}

	if exists {
		fmt.Println("Comments column EXISTS")
	} else {
		fmt.Println("Comments column MISSING")
	}
}
