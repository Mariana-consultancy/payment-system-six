package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// Add funds
func (u *HTTPHandler) AddFunds(c *gin.Context) {
	var fundRequest *models.FundRequest
	var Transaction models.Transaction
	if err := c.ShouldBind(&fundRequest); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	fundRequest.CardNumber = strings.TrimSpace(fundRequest.CardNumber)
	fundRequest.CardExpiry = strings.TrimSpace(fundRequest.CardExpiry)
	fundRequest.CardCVV = strings.TrimSpace(fundRequest.CardCVV)

	if fundRequest.CardNumber == "" {
		util.Response(c, "Card Number must not be empty", 400, nil, nil)
		return
	}
	if fundRequest.CardExpiry == "" {
		util.Response(c, "Card expiry must not be empty", 400, nil, nil)
		return
	}
	if fundRequest.CardCVV == "" {
		util.Response(c, "CVV must not be empty", 400, nil, nil)
		return
	}
	if fundRequest.Amount <= 0 {
		util.Response(c, "Amount must be greater than zero", 400, nil, nil)
		return
	}

	err := u.Repository.ValidateCard(fundRequest.CardNumber, fundRequest.CardExpiry, fundRequest.CardCVV)
	if err != nil {
		util.Response(c, "Invalid Card Details", 400, err.Error(), nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	Transaction.BalanceBefore = user.Balance
	user.Balance += fundRequest.Amount

	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}
	Transaction.BalanceAfter = user.Balance
	Transaction.AccountNumber = user.AccountNumber
	Transaction.PayerAccountNumber = user.AccountNumber
	Transaction.PayeeAccountNumber = user.AccountNumber
	Transaction.TransactionType = "Credit"
	Transaction.TransactionAmount = fundRequest.Amount
	Transaction.TransactionDate = time.Now()

	err = u.Repository.RecordTransaction(&Transaction)
	if err != nil {
		util.Response(c, "There is an error occured.", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Funds added successfully", 200, nil, nil)

}
