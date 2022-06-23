package api

import (
	"database/sql"
	"errors"
	"fmt"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) createTransfer(ctx *gin.Context) {
	var req transferRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}

	userPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	fromCustomer, isValid := server.validAccount(ctx, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}
	_, isValid = server.validAccount(ctx, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}

	//case have account but don't have permission for that account
	if fromCustomer.CustomerName != userPayload.Username {
		err := errors.New("unable to Transfer from not user account ")
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}

	arg := db.TransferTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, err := server.store.MakeTransferTx(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	ctx.JSON(http.StatusOK, result)

}

func (server *Server) validAccount(ctx *gin.Context, customerID int64, currency string) (db.CustomerAccount, bool) {
	customer, err := server.store.GetCustomer(ctx, customerID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errerrorReturn(err))
			return customer, false
		}

		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return customer, false
	}

	if customer.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", customer.ID, customer.Currency, currency)
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return customer, false
	}
	return customer, true
}
