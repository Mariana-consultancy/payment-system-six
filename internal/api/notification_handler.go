package api

import (
	"payment-system-six/internal/models"
	"payment-system-six/internal/util"

	"github.com/gin-gonic/gin"
)

func (u *HTTPHandler) GetAllNotifications(c *gin.Context) {

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	allNotifications, err := u.Repository.GetNotificationsByUserID(user.ID)
	if err != nil {
		util.Response(c, "Internal server error", 500, err.Error(), nil)
		return
	}
	if len(allNotifications.Notification) == 0 {
		util.Response(c, "No record found", 200, allNotifications, nil)
		return
	}
	util.Response(c, "Record found", 200, allNotifications, nil)

}

func (u *HTTPHandler) GetNotification(c *gin.Context) {
	var getNotification models.GetNotification

	if err := c.ShouldBind(&getNotification); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	notification, err := u.Repository.GetNotificationByNotificationID(getNotification.NotificationID)
	if err != nil {
		util.Response(c, "Notification ID does not exist", 400, err.Error(), nil)
		return
	}

	util.Response(c, "Record found", 200, notification, nil)

}

func (u *HTTPHandler) ReadNotification(c *gin.Context) {
	var getNotification models.GetNotification

	if err := c.ShouldBind(&getNotification); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	notification, err := u.Repository.GetNotificationByNotificationID(getNotification.NotificationID)
	if err != nil {
		util.Response(c, "Notification ID does not exist", 404, err.Error(), nil)
		return
	}

	notification.Status = "Read"
	err = u.Repository.UpdateNotification(notification)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Record updated", 200, nil, nil)

}

func (u *HTTPHandler) ReadAllNotifications(c *gin.Context) {

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	err = u.Repository.UpdateAllNotificationsByUserID(user.ID)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Records Updated", 200, nil, nil)

}

func (u *HTTPHandler) DeleteNotification(c *gin.Context) {
	var getNotification models.GetNotification

	if err := c.ShouldBind(&getNotification); err != nil {
		util.Response(c, "invalid request", 400, err.Error(), nil)
		return
	}

	notification, err := u.Repository.GetNotificationByNotificationID(getNotification.NotificationID)
	if err != nil {
		util.Response(c, "Notification ID does not exist", 404, err.Error(), nil)
		return
	}

	err = u.Repository.DeleteNotification(notification)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Record Deleted", 200, nil, nil)

}

func (u *HTTPHandler) DeleteAllNotifications(c *gin.Context) {

	user, err := u.GetUserFromContext(c)
	if err != nil {
		util.Response(c, "User not logged in", 400, err.Error(), nil)
		return
	}

	err = u.Repository.DeleteAllNotificationByUserID(user.ID)
	if err != nil {
		util.Response(c, "There is an error occured", 500, err.Error(), nil)
		return
	}

	util.Response(c, "Records Deleted", 200, nil, nil)

}
