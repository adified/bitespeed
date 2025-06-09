package main

import (
	"context"
	"fmt"
	"os"

	"github.com/adified/bitespeed/config"
	"github.com/jackc/pgx/v5"
)

func main() {
	config, _ := config.Loadconfig()
	// fmt.Println(config.DB_Url)

	conn, err := pgx.Connect(context.Background(), config.DB_Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())
}
