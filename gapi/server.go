package gapi

import (
	db "db/db/sqlc"
	"db/db/util"
	"db/pb"
	"db/token"
	"fmt"
)

// Server gRPC
type Server struct {
	pb.UnimplementedSimpleBankServer
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
}

// Tao 1 new server gRPC
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
	return server, nil
}
