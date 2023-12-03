package database

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5"
	"time"
	"wbLab0/internal/configuration"
	"wbLab0/internal/models"
)

type Client interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	Close(context.Context) error
}

func NewClient(ctx context.Context, maxAttempts int, sc configuration.StConfig) (conn *pgx.Conn, err error) {
	connectionUrl := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", sc.Username, sc.Password, sc.Host, sc.Port, sc.Database)
	err = attemptDatabaseConnection(func() error {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		conn, err = pgx.Connect(ctx, connectionUrl)
		if err != nil {
			return err
		}

		return nil
	}, maxAttempts, 5*time.Second)

	if err != nil {
		fmt.Printf("Error connecting to database")
	}

	return conn, nil
}

func attemptDatabaseConnection(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--

			continue
		}
		return nil
	}
	return
}

func AddMessageToDatabase(db Client, testStruct models.IntTest) error {
	ctx := context.Background()
	tx, err := db.Begin(ctx)
	defer db.Close(ctx)

	insertItemQuery := "INSERT INTO receivedcodes (code) VALUES ($1);"
	_, err = tx.Exec(ctx, insertItemQuery, testStruct.Code)
	if err != nil {
		tx.Rollback(ctx)
		return errors.New(fmt.Sprintf("Insertion failed (%v)\n", err))
	}

	err = tx.Commit(ctx)
	if err != nil {
		return errors.New(fmt.Sprintf("Commit failed: (%v)\n", err))
	} else {
		fmt.Printf("Insertion succeded!")
		return nil
	}
}
