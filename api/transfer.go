package api

import (
	"database/sql"
	db "db/db/sqlc"
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
	if !s.validAccount(c, req.FromAccountID, req.Currency) {
		return
	}
	if !s.validAccount(c, req.ToAccountID, req.Currency) {
		return
	}
	//Insert vào database
	arg := db.TransfersTxParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	results, err := s.store.TransfersTX(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	c.JSON(http.StatusOK, results)
}

func (s *Server) validAccount(c *gin.Context, accountID int64, currency string) bool {
	account, err := s.store.GetAccount(c, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return false
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return false
	}

	if account.Currency != currency {
		err := fmt.Errorf("account [%d] currency mismatch: %s vs %s", accountID, account.Currency, currency)
		c.JSON(http.StatusBadRequest, errResponse(err))
		return false
	}
	return true
}
