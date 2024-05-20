package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Create User
func (u *HTTPHandler) CreateUser(c *gin.Context) {
	var user *models.User
	if err := c.ShouldBind(&user); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	if strings.TrimSpace(user.Email) == "" || strings.TrimSpace(user.Password) == "" {
		util.Response(c, "Email and Password must not be empty", 400, nil, nil)
		return
	}

	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	isEmailExist, _ := u.Repository.FindUserByEmail(user.Email)
	if isEmailExist != nil {
		util.Response(c, "Email already exist", 400, nil, nil)
		return
	}

	if !util.ValidatePassword(user.Password) {
		util.Response(c, "Password acceptence criteria not matched. Password must be At least 6 characters long , Contains at least one uppercase letter, Contains at least one number, Contains at least one special character", 400, nil, nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.Password = hashedPassword

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}
