package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

type Sensor struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
}

type SensorData struct {
	Timestamp time.Time
	DeviceId  string
	Value     int
}

// connect to database using a single connection
func main() {
	r := rand.New(rand.NewSource(42)) // Custom seed

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	sensor := Sensor{
		ID:   os.Getenv("COLLECTOR_ID"),
		Name: os.Getenv("COLLECTOR_NAME"),
	}

	err = registerSensor(os.Getenv("FASTAPI_BACKEND_URL"), sensor)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
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
	var sensorData SensorData

	sensorData.DeviceId = sensor.ID
	sensorData.Timestamp = time.Now()
	sensorData.Value = r.Intn(100) // Random int between 0 and 99

	queryInsertMetadata := `INSERT INTO sensor_data (device_id, timestamp, value) VALUES ($1, $2, $3);`

	_, erra := conn.Exec(ctx, queryInsertMetadata, sensorData.DeviceId, sensorData.Timestamp, sensorData.Value)
	if erra != nil {
		fmt.Fprintf(os.Stderr, "Unable to insert data into database: %v\n", erra)
		os.Exit(1)
	}
	fmt.Printf("Inserted sensor data (%s, %s, %d) into database \n", sensorData.DeviceId, sensorData.Timestamp, sensorData.Value)
}

func registerSensor(baseURL string, sensor Sensor) error {
	// Encode to JSON
	payload, err := json.Marshal(sensor)
	if err != nil {
		return err
	}

	// Create POST request
	resp, err := http.Post(baseURL+"/sensors/", "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to register sensor: status %s", resp.Status)
	}

	fmt.Println("Sensor registered successfully")
	return nil
}
