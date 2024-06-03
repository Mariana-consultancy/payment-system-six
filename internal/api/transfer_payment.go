package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

// Send funds
func (u *HTTPHandler) TransferPayment(c *gin.Context) {
	var transferPayment models.TransferPayment
	var payerTransaction models.Transaction
	var payeeTransaction models.Transaction
	if err := c.ShouldBind(&transferPayment); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if transferPayment.Amount <= 0 {
		util.Response(c, "Amount must be greater than zero", 400, nil, nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	payee, err := u.Repository.GetUserByAccountNumber(transferPayment.AccountNumber)
	if err != nil {
		util.Response(c, "Payee not found. Please enter valid account number", 400, err.Error(), nil)
		return
	}

	if user.AccountNumber == payee.AccountNumber {
		util.Response(c, "Payment unsuccessful due to payer account and payee account are smae.", 400, nil, nil)
		return
	}

	if user.Balance-transferPayment.Amount < 0 {
		util.Response(c, "Payment unsuccessful due to insufficient funds.", 400, nil, nil)
		return
	}

	payeeTransaction.BalanceBefore = payee.Balance
	payee.Balance += transferPayment.Amount
	err = u.Repository.UpdateUser(payee)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	payerTransaction.BalanceBefore = user.Balance
	user.Balance -= transferPayment.Amount
	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	payerTransaction.BalanceAfter = user.Balance
	payerTransaction.AccountNumber = user.AccountNumber
	payerTransaction.PayerAccountNumber = user.AccountNumber
	payerTransaction.PayeeAccountNumber = payee.AccountNumber
	payerTransaction.TransactionType = "Debit"
	payerTransaction.TransactionAmount = transferPayment.Amount
	payerTransaction.TransactionDate = time.Now()

	payeeTransaction.BalanceAfter = payee.Balance
	payeeTransaction.AccountNumber = payee.AccountNumber
	payeeTransaction.PayerAccountNumber = user.AccountNumber
	payeeTransaction.PayeeAccountNumber = payee.AccountNumber
	payeeTransaction.TransactionType = "Credit"
	payeeTransaction.TransactionAmount = transferPayment.Amount
	payeeTransaction.TransactionDate = time.Now()

	err = u.Repository.RecordTransaction(&payerTransaction)
	if err != nil {
		util.Response(c, "There is an error occured.", 500, err.Error(), nil)
		return
	}

	err = u.Repository.RecordTransaction(&payeeTransaction)
	if err != nil {
		util.Response(c, "There is an error occured.", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Funds transfer successfully", 200, nil, nil)

}
