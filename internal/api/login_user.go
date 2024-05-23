package api

import (
	"os"
	"payment-system-six/internal/middleware"
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

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

	if user.IsLocked {

		if time.Now().After(user.LockedAt.Add(time.Hour*2)) || time.Now().Equal(user.LockedAt.Add(time.Hour*2)) {
			user.IsLocked = false
			user.LoginCounter = 0
		} else {
			util.Response(c, "Your Account has been locked. Please try again after 2 hours", 401, nil, nil)
			return
		}

	}

	if user.LoginCounter >= 3 {
		user.IsLocked = true
		user.LockedAt = time.Now()
		user.UpdatedAt = time.Now()
		err = u.Repository.UpdateUser(user)
		if err != nil {
			util.Response(c, "There is an error occured", 500, err.Error(), nil)
			return
		}

		util.Response(c, "Your Account has been locked after 3 attempts. Please try again after 2 hours", 401, nil, nil)
		return
	}

	// Verify the password
	match := util.CheckPasswordHash(loginRequest.Password, user.Password)
	if !match {
		user.LoginCounter++
		user.UpdatedAt = time.Now()
		err = u.Repository.UpdateUser(user)
		if err != nil {
			util.Response(c, "There is an error occured", 500, err.Error(), nil)
			return
		}
		util.Response(c, "Incorrect password", 401, nil, nil)
		return
	}

	user.LoginCounter = 0
	user.UpdatedAt = time.Now()
	err = u.Repository.UpdateUser(user)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

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

	c.Header("access_token", *accessToken)
	c.Header("refresh_token", *refreshToken)

	util.Response(c, "Login successful", 200, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}, nil)
}

// Protected Route
func (u *HTTPHandler) GetUserByEmail(c *gin.Context) {

	_, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 500, err.Error(), nil)
		return
	}

	email := c.Query("email")
	email = strings.TrimSpace(email)
	if email == "" {
		util.Response(c, "Email is required", 400, nil, nil)
		return
	}

	user, err := u.Repository.FindUserByEmail(email)
	if err != nil {
		util.Response(c, "User not found", 404, err.Error(), nil)
		return
	}

	util.Response(c, "User Found", 200, user, nil)

}
