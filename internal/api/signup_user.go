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

	user.FirstName = strings.TrimSpace(user.FirstName)
	user.Lastname = strings.TrimSpace(user.Lastname)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)
	user.DateOfBirth = strings.TrimSpace(user.DateOfBirth)
	user.Phone = strings.TrimSpace(user.Phone)
	user.Address = strings.TrimSpace(user.Address)
	if user.FirstName == "" {
		util.Response(c, "First name must not be empty", 400, nil, nil)
		return
	}
	if user.Lastname == "" {
		util.Response(c, "Last name must not be empty", 400, nil, nil)
		return
	}
	if user.Email == "" {
		util.Response(c, "Email must not be empty", 400, nil, nil)
		return
	}
	if user.Password == "" {
		util.Response(c, "Password must not be empty", 400, nil, nil)
		return
	}
	if user.DateOfBirth == "" {
		util.Response(c, "Date of birth must not be empty", 400, nil, nil)
		return
	}
	if user.Phone == "" {
		util.Response(c, "Phone must not be empty", 400, nil, nil)
		return
	}
	if user.Address == "" {
		util.Response(c, "Address must not be empty", 400, nil, nil)
		return
	}

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
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}
