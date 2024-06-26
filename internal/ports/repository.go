package ports

import "payment-system-six/internal/models"

type Repository interface {
	GenerateUserAccountNumber() (uint, error)
	FindUserByEmail(email string) (*models.User, error)
	GetUserByID(userID uint) (*models.User, error)
	FindAllUsers() ([]models.User, error)
	GetUserByAccountNumber(accountNumber uint) (*models.User, error)
	FindAdminByEmail(email string) (*models.Admin, error)
	CreateUser(user *models.User) error
	CreateAdmin(admin *models.Admin) error
	UpdateUser(user *models.User) error
	UpdateAdmin(user *models.Admin) error
	ValidateCard(cardNumber, cardExpiry, cardCVV string) error
	RecordTransaction(transaction *models.Transaction) error
	RequestFunds(paymentRequest *models.PaymentRequests) error
	GetAllPaymentRequestsByAccountNumber(accountNumber uint) (*[]models.PaymentRequests, error)
	GetPaymentRequestByRequestID(requestID uint) (*models.PaymentRequests, error)
	UpdatePaymentRequest(paymentRequest *models.PaymentRequests) error
	DeletePaymentRequest(paymentRequest *models.PaymentRequests) error
	CreateNotification(notification *models.Notification) error
	GetNotificationsByUserID(userID uint) (*models.NotificationDetails, error)
	GetNotificationByNotificationID(notificationID uint) (*models.Notification, error)
	UpdateNotification(notification *models.Notification) error
	UpdateAllNotificationsByUserID(userID uint) error
	DeleteNotification(notification *models.Notification) error
	DeleteAllNotificationByUserID(userID uint) error
	DepositFunds(depositRequest *models.DepositRequests) error
	GetAllDepositRequests() ([]models.DepositRequests, error)
	GetAllDepositRequestsByAccountNumber(accountNumber uint) (*[]models.DepositRequests, error)
	GetDepositRequestByRequestID(requestID uint) (*models.DepositRequests, error)
	UpdateDepositRequest(depositRequest *models.DepositRequests) error
	GenerateStatement(accountNumber uint) (*models.StatementDetails, error)
	BlacklistToken(token *models.BlacklistTokens) error
	TokenInBlacklist(token *string) bool
}
