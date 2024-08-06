package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jessicacastro/ama-application/go/internal/api"
	"github.com/jessicacastro/ama-application/go/internal/store/pgstore"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}

	context := context.Background()

	pool, err := pgxpool.New(context, fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		os.Getenv("WS_DATABASE_HOST"),
		os.Getenv("WS_DATABASE_PORT"),
		os.Getenv("WS_DATABASE_USER"),
		os.Getenv("WS_DATABASE_PASSWORD"),
		os.Getenv("WS_DATABASE_NAME"),
	))

	if err != nil {
		panic(err)
	}

	// defer make sure the pool is closed when the main function ends before the program exits, like a cleanup function
	defer pool.Close()

	if err := pool.Ping(context); err != nil {
		panic(err)
	}

	handler := api.NewAPIHandler(pgstore.New(pool))

	go func() {
		if err := http.ListenAndServe(":8080", handler); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				panic(err)
			}
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
}
