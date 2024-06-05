package api

import (
	"payment-system-six/internal/util"
	"strings"

	"github.com/gin-gonic/gin"
)

// Protected Route
func (u *HTTPHandler) GetUserByEmail(c *gin.Context) {

	_, err := u.GetAdminFromContext(c)
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

func (u *HTTPHandler) GetAllUsers(c *gin.Context) {

	_, err := u.GetAdminFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	allUsers, err := u.Repository.FindAllUsers()
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	if len(allUsers) == 0 {
		util.Response(c, "No record found", 200, allUsers, nil)
		return
	}
	util.Response(c, "Record found", 200, allUsers, nil)
}
