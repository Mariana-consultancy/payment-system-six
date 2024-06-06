package repository

import (
	"payment-system-six/internal/models"
	"strconv"
)

func (p *Postgres) CreateNotification(notification *models.Notification) error {
	if err := p.DB.Create(notification).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) GetNotificationsByUserID(userID uint) (*models.NotificationDetails, error) {
	notifications := &models.NotificationDetails{}
	if err := p.DB.Where("receiver_id = ?", userID).Find(&notifications.Notification).Error; err != nil {
		return nil, err
	}

	notifications.NotificationsCountTotal = len(notifications.Notification)
	notifications.NotificationMessage = "You have " + strconv.Itoa(notifications.NotificationsCountTotal) + " total notifications"

	notifications.NotificationsCountNew = 0
	for _, notification := range notifications.Notification {
		if notification.Status == "Unread" {
			notifications.NotificationsCountNew++
		}
	}
	notifications.NotificationMessageNew = "You have " + strconv.Itoa(notifications.NotificationsCountNew) + " new notifications"

	return notifications, nil
}

func (p *Postgres) GetNotificationByNotificationID(notificationID uint) (*models.Notification, error) {
	notification := &models.Notification{}
	if err := p.DB.Where("ID = ?", notificationID).First(&notification).Error; err != nil {
		return nil, err
	}
	return notification, nil
}

func (p *Postgres) UpdateNotification(notification *models.Notification) error {
	if err := p.DB.Save(notification).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) UpdateAllNotificationsByUserID(userID uint) error {
	notifications := &models.Notification{}
	if err := p.DB.Model(notifications).Where("receiver_id = ? AND deleted_at IS NULL", userID).Updates(models.Notification{Status: "Read"}).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteNotification(notification *models.Notification) error {
	if err := p.DB.Delete(notification).Error; err != nil {
		return err
	}
	return nil
}

func (p *Postgres) DeleteAllNotificationByUserID(userID uint) error {
	notifications := &[]models.Notification{}
	if err := p.DB.Where("receiver_id = ? AND deleted_at IS NULL", userID).Find(&notifications).Error; err != nil {
		return err
	}
	if len(*notifications) != 0 {
		if err := p.DB.Delete(notifications).Error; err != nil {
			return err
		}
		return nil
	}
	return nil
}
