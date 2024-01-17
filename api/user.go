package api

import (
	"database/sql"
	db "db/db/sqlc"
	"db/db/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
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

	arg := db.CreateUserParams{
		Username:     req.Username,
		HashPassword: hashedPassword,
		FullName:     req.FullName,
		Email:        req.Email,
	}

	user, err := s.store.CreateUser(c, arg)
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
	AccessToken  string       `json:"access_token"`
	UserResponse userResponse `json:"user"`
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

	token, err := s.tokenMaker.CreateToken(req.Username, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	userResponse := newUsersResponse(user)
	res = loginUserResponse{
		AccessToken:  token,
		UserResponse: userResponse,
	}
	c.JSON(http.StatusOK, res)
}
