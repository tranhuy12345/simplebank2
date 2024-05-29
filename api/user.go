package api

import (
	"database/sql"
	db "db/db/sqlc"
	"db/db/util"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (s *Server) createUser(c *gin.Context) {
	var req createUserRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	fmt.Println("Vo1")
	arg := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: hashedPassword,
		FullName:     req.FullName,
		Email:        req.Email,
	}

	user, err := s.store.CreateUser(c, arg)
	fmt.Println("Vo2")
	if err != nil {
		errPq, ok := err.(*pq.Error)
		if ok {
			c.JSON(http.StatusInternalServerError, errResponse(errPq))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	responseUser := newUsersResponse(user)

	c.JSON(http.StatusOK, responseUser)
}

func newUsersResponse(user db.Users) userResponse {
	responseUser := userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt.Time,
	}
	return responseUser
}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionId             uuid.UUID    `json:"session_id"`
	AccessToken           string       `json:"access_token"`
	AccessTokenExpiresAt  time.Time    `json:"access_token_expires_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time    `json:"refresh_token_expires_at"`
	UserResponse          userResponse `json:"user"`
}

func (s *Server) login(c *gin.Context) {
	var req loginUserRequest
	var res loginUserResponse
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	user, err := s.store.GetUser(c, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	//Tao access Token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	//Tao refresh token
	refreshToken, refreshPayload, err := s.tokenMaker.CreateToken(
		req.Username,
		s.config.RefreshTokenDuration,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	//fmt.Println("Refresh token", time.Unix(refreshPayload.ExpiredAt, 0))
	//Tao sessions

	arg := db.CreateSessionParams{
		ID:           refreshPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    c.Request.UserAgent(),
		ClientIp:     c.ClientIP(),
		IsBlocked:    false,
		ExpiresAt:    time.Unix(refreshPayload.ExpiredAt, 0),
	}

	sessions, err := s.store.CreateSession(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	userResponse := newUsersResponse(user)
	res = loginUserResponse{
		SessionId:             sessions.ID,
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  time.Unix(accessPayload.ExpiredAt, 0),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: sessions.ExpiresAt,
		UserResponse:          userResponse,
	}
	c.JSON(http.StatusOK, res)
}

// Path: api/user.go
