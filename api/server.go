package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/mgo/banksystem/db/sqlc"
)

type Server struct {
	store  *db.Queries
	router *gin.Engine
}

func NewServer(store *db.Queries) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/customer", server.makeNewCustomerfunc)
	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
