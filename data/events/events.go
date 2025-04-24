package events

import (
	"context"
	t "data/types"
	u "data/util"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"os"
)

var connection *pgxpool.Pool

// Read all events from postgres
func Read() ([]t.NodeActionEventResponse, error) {
	c, err := getConnection()
	if err != nil {
		return nil, err
	}

	logger, err := u.GetLogger()
	if err != nil {
		return nil, err
	}

	rows, err := c.Query(context.Background(), "select hostname, type, data, created_at from events.actions")
	if err != nil {
		return nil, err
	}

	var actions []t.NodeActionEventResponse
	for rows.Next() {
		var action t.NodeActionEventResponse
		err = rows.Scan(&action.Hostname, &action.Type, &action.Data, &action.CreatedAt)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}

	logger.Info("All events read successfully")

	return actions, nil
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
