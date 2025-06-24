package config

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

func ConnectDB(ctx context.Context) (*pgx.Conn, error) {
	log.Debug("Database connection")

	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	if host == "" || port == "" || user == "" || dbname == "" {
		return nil, fmt.Errorf("Database connection parameters not specified")
	}

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	conn, err := pgx.Connect(ctx, dbURL)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(ctx); err != nil {
		return nil, err
	}

	return conn, nil
}
