package api

import (
	db "db/db/sqlc"
	"db/db/util"
	"net/http"

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
	Username string `json:"username" binding:"required,alphanum"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
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

	responseUser := userResponse{
		Username: user.Username,
		FullName: user.FullName,
		Email:    user.Email,
	}

	c.JSON(http.StatusOK, responseUser)

}
