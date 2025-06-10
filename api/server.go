package api

import (
	db "github.com/adified/bitespeed/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Server serves HTTP requests for our banking service.
type Server struct {
	Querier db.Queries
	router  *gin.Engine
}

func NewServer(dbpool *pgxpool.Pool) *Server {
	Querier := db.New(dbpool)
	server := &Server{
		Querier: *Querier,
	}
	return server
}

func (server *Server) SetupRouter() {
	router := gin.Default()
	router.POST("/identify", CheckifExists(server))

	server.router = router
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
