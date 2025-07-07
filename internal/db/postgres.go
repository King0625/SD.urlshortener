package db

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

func InitPostgres(dsn string) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s", os.Getenv("POSTGRES_DSN"))
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close(context.Background())
	return conn, nil
}
