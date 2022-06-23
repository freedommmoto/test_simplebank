package api

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/freedommmoto/test_simplebank/token"
	"github.com/lib/pq"
	"net/http"

	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/gin-gonic/gin"
)

type makeNewCustomer struct {
	CustomerName string `json:"customer_name" binding:"required"`
	Currency     string `json:"currency" binding:"required,currency" `
}

type GetCustomerInput struct {
	CustomerID int64 `uri:"id" binding:"required,min=1"`
}
type GetCustomersInput struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=1,max=20"`
}

func (server *Server) makeNewCustomerfunc(ctx *gin.Context) {
	var req makeNewCustomer
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}
	//get payload from token
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	arg := db.CreateCustomerParams{
		CustomerName: authPayload.Username, //used payload username so user can't make a new customer with other username
		Currency:     req.Currency,
		Balance:      0,
	}

	customer, err := server.store.CreateCustomer(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			//log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errerrorReturn(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func (server *Server) listCustomerByID(ctx *gin.Context) {
	//fmt.Println("listCustomerByID ja ja")
	var req GetCustomerInput
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}

	//fmt.Println(req.CustomerID)
	customer, err := server.store.GetCustomer(ctx, req.CustomerID)
	fmt.Println(customer, "customer")

	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errerrorReturn(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	if authPayload.Username != customer.CustomerName {
		errText := errors.New("unable to get customer don't belong to you user")
		ctx.JSON(http.StatusUnauthorized, errText)
		return
	}

	ctx.JSON(http.StatusOK, customer)
}

func (server *Server) listCustomer(ctx *gin.Context) {
	var req GetCustomersInput
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}
	authPayload := ctx.MustGet(authPayloadKey).(*token.Payload)
	arg := db.ListCustomerWithOwnerParams{
		CustomerName: authPayload.Username,
		Limit:        req.PageSize,
		Offset:       (req.PageID - 1) * req.PageSize,
	}

	customerList, err := server.store.ListCustomerWithOwner(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
	}

	//fmt.Printf("%v", customerList)
	ctx.JSON(http.StatusOK, customerList)
}

func (server *Server) listCustomerWithAdmin(ctx *gin.Context) {
	var req GetCustomersInput
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}
	arg := db.ListCustomerParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	customerList, err := server.store.ListCustomer(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
	}

	//fmt.Printf("%v", customerList)
	ctx.JSON(http.StatusOK, customerList)
}

func errerrorReturn(err error) gin.H {
	return gin.H{"error": err.Error()}
}
