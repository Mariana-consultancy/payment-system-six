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
		r.POST("/create", handler.CreateUser)
		r.POST("/login", handler.LoginUser)
	}

	// AuthorizeAdmin authorizes all the authorized users haldlers
	authorizeAdmin := r.Group("/admin")
	authorizeAdmin.Use(middleware.AuthorizeAdmin(repository.FindUserByEmail, repository.TokenInBlacklist))
	{
		authorizeAdmin.GET("/user", handler.GetUserByEmail)
	}

	return router
}
