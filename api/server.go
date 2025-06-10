package api

import (
	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	pool    *pgxpool.Pool
	Querier *db.Queries
	router  *gin.Engine
}

func NewServer(dbpool *pgxpool.Pool) *Server {
	server := &Server{
		pool:    dbpool,         // Store the pool
		Querier: db.New(dbpool), // Create a querier instance
	}
	return server
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/identify", server.CheckifExists())

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
