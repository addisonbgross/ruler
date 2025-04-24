package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	t "node/types"
	u "node/util"
	"os"
	"time"
)

var connection *pgxpool.Pool

// Push a node action to postgres
func Push(event t.NodeActionEvent) error {
	c := getConnection()
	if c == nil {
		return errors.New("failed to get event db")
	}

	jsonData, err := json.Marshal(event.Data)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO events.actions (hostname, type, data, created_at)
		VALUES ($1, $2, $3, $4)
	`

	createdAt := time.Now() // Use current time for `created_at` explicitly
	_, err = connection.Exec(context.Background(), insertQuery, event.Hostname, event.Type, jsonData, createdAt)
	if err != nil {
		return err
	}

	sugar, err := u.GetLogger()
	if err != nil {
		return err
	}

	sugar.Info(fmt.Sprintf("%s action pushed successfully", event.Type.String()))

	return nil
}

// CloseEventQueue closes the Postgres pool connection
func CloseEventQueue() {
	if connection != nil {
		connection.Close()
	}
}

// getConnection initializes or retrieves a Postgres connection instance.
// If the connection is already initialized, it returns the existing connection.
func getConnection() *pgxpool.Pool {
	if connection != nil {
		if connection.Ping(context.Background()) == nil {
			return connection
		}
		connection.Close()
		connection = nil
	}

	host := os.Getenv("POSTGRES_HOST")
	if host == "" {
		host = "queue"
	}

	port := os.Getenv("POSTGRES_PORT")
	if port == "" {
		port = "5432"
	}

	user := os.Getenv("POSTGRES_USER")
	if user == "" {
		user = "postgres"
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	if password == "" {
		password = "password"
	}

	db := os.Getenv("POSTGRES_DB")
	if db == "" {
		db = "events"
	}

	var err error
	connection, err = pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db))
	if err != nil {
		log.Default().Printf("Unable to connect to database: %v\n", err)
	}
	log.Default().Printf("Connected to postgres database: %v\n", db)

	return connection
}
