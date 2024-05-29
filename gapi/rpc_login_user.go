package gapi

import (
	"context"
	"database/sql"
	db "db/db/sqlc"
	"db/db/util"
	"db/pb"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) LoginUser(c context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := s.store.GetUser(c, req.GetUsername())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user %s not found", req.Username)
		}
		return nil, status.Errorf(codes.Internal, "failed to get user %s", err)
	}
	err = util.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "password is not wrong")
	}

	//Tao access Token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(req.GetUsername(), s.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed %s", err)
	}

	//Tao refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		req.GetUsername(),
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed %s", err)
	}
	//fmt.Println("Refresh token", time.Unix(refreshPayload.ExpiredAt, 0))
	//Tao sessions

	metadata := s.extractMetadata(c)

	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    metadata.UserAgent,
		ClientIp:     metadata.ClientIp,
		IsBlocked:    false,
		ExpiresAt:    time.Unix(refreshPayload.ExpiredAt, 0),
	}

	sessions, err := s.store.CreateSession(c, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "sessions is %s", err)
	}
	res := &pb.LoginUserResponse{
		SessionId:             sessions.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(time.Unix(accessPayload.ExpiredAt, 0)),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(sessions.ExpiresAt),
		User:                  converterUser(user),
	}
	return res, nil
}
