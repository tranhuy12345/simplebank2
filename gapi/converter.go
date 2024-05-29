package gapi

import (
	db "db/db/sqlc"
	"db/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func converterUser(user db.Users) *pb.User {
	return &pb.User{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: timestamppb.New(user.PasswordChangedAt),
		CreatedAt:         timestamppb.New(user.CreatedAt.Time),
	}
}
