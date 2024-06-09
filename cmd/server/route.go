package server

import (
	"payment-system-six/internal/api"
	"payment-system-six/internal/middleware"
	"payment-system-six/internal/ports"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter is where router endpoints are called
func SetupRouter(handler *api.HTTPHandler, repository ports.Repository) *gin.Engine {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST", "GET", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r := router.Group("/")
	{
		r.GET("/", handler.Readiness)
	}

	user := r.Group("/user")
	{
		user.POST("/create", handler.CreateUser)
		user.POST("/login", handler.LoginUser)
	}

	// AuthorizeUser authorizes all the authorized users haldlers
	user.Use(middleware.AuthorizeUser(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		user.POST("/addfunds", handler.AddFunds)
		user.POST("/depositfunds", handler.DepositFunds)
		user.POST("/transferpayment", handler.TransferPayment)
		user.POST("/requestpayment", handler.RequestPayment)
		user.GET("/getallpaymentrequests", handler.GetAllPaymentRequests)
		user.POST("/getpaymentrequest", handler.GetPaymentRequest)
		user.PUT("/approvepaymentrequest", handler.ApprovePaymentRequest)
		user.PUT("/declinepaymentrequest", handler.DeclinePaymentRequest)
		user.DELETE("/deletepaymentrequest", handler.DeletePaymentRequest)
		user.POST("/generatestatement", handler.GenerateStatementUser)
		user.GET("/getallnotifications", handler.GetAllNotifications)
		user.POST("/getnotification", handler.GetNotification)
		user.PUT("/readnotification", handler.ReadNotification)
		user.PUT("/readAllnotifications", handler.ReadAllNotifications)
		user.DELETE("/deletenotification", handler.DeleteNotification)
		user.DELETE("/deleallnotifications", handler.DeleteAllNotifications)
		user.POST("/logout", handler.Logout)
	}

	// AuthorizeAdmin authorizes all the authorized users haldlers
	Admin := r.Group("/admin")
	{
		Admin.POST("/create", handler.CreateAdmin)
		Admin.POST("/login", handler.LoginAdmin)
	}
	Admin.Use(middleware.AuthorizeAdmin(repository.FindAdminByEmail, repository.TokenInBlacklist))
	{
		Admin.GET("/getallusers", handler.GetAllUsers)
		Admin.GET("/getalldepositrequests", handler.GetAllDepositRequests)
		Admin.POST("/getalldepositrequestsbyaccountnumber", handler.GetAllDepositRequestsByAccountNumber)
		Admin.POST("/getdepositrequestbyrequestid", handler.GetDepositRequestByRequestID)
		Admin.PUT("/approvedepositrequest", handler.ApproveDepositRequest)
		Admin.PUT("/declinedepositrequest", handler.DeclineDepositRequest)
		Admin.POST("/generatestatement", handler.GenerateStatementAdmin)
		Admin.POST("/logout", handler.Logout)
	}

	return router
}
