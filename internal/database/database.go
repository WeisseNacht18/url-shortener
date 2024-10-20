package database

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var (
	Database *sql.DB
)

func NewConnection(dsn string) (err error) {
	Database, err = sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalln(err)
	}
	return
}

func CloseConnection() {
	Database.Close()
}

func CheckConnection() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := Database.PingContext(ctx)

	return err
}
