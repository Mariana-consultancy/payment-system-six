package api

import (
	"os"
	"payment-system-six/internal/middleware"
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"

	"github.com/dgrijalva/jwt-go"
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

	// Generate Account number
	accountNumber, err := u.Repository.GenerateUserAccountNumber()
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	user.AccountNumber = accountNumber

	err = u.Repository.CreateUser(user)
	if err != nil {
		util.Response(c, "User not created", 500, err.Error(), nil)
		return
	}
	util.Response(c, "User created", 200, nil, nil)

}

// Login User
func (u *HTTPHandler) LoginUser(c *gin.Context) {
	var loginRequest *models.LoginRequestUser
	err := c.ShouldBind(&loginRequest)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	loginRequest.Email = strings.TrimSpace(loginRequest.Email)
	loginRequest.Password = strings.TrimSpace(loginRequest.Password)

	if loginRequest.Email == "" {
		util.Response(c, "Email must not be empty", 400, nil, nil)
		return
	}
	if loginRequest.Password == "" {
		util.Response(c, "Password must not be empty", 400, nil, nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(loginRequest.Email)
	if err != nil {
		util.Response(c, "Email does not exist", 404, err.Error(), nil)
		return
	}

	/* if user.IsLocked {

		if time.Now().After(user.LockedAt.Add(time.Hour*2)) || time.Now().Equal(user.LockedAt.Add(time.Hour*2)) {
			user.IsLocked = false
			user.LoginCounter = 0
		} else {
			util.Response(c, "Your Account has been locked. Please try again after 2 hours", 401, nil, nil)
			return
		}

	} */

	/* if user.LoginCounter >= 3 {
		user.IsLocked = true
		user.LockedAt = time.Now()
		err = u.Repository.UpdateUser(user)
		if err != nil {
			util.Response(c, "There is an error occured", 500, err.Error(), nil)
			return
		}

		util.Response(c, "Your Account has been locked after 3 attempts. Please try again after 2 hours", 401, nil, nil)
		return
	} */

	// Verify the password
	match := util.CheckPasswordHash(loginRequest.Password, user.Password)
	if !match {
		//user.LoginCounter++
		err = u.Repository.UpdateUser(user)
		if err != nil {
			util.Response(c, "There is an error occured", 500, err.Error(), nil)
			return
		}
		util.Response(c, "Incorrect password", 401, nil, nil)
		return
	}

	/* user.LoginCounter = 0
	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	} */

	accessClaims, refreshClaims := middleware.GenerateClaims(user.Email)

	secret := os.Getenv("JWT_SECRET")

	accessToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, accessClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating access token", 500, err.Error(), nil)
		return
	}

	refreshToken, err := middleware.GenerateToken(jwt.SigningMethodHS256, refreshClaims, &secret)
	if err != nil {
		util.Response(c, "Error generating refresh token", 500, err.Error(), nil)
		return
	}

	/* notifications, err := u.Repository.GetNotificationsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	} */

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user": user,
		//"notification_details": notifications,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

func (u *HTTPHandler) GenerateStatementUser(c *gin.Context) {

	var userAccountNumber *models.UserAccountNumber
	err := c.ShouldBind(&userAccountNumber)
	if err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	_, err = u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	statement, err := u.Repository.GenerateStatement(userAccountNumber.AccountNumber)
	if err != nil {
		util.Response(c, "Statement not found", 400, err.Error(), nil)
		return
	}

	util.Response(c, "Record found", 200, statement, nil)
}
