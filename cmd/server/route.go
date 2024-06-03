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

	// AuthorizeAdmin authorizes all the authorized users haldlers
	user.Use(middleware.AuthorizeAdmin(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		user.POST("/addfunds", handler.AddFunds)
		user.POST("/transferpayment", handler.TransferPayment)
		user.POST("/requestpayment", handler.RequestPayment)
		user.GET("/getallpaymentrequests", handler.GetAllPaymentRequests)
		user.POST("/getpaymentrequest", handler.GetPaymentRequest)
		user.PUT("/approvepaymentrequest", handler.ApprovePaymentRequest)
		user.PUT("/declinepaymentrequest", handler.DeclinePaymentRequest)
		user.DELETE("/deletepaymentrequest", handler.DeletePaymentRequest)
		user.GET("/getallnotifications", handler.GetAllNotifications)
		user.POST("/getnotification", handler.GetNotification)
		user.PUT("/readnotification", handler.ReadNotification)
		user.PUT("/readAllnotifications", handler.ReadAllNotifications)
		user.DELETE("/deletenotification", handler.DeleteNotification)
		user.DELETE("/deleallnotifications", handler.DeleteAllNotifications)
	}

	// AuthorizeAdmin authorizes all the authorized users haldlers
	authorizeAdmin := r.Group("/admin")
	{
		authorizeAdmin.POST("/create", handler.CreateAdmin)
		authorizeAdmin.POST("/login", handler.LoginAdmin)
	}
	authorizeAdmin.Use(middleware.AuthorizeAdmin(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		authorizeAdmin.GET("/user", handler.GetUserByEmail)
	}

	return router
}
