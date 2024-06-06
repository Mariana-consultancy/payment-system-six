package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"time"

	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) GetAllPaymentRequests(c *gin.Context) {

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	paymentRequests, err := u.Repository.GetAllPaymentRequestsByAccountNumber(user.AccountNumber)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	if len(*paymentRequests) == 0 {
		util.Response(c, "No record found", 200, paymentRequests, nil)
		return
	}
	util.Response(c, "Record found", 200, paymentRequests, nil)

}

func (u *HTTPHandler) GetPaymentRequest(c *gin.Context) {
	var getPaymentRequest models.GetPaymentRequest

	if err := c.ShouldBind(&getPaymentRequest); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	paymentRequest, err := u.Repository.GetPaymentRequestByRequestID(getPaymentRequest.PaymentRequestID)
	if err != nil {
		util.Response(c, "Payment request ID does not exist", 400, err.Error(), nil)
		return
	}

	util.Response(c, "Record found", 200, paymentRequest, nil)

}

func (u *HTTPHandler) ApprovePaymentRequest(c *gin.Context) {
	var getPaymentRequest models.GetPaymentRequest
	var notification models.Notification
	var payerTransaction models.Transaction
	var payeeTransaction models.Transaction

	if err := c.ShouldBind(&getPaymentRequest); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	paymentRequest, err := u.Repository.GetPaymentRequestByRequestID(getPaymentRequest.PaymentRequestID)
	if err != nil {
		util.Response(c, "Payment request ID does not exist", 400, err.Error(), nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	if user.AccountNumber == paymentRequest.RequesterAccountNumber {
		util.Response(c, "You can not approve this request because you are the requester of this request", 400, nil, nil)
		return
	}

	if paymentRequest.IsCompleted {
		util.Response(c, "You can not change this request because this request is already completed", 400, nil, nil)
		return
	}

	requester, err := u.Repository.GetUserByAccountNumber(paymentRequest.RequesterAccountNumber)
	if err != nil {
		util.Response(c, "Requester not found", 400, err.Error(), nil)
		return
	}

	if user.Balance-paymentRequest.Amount < 0 {
		util.Response(c, "Payment unsuccessful due to insufficient funds.", 400, nil, nil)
		return
	}

	payeeTransaction.BalanceBefore = requester.Balance
	requester.Balance += paymentRequest.Amount
	err = u.Repository.UpdateUser(requester)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	payerTransaction.BalanceBefore = user.Balance
	user.Balance -= paymentRequest.Amount
	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	payerTransaction.BalanceAfter = user.Balance
	payerTransaction.AccountNumber = user.AccountNumber
	payerTransaction.PayerAccountNumber = user.AccountNumber
	payerTransaction.PayeeAccountNumber = requester.AccountNumber
	payerTransaction.TransactionType = "Debit"
	payerTransaction.TransactionAmount = paymentRequest.Amount
	payerTransaction.TransactionDate = time.Now()

	payeeTransaction.BalanceAfter = requester.Balance
	payeeTransaction.AccountNumber = requester.AccountNumber
	payeeTransaction.PayerAccountNumber = user.AccountNumber
	payeeTransaction.PayeeAccountNumber = requester.AccountNumber
	payeeTransaction.TransactionType = "Credit"
	payeeTransaction.TransactionAmount = paymentRequest.Amount
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

	paymentRequest.Status = "Approved"
	paymentRequest.IsCompleted = true
	err = u.Repository.UpdatePaymentRequest(paymentRequest)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	notification.ReceiverID = requester.ID
	notification.SenderID = user.ID
	notification.Title = "Payment Request approved"
	notification.Message = "Your Payment Request has been approved by " + user.FirstName + " " + user.Lastname + " and added to your account balance."
	notification.NotificationType = "Payment Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Payment request approved successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}

	notification = models.Notification{}
	notification.ReceiverID = user.ID
	notification.SenderID = user.ID
	notification.Title = "Funds Transfered"
	notification.Message = "Your account has been Debit."
	notification.NotificationType = "Payment Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Payment request approved successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}
	util.Response(c, "Request approved successfully", 200, nil, nil)

}

func (u *HTTPHandler) DeclinePaymentRequest(c *gin.Context) {
	var getPaymentRequest models.GetPaymentRequest
	var notification models.Notification

	if err := c.ShouldBind(&getPaymentRequest); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	paymentRequest, err := u.Repository.GetPaymentRequestByRequestID(getPaymentRequest.PaymentRequestID)
	if err != nil {
		util.Response(c, "Payment request ID does not exist", 400, err.Error(), nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	if user.AccountNumber == paymentRequest.RequesterAccountNumber {
		util.Response(c, "You can not decline this request because you are the requester of this request", 400, nil, nil)
		return
	}

	if paymentRequest.IsCompleted {
		util.Response(c, "You can not change this request because this request is already completed", 400, nil, nil)
		return
	}

	requester, err := u.Repository.GetUserByAccountNumber(paymentRequest.RequesterAccountNumber)
	if err != nil {
		util.Response(c, "Requester not found", 400, err.Error(), nil)
		return
	}

	paymentRequest.Status = "Declined"
	paymentRequest.IsCompleted = true
	err = u.Repository.UpdatePaymentRequest(paymentRequest)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	notification.ReceiverID = requester.ID
	notification.SenderID = user.ID
	notification.Title = "Payment Request declined"
	notification.Message = "Your Payment Request has been declined by " + user.FirstName + " " + user.Lastname
	notification.NotificationType = "Payment Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Payment request declined successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}
	util.Response(c, "Request declined successfully", 200, nil, nil)

}

func (u *HTTPHandler) DeletePaymentRequest(c *gin.Context) {
	var getPaymentRequest models.GetPaymentRequest
	var notification models.Notification

	if err := c.ShouldBind(&getPaymentRequest); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	paymentRequest, err := u.Repository.GetPaymentRequestByRequestID(getPaymentRequest.PaymentRequestID)
	if err != nil {
		util.Response(c, "Payment request ID does not exist", 400, err.Error(), nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	if user.AccountNumber == paymentRequest.RecipientAccountNumber {
		util.Response(c, "You can not delete this request because you are not requester of this request", 400, nil, nil)
		return
	}

	if paymentRequest.IsCompleted {
		util.Response(c, "You can not delete this request because this request is already completed", 400, nil, nil)
		return
	}

	recipient, err := u.Repository.GetUserByAccountNumber(paymentRequest.RecipientAccountNumber)
	if err != nil {
		util.Response(c, "Requester not found", 400, err.Error(), nil)
		return
	}

	err = u.Repository.DeletePaymentRequest(paymentRequest)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	notification.ReceiverID = recipient.ID
	notification.SenderID = user.ID
	notification.Title = "Payment Request Deleted"
	notification.Message = "Payment Request has been deleted by " + user.FirstName + " " + user.Lastname
	notification.NotificationType = "Payment Request"
	notification.Status = "Unread"

	err = u.Repository.CreateNotification(&notification)
	if err != nil {
		util.Response(c, "Payment request deleted successfully but Notification unsuccessful", 200, err.Error(), nil)
		return
	}
	util.Response(c, "Payment request deleted successfully", 200, nil, nil)

}
