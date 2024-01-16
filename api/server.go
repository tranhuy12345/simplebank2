package api

import (
	db "db/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

// Tao 1 new server
func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	validator, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validator.RegisterValidation("currency", validCurrency)
	}

	router.POST("/users", server.createUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)
	router.PUT("/accounts/:id", server.updateAccount)
	router.DELETE("/accounts/:id", server.deleteAccount)

	router.POST("/transfers", server.createTransfers)
	server.router = router
	return server
}

// Start server
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

// Funtion tra ve loi
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
