package api

import (
	"errors"
	"net/http"

	db "github.com/daniel-vuky/golang-bank-app/db/sqlc"
	"github.com/daniel-vuky/golang-bank-app/token"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// getAccountRequest defines the request body for createAccount handler
type createAccountRequest struct {
	Currency string `json:"currency" binding:"required,currency"`
}

// Server serves HTTP create account requests
func (server *Server) createAccount(c *gin.Context) {
	var req createAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Balance:  0,
		Currency: req.Currency,
	}
	account, err := server.store.CreateAccount(c, arg)
	if err != nil {
		if pqErr, ok := err.(*pgconn.PgError); ok {
			switch pqErr.Code {
			case "23503", "23505":
				c.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		c.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)
}

// getAccountRequest defines the request body for getAccount handler
type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

// Server serves HTTP get account requests
func (server *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	account, getAccountErr := server.store.GetAccount(c, req.ID)
	if getAccountErr != nil {
		if errors.Is(getAccountErr, pgx.ErrNoRows) {
			c.JSON(http.StatusNotFound, errorResponse(getAccountErr))
			return
		}
		c.JSON(http.StatusInternalServerError, errorResponse(getAccountErr))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authPayload.Username {
		err := errors.New("account not belong to authenticated user")
		c.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)
}

// listAccountRequest defines the request body for listAccounts handler
type listAccountRequest struct {
	CurrentPage int32 `form:"current_page" binding:"required,min=1"`
	PageSize    int32 `form:"page_size" binding:"required,min=5,max=10"`
}

// Server serves HTTP list accounts requests
func (server *Server) listAccounts(c *gin.Context) {
	var req listAccountRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authPayload.Username,
		Limit:  req.PageSize,
		Offset: (req.CurrentPage - 1) * req.PageSize,
	}
	listAccount, listAccountErr := server.store.ListAccounts(c, arg)
	if listAccountErr != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(listAccountErr))
		return
	}
	c.JSON(http.StatusOK, listAccount)
}
