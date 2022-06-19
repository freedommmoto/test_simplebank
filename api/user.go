package api

import (
	"database/sql"
	db "github.com/freedommmoto/test_simplebank/db/sqlc"
	"github.com/freedommmoto/test_simplebank/tool"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type makeNewUser struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type UserResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
}

type CreateUserParams struct {
	Username       string `json:"username"`
	HashedPassword string `json:"hashed_password"`
	FullName       string `json:"full_name"`
	Email          string `json:"email"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	User        UserResponse `json:"user"`
}

func (server *Server) makeNewUser(ctx *gin.Context) {
	var req makeNewUser
	if err := ctx.ShouldBindJSON(&req); err != nil {
		//same like customer return bad request in case not pass validation
		ctx.JSON(http.StatusBadRequest, errerrorReturn(err))
		return
	}

	hashPassword, err := tool.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashPassword,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errerrorReturn(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}

	ctx.JSON(http.StatusOK, tranFerUserToUserResponse(user))
}

func tranFerUserToUserResponse(user db.User) UserResponse {
	newUser := UserResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
	}
	return newUser
}

func (server *Server) loginProcess(ctx *gin.Context) {
	var req LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, err)
	}
	userFromDB, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errerrorReturn(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
	}

	err = tool.CheckPassword(req.Password, userFromDB.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errerrorReturn(err))
		return
	}

	accessToken, err := server.tokenMaker.CreateToken(userFromDB.Username, server.config.TokenLiftTimeConfig)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errerrorReturn(err))
		return
	}
	passLoginResponse := LoginResponse{
		AccessToken: accessToken,
		User:        tranFerUserToUserResponse(userFromDB),
	}
	ctx.JSON(http.StatusOK, passLoginResponse)
}
