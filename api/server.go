package api

import (
	db "db/db/sqlc"
	"db/db/util"
	"db/token"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// Tao 1 new server
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("Don't create token maker %w", err)
	}

	server := &Server{
		store:      store,
		config:     config,
		tokenMaker: tokenMaker,
	}
	server.setupRouter()
	validator, ok := binding.Validator.Engine().(*validator.Validate)
	if ok {
		validator.RegisterValidation("currency", validCurrency)
	}

	return server, nil
}

// Setup cac route
func (s *Server) setupRouter() {
	router := gin.Default()
	router.POST("/users/login", s.login)
	router.POST("/users", s.createUser)

	authRoutes := router.Group("/").Use(authMiddleware(s.tokenMaker))

	authRoutes.POST("/accounts", s.createAccount)
	authRoutes.GET("/accounts/:id", s.getAccount)
	authRoutes.GET("/accounts", s.listAccount)
	authRoutes.PUT("/accounts/:id", s.updateAccount)
	authRoutes.DELETE("/accounts/:id", s.deleteAccount)

	authRoutes.POST("/transfers", s.createTransfers)
	s.router = router
}

// Start server
func (server *Server) Start(address string) error {
	return server.router.Run(address)

}

// Funtion tra ve loi
func errResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
