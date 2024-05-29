package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type renewAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type renewAccessTokenResponse struct {
	AccessToken          string    `json:"access_token"`
	AccessTokenExpiresAt time.Time `json:"access_token_expires_at"`
}

func (s *Server) renewAccessToken(c *gin.Context) {
	var req renewAccessTokenRequest
	var res renewAccessTokenResponse
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	refreshPayload, err := s.tokenMaker.VerifyToken(req.RefreshToken)

	if err != nil {
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	session, err := s.store.GetSession(c, refreshPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	if session.IsBlocked {
		err := fmt.Errorf("block session")
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if session.Username != refreshPayload.Username {
		err := fmt.Errorf("username mismatch")
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if req.RefreshToken != session.RefreshToken {
		err := fmt.Errorf("refresh token mismatch")
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	if time.Now().After(session.ExpiresAt) {
		err := fmt.Errorf("refresh token expired")
		c.JSON(http.StatusUnauthorized, errResponse(err))
		return
	}

	//Tao access Token
	accessToken, accessPayload, err := s.tokenMaker.CreateToken(session.Username, s.config.AccessTokenDuration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	res = renewAccessTokenResponse{
		AccessToken:          accessToken,
		AccessTokenExpiresAt: time.Unix(accessPayload.ExpiredAt, 0),
	}
	c.JSON(http.StatusOK, res)
}
