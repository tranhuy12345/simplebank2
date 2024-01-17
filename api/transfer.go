package api

import (
	"database/sql"
	db "db/db/sqlc"
	"db/token"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type transferRequest struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"` //gt greater than
	Currency      string `json:"currency" binding:"required,currency"`
}

// Xu ly create Transfer
func (s *Server) createTransfers(c *gin.Context) {
	//Lấy dữ liệu từ client
	var req transferRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//Check currency
	from_account, isValid := s.validAccount(c, req.FromAccountID, req.Currency)
	if !isValid {
		return
	}

	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if from_account.Owner != authorPayLoad.Username {
		err := errors.New("Account is not belong to user is logged")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	to_account, isValid := s.validAccount(c, req.ToAccountID, req.Currency)
	if !isValid {
		return
	}
	//Insert vào database
	arg := db.TransfersTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   to_account.ID,
		Amount:        req.Amount,
	}
	results, err := s.store.TransfersTX(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string) (db.Accounts, bool) {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return account, false
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return account, false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errResponse(err))
		return account, false
	}
	return account, true
}
