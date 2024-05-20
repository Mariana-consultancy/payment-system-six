package util

import (
	"net/http"
	"regexp"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// Response is customized to help return all responses need
func Response(c *gin.Context, message string, status int, data interface{}, errs []string) {
	responsedata := gin.H{
		"message":   message,
		"data":      data,
		"errors":    errs,
		"status":    http.StatusText(status),
		"timestamp": time.Now().Format("2006-01-02 15:04:05"),
	}

	c.IndentedJSON(status, responsedata)
}

// HashPassword takes a plaintext password and returns the hashed password or an error
func HashPassword(password string) (string, error) {

	// Use bcrypt to generate a hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash compares a plaintext password with its hashed version and returns a boolean value
func CheckPasswordHash(password, hashedPassword string) bool {

	// Compare the hashed password with the given password
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

// validatePassword checks if the password meets the specified criteria
func ValidatePassword(password string) bool {
	// Check the length of the password
	if len(password) < 6 {
		return false
	}

	// Compile regular expressions for each requirement
	uppercasePattern := `[A-Z]`
	numberPattern := `[0-9]`
	specialCharPattern := `[\W_]`

	// Create regex objects
	uppercaseRegex := regexp.MustCompile(uppercasePattern)
	numberRegex := regexp.MustCompile(numberPattern)
	specialCharRegex := regexp.MustCompile(specialCharPattern)

	// Check if the password contains at least one uppercase letter, one number, and one special character
	if !uppercaseRegex.MatchString(password) {
		return false
	}
	if !numberRegex.MatchString(password) {
		return false
	}
	if !specialCharRegex.MatchString(password) {
		return false
	}

	return true
}
