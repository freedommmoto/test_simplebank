package api

import (
	"net/http"

	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type makeNewCustomer struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Currency     string `json:"currency" binding:"required,oneof=USD EUR" `
}

func (server *Server) makeNewCustomerfunc(ctx *gin.Context) {
	var req makeNewCustomer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}

	arg := db.CreateCustomerParams{
		CustomerName: req.CustomerName,
		Currency:     req.Currency,
		Balance:      0,
	}

	customer, err := server.store.CreateCustomer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func errerrorReturn(err error) gin.H {
	return gin.H{"error": err.Error()}
}
