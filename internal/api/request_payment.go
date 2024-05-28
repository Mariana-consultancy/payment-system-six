package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Request funds
func (u *HTTPHandler) RequestPayment(c *gin.Context) {
	var requestPayment models.RequestPayment
	var paymentRequest models.PaymentRequests
	if err := c.ShouldBind(&requestPayment); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	requestPayment.RecipientEmail = strings.TrimSpace(requestPayment.RecipientEmail)

	if requestPayment.RecipientEmail == "" {
		util.Response(c, "Recipient email must not be empty", 400, nil, nil)
		return
	}

	if requestPayment.Amount <= 0 {
		util.Response(c, "Amount must be greater than zero", 400, nil, nil)
		return
	}

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	recipient, err := u.Repository.FindUserByEmail(requestPayment.RecipientEmail)
	if err != nil {
		util.Response(c, "Recipient not found. Please enter valid Recipient email address", 404, err.Error(), nil)
		return
	}

	if user.Email == recipient.Email {
		util.Response(c, "Request unsuccessful due to requester account and recipient account are smae.", 400, nil, nil)
		return
	}

	paymentRequest.RequesterID = user.ID
	paymentRequest.RecipientID = recipient.ID
	paymentRequest.Amount = requestPayment.Amount
	paymentRequest.Desscription = requestPayment.Desscription
	paymentRequest.Status = "Pending"

	err = u.Repository.RequestFunds(&paymentRequest)
	if err != nil {
		util.Response(c, "Payment request unsuccessful", 400, err.Error(), nil)
		return
	}
	util.Response(c, "Payment request sent successfully", 200, nil, nil)

}
