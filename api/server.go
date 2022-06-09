package api

import (
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.GET("/customer", server.listCustomer)
	router.GET("/customer/id/:id", server.listCustomerByID)
	router.POST("/customer", server.makeNewCustomerfunc)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
