package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// connect to database using a single connection
func main() {
	type sensorData struct {
		Timestamp time.Time
		DeviceId  string
		Value     int
	}

	r := rand.New(rand.NewSource(42)) // Custom seed

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	/***********************************************/
	/* Single Connection to TimescaleDB/ PostgreSQL */
	/***********************************************/
	ctx := context.Background()
	connStr := os.Getenv("DATABASE_CONNECTION_STRING")
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

	/********************************************/
	/* INSERT into  relational table            */
	/********************************************/
	var data sensorData

	data.DeviceId = os.Getenv("COLLECTOR_ID")
	data.Timestamp = time.Now()
	data.Value = r.Intn(100) // Random int between 0 and 99

	queryInsertMetadata := `INSERT INTO sensor_data (device_id, timestamp, value) VALUES ($1, $2, $3);`

	_, erra := conn.Exec(ctx, queryInsertMetadata, data.DeviceId, data.Timestamp, data.Value)
	if erra != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", erra)
		os.Exit(1)
	}
	fmt.Printf("Inserted sensor data (%s, %s, %d) into database \n", data.DeviceId, data.Timestamp, data.Value)
}
