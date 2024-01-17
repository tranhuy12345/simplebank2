package api

import (
	"database/sql"
	db "db/db/sqlc"
	"db/token"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createAccountRequest struct {
	//Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

// Xu ly create Account
func (s *Server) createAccount(c *gin.Context) {
	//Lấy dữ liệu từ client
	var req createAccountRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	//Insert vào database
	//Vì có middleware nên owner phải lấy từ payload ra để user nào thì chỉ tạo được account của user đó
	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.CreateAccountParams{
		Owner:    authorPayLoad.Username,
		Currency: req.Currency,
		Balance:  0,
	}
	accounts, err := s.store.CreateAccount(c, arg)
	if err != nil {
		fmt.Println(err)
		pgErr, ok := err.(*pq.Error)
		if ok {
			switch pgErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				c.JSON(http.StatusForbidden, errResponse(pgErr))
			}
		}
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	c.JSON(http.StatusOK, accounts)
}

type getAccountRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (s *Server) getAccount(c *gin.Context) {
	var req getAccountRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authorPayLoad.Username {
		err := errors.New("Account is not belong to user is logged")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	c.JSON(http.StatusOK, account)

}

type listAccountRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (s *Server) listAccount(c *gin.Context) {
	var req listAccountRequest
	err := c.ShouldBindQuery(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsParams{
		Owner:  authorPayLoad.Username,
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	accounts, err := s.store.ListAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	c.JSON(http.StatusOK, accounts)

}

type updateAccountRequest struct {
	Owner    string `json:"owner"`
	Currency string `json:"currency" binding:"oneof=USD EUR"`
	Balance  int64  `json:"balance"`
}

func (s *Server) updateAccount(c *gin.Context) {
	var req getAccountRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := s.store.GetAccountForUpdate(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authorPayLoad.Username {
		err := errors.New("Account is not belong to user is logged")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	var body updateAccountRequest
	err = c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}
	arg := db.UpdateAccountsParams{
		Balance: func() int64 {
			if body.Balance != 0 {
				return body.Balance
			}
			return account.Balance
		}(),
		Owner: func() string {
			if body.Owner != "" {
				return body.Owner
			}
			return account.Owner
		}(),
		Currency: body.Currency,
		ID:       account.ID,
	}
	accountUpdated, err := s.store.UpdateAccounts(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}

	c.JSON(http.StatusCreated, accountUpdated)
}

func (s *Server) deleteAccount(c *gin.Context) {
	var req getAccountRequest
	err := c.ShouldBindUri(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	account, err := s.store.GetAccount(c, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, errResponse(err))
			return
		}

		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	authorPayLoad := c.MustGet(authorizationPayloadKey).(*token.Payload)
	if account.Owner != authorPayLoad.Username {
		err := errors.New("Account is not belong to user is logged")
		c.JSON(http.StatusBadRequest, errResponse(err))
		return
	}

	arg := db.DeleteAccountsTXParams{
		ID: account.ID,
	}
	err = s.store.DeleteAccountsTX(c, arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errResponse(err))
		return
	}
	c.JSON(http.StatusAccepted, account)
}
