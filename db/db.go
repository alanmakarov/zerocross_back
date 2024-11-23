package db

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

var Pool *pgxpool.Pool

// Инициализация пула соединений
func InitDB() error {
	var err error
	Pool, err = pgxpool.New(context.Background(), "postgres://alex@localhost:5432/go_test")
	if err != nil {
		return fmt.Errorf("unable to create pool of connections to database: %v", err)
	}
	conn, err := Pool.Acquire(context.Background())
	defer conn.Release()
	
	if err != nil {
		return fmt.Errorf("unable to connect to database: %v", err)
	}
	return nil
}

// Закрытие пула соединений
func CloseDB() {
	Pool.Close()
	log.Print("DB Closed")
}
