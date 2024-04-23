package api

import (
	"github.com/gin-gonic/gin"
)

// LoadUserRoutes load all the user routes
func LoadUserRoutes(router *gin.Engine, server *Server) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/", server.createUser)
		userGroup.POST("/login", server.login)
		userGroup.POST("/tokens/renew", server.renewAccessToken)
	}
}

// LoadAccountRoutes load all the account routes
func LoadAccountRoutes(router *gin.Engine, server *Server) {
	accountGroup := router.Group("/accounts")
	accountGroup.Use(authMiddleware(server.tokenMaker))
	{
		accountGroup.POST("/", server.createAccount)
		accountGroup.GET("/:id", server.getAccount)
		accountGroup.GET("/", server.listAccounts)
	}
}

// LoadTransferRoutes load all the transfer routes
func LoadTransferRoutes(router *gin.Engine, server *Server) {
	transferGroup := router.Group("/transfers")
	transferGroup.Use(authMiddleware(server.tokenMaker))
	{
		transferGroup.POST("/", server.createTransfer)
	}
}
