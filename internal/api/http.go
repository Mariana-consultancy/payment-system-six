package api

import (
	"fmt"
	"payment-system-six/internal/models"
	"payment-system-six/internal/ports"

	"github.com/gin-gonic/gin"
)

type HTTPHandler struct {
	Repository ports.Repository
}

func NewHTTPHandler(repository ports.Repository) *HTTPHandler {
	return &HTTPHandler{
		Repository: repository,
	}
}

func (u *HTTPHandler) GetUserFromContext(c *gin.Context) (*models.User, error) {
	contextUser, exists := c.Get("user")
	if !exists {
		return nil, fmt.Errorf("error getting user from context")
	}
	user, ok := contextUser.(*models.User)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return user, nil
}

func (u *HTTPHandler) GetAdminFromContext(c *gin.Context) (*models.Admin, error) {
	contextAdmin, exists := c.Get("admin")
	if !exists {
		return nil, fmt.Errorf("error getting admin from context")
	}
	admin, ok := contextAdmin.(*models.Admin)
	if !ok {
		return nil, fmt.Errorf("an error occurred")
	}
	return admin, nil
}

/*func (u *HTTPHandler) GetTokenFromContext(c *gin.Context) (string, error) {
	tokenI, exists := c.Get("access_token")
	if !exists {
		return "", fmt.Errorf("error getting access token")
	}
	tokenstr := tokenI.(string)
	return tokenstr, nil
} */
