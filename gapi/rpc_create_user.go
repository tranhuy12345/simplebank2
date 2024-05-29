package gapi

import (
	"context"
	db "db/db/sqlc"
	"db/db/util"
	"db/pb"
	"fmt"

	"github.com/lib/pq"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) CreateUser(c context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to hash password, %s", err)
	}
	arg := db.CreateUserParams{
		Username:     req.GetUsername(),
		HashPassword: hashedPassword,
		FullName:     req.GetFullName(),
		Email:        req.GetEmail(),
	}

	user, err := s.store.CreateUser(c, arg)
	fmt.Println("Vo2")
	if err != nil {
		errPq, ok := err.(*pq.Error)
		if ok {
			switch errPq.Code.Name() {
			case "unique_violation":
				return nil, status.Errorf(codes.AlreadyExists, "username %s already exists", req.Username)
			default:
				return nil, status.Errorf(codes.Internal, "failed to create user, %s", errPq)
			}
		}
		return nil, status.Errorf(codes.Internal, "failed to create user %s", user)
	}
	responseUser := &pb.CreateUserResponse{
		User: converterUser(user),
	}
	return responseUser, nil

}
