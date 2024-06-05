package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) GetAllDepositRequests(c *gin.Context) {

	depositRequests, err := u.Repository.GetAllDepositRequests()
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	if len(depositRequests) == 0 {
		util.Response(c, "No record found", 200, depositRequests, nil)
		return
	}
	util.Response(c, "Record found", 200, depositRequests, nil)

}

func (u *HTTPHandler) GetAllDepositRequestsByAccountNumber(c *gin.Context) {
	var userAccountNumber *models.UserAccountNumber
	err := c.ShouldBind(&userAccountNumber)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	depositRequests, err := u.Repository.GetAllDepositRequestsByAccountNumber(userAccountNumber.AccountNumber)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	if len(*depositRequests) == 0 {
		util.Response(c, "No record found", 200, depositRequests, nil)
		return
	}
	util.Response(c, "Record found", 200, depositRequests, nil)

}

func (u *HTTPHandler) GetDepositRequestByRequestID(c *gin.Context) {
	var requestID *models.RequestID
	err := c.ShouldBind(&requestID)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	depositRequests, err := u.Repository.GetDepositRequestByRequestID(requestID.RequestID)
	if err != nil {
		util.Response(c, "Deposit request ID does not exist", 400, err.Error(), nil)
		return
	}

	util.Response(c, "Record found", 200, depositRequests, nil)

}

func (u *HTTPHandler) ApproveDepositRequest(c *gin.Context) {
	var requestID models.RequestID
	var notification models.Notification
	var payeeTransaction models.Transaction

	if err := c.ShouldBind(&requestID); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	depositRequest, err := u.Repository.GetDepositRequestByRequestID(requestID.RequestID)
	if err != nil {
		util.Response(c, "Deposit request ID does not exist", 400, err.Error(), nil)
		return
	}

	user, err := u.GetAdminFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	if depositRequest.IsCompleted {
		util.Response(c, "You can not change this request because this request is already completed", 400, nil, nil)
		return
	}

	depositor, err := u.Repository.GetUserByAccountNumber(depositRequest.DepositorAccountNumber)
	if err != nil {
		util.Response(c, "Depositor not found", 400, err.Error(), nil)
		return
	}

	payeeTransaction.BalanceBefore = depositor.Balance
	depositor.Balance += depositRequest.Amount
	err = u.Repository.UpdateUser(depositor)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	payeeTransaction.BalanceAfter = depositor.Balance
	payeeTransaction.AccountNumber = depositor.AccountNumber
	payeeTransaction.PayerAccountNumber = depositor.AccountNumber
	payeeTransaction.PayeeAccountNumber = depositor.AccountNumber
	payeeTransaction.TransactionType = "Credit"
	payeeTransaction.TransactionAmount = depositRequest.Amount
	payeeTransaction.TransactionDate = time.Now()

	err = u.Repository.RecordTransaction(&payeeTransaction)
	if err != nil {
		util.Response(c, "There is an error occured.", 500, err.Error(), nil)
		return
	}

	depositRequest.Status = "Approved"
	depositRequest.IsCompleted = true
	err = u.Repository.UpdateDepositRequest(depositRequest)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	notification.ReceiverID = depositor.ID
	notification.SenderID = user.ID
	notification.Title = "Deposit Request approved"
	notification.Message = "Your Deposit Request has been approved by Admin and added to your account balance."
	notification.NotificationType = "Deposit Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Deposit request approved successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}
	util.Response(c, "Request approved successfully", 200, nil, nil)

}

func (u *HTTPHandler) DeclineDepositRequest(c *gin.Context) {
	var requestID models.RequestID
	var notification models.Notification

	if err := c.ShouldBind(&requestID); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	depositRequest, err := u.Repository.GetDepositRequestByRequestID(requestID.RequestID)
	if err != nil {
		util.Response(c, "Deposit request ID does not exist", 400, err.Error(), nil)
		return
	}

	user, err := u.GetAdminFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	if depositRequest.IsCompleted {
		util.Response(c, "You can not change this request because this request is already completed", 400, nil, nil)
		return
	}

	depositor, err := u.Repository.GetUserByAccountNumber(depositRequest.DepositorAccountNumber)
	if err != nil {
		util.Response(c, "Depositor not found", 400, err.Error(), nil)
		return
	}

	depositRequest.Status = "Declined"
	depositRequest.IsCompleted = true
	err = u.Repository.UpdateDepositRequest(depositRequest)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	notification.ReceiverID = depositor.ID
	notification.SenderID = user.ID
	notification.Title = "Deposit Request declined"
	notification.Message = "Your Deposit Request has been declined by Admin"
	notification.NotificationType = "Deposit Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Deposit request declined successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}
	util.Response(c, "Request declined successfully", 200, nil, nil)

}
