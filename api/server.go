package api

import (
	"fmt"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/token"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     tool.ConfigObject
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config tool.ConfigObject, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenConfigKey)
	if err != nil {
		return nil, fmt.Errorf("unable to run maker token %w", err)
	}

	router := gin.Default()
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
		router:     router,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	router.POST("/user", server.makeNewUser)
	router.POST("/login", server.loginProcess)

	authRouter := router.Group("/").Use(authMiddleware(server.tokenMaker))

	authRouter.GET("/customer", server.listCustomer)
	authRouter.GET("/customer/id/:id", server.listCustomerByID)
	authRouter.POST("/customer", server.makeNewCustomerfunc)
	authRouter.POST("/transfers", server.createTransfer)

	server.router = router
	return server, nil
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
