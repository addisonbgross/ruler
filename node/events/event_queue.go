package events

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	t "node/types"
	u "node/util"
	"os"
	"time"
)

var connection *pgxpool.Pool

// Push a node action to postgres
func Push(event t.NodeActionEvent) error {
	c, err := getConnection()
	if err != nil {
		return err
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
	_, err = c.Exec(context.Background(), insertQuery, event.Hostname, event.Type, jsonData, createdAt)
	if err != nil {
		return err
	}

	logger, err := u.GetLogger()
	if err != nil {
		return err
	}

	logger.Info(fmt.Sprintf("%s action pushed successfully", event.Type.String()))

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
func getConnection() (*pgxpool.Pool, error) {
	if connection != nil {
		if connection.Ping(context.Background()) == nil {
			return connection, nil
		}
		connection.Close()
		connection = nil
	}

	logger, err := u.GetLogger()
	if err != nil {
		return nil, err
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

	connection, err = pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, db))
	if err != nil {
		return nil, err
	}
	logger.Info(fmt.Sprintf("Connected to postgres database: %v", db))

	return connection, nil
}
