package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

// connect to database using a single connection
func main() {
	/***********************************************/
	/* Single Connection to TimescaleDB/ PostgreSQL */
	/***********************************************/
	ctx := context.Background()
	connStr := "postgres://postgres:password@localhost/timescale"
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	//run a simple query to check our connection
	var greeting string
	err = conn.QueryRow(ctx, "select 'Hello, Timescale!'").Scan(&greeting)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(greeting)
}
