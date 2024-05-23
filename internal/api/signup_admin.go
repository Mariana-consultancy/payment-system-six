package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Create Admin
func (u *HTTPHandler) CreateAdmin(c *gin.Context) {
	var admin *models.Admin
	if err := c.ShouldBind(&admin); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	admin.FirstName = strings.TrimSpace(admin.FirstName)
	admin.Lastname = strings.TrimSpace(admin.Lastname)
	admin.Email = strings.TrimSpace(admin.Email)
	admin.Password = strings.TrimSpace(admin.Password)
	admin.DateOfBirth = strings.TrimSpace(admin.DateOfBirth)
	admin.Phone = strings.TrimSpace(admin.Phone)
	admin.Address = strings.TrimSpace(admin.Address)
	if admin.FirstName == "" {
		util.Response(c, "First name must not be empty", 400, nil, nil)
		return
	}
	if admin.Lastname == "" {
		util.Response(c, "Last name must not be empty", 400, nil, nil)
		return
	}
	if admin.Email == "" {
		util.Response(c, "Email must not be empty", 400, nil, nil)
		return
	}
	if admin.Password == "" {
		util.Response(c, "Password must not be empty", 400, nil, nil)
		return
	}
	if admin.DateOfBirth == "" {
		util.Response(c, "Date of birth must not be empty", 400, nil, nil)
		return
	}
	if admin.Phone == "" {
		util.Response(c, "Phone must not be empty", 400, nil, nil)
		return
	}
	if admin.Address == "" {
		util.Response(c, "Address must not be empty", 400, nil, nil)
		return
	}

	isEmailExist, _ := u.Repository.FindAdminByEmail(admin.Email)
	if isEmailExist != nil {
		util.Response(c, "Email already exist", 400, nil, nil)
		return
	}

	if !util.ValidatePassword(admin.Password) {
		util.Response(c, "Password acceptence criteria not matched. Password must be At least 6 characters long , Contains at least one uppercase letter, Contains at least one number, Contains at least one special character", 400, nil, nil)
		return
	}

	// Hash the password
	hashedPassword, err := util.HashPassword(admin.Password)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	admin.Password = hashedPassword

	err = u.Repository.CreateAdmin(admin)
	if err != nil {
		util.Response(c, "Admin not created", 400, err.Error(), nil)
		return
	}
	util.Response(c, "Admin created", 200, nil, nil)

}
