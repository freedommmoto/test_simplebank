package api

import (
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	// ทำถึงนี้ https://www.youtube.com/watch?v=B3xnJI2lHmc&list=PLy_6D98if3ULEtXtNSY_2qN21VCKgoQAE&index=17
	// เรียนได้ครึ่งทางแล้ว
	router.GET("/customer", server.listCustomer)
	router.GET("/customer/id/:id", server.listCustomerByID)
	router.POST("/customer", server.makeNewCustomerfunc)
	router.POST("/transfers", server.createTransfer)
	router.POST("/user", server.makeNewUser)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
