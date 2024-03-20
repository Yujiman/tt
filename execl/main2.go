package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/jackc/pgx/v4" // Import the pgx driver
)

func main() {

	connString := "user=postgres password=0SQwGDbjfVyve9HP dbname=halyk_stage host=read.db.halyk-travel.com port=5432 sslmode=disable"

	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Println("Error connecting to the database:", err)
		return
	}
	defer conn.Close(context.Background())

	// Read SQL queries from the file
	sqlBytes, err := ioutil.ReadFile("result.sql")
	if err != nil {
		fmt.Println("Error reading result.sql file:", err)
		return
	}

	// Split SQL queries by semicolon
	sqlStatements := strings.Split(string(sqlBytes), ";")

	// Execute each SQL query
	for _, stmt := range sqlStatements {
		if strings.TrimSpace(stmt) == "" {
			continue
		}

		_, err := conn.Exec(context.Background(), stmt)
		if err != nil {
			fmt.Println("Error executing SQL query:", err)
			continue
		}
		fmt.Println("Executed SQL query:", stmt)
	}
}
