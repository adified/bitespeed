package main

import (
	"context"
	"fmt"
	"os"

	"github.com/adified/bitespeed/api"
	"github.com/adified/bitespeed/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	config, _ := config.Loadconfig()
	// fmt.Println(config.DB_Url)

	dbpool, err := pgxpool.New(context.Background(), config.DB_Url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer dbpool.Close()

	server := api.NewServer(dbpool)
	server.SetupRouter()

	server.Start(config.Address)
}
