package models

import "gorm.io/gorm"

type Notification struct {
	gorm.Model
	ReceiverID       uint   `json:"receiver_id"`
	SenderID         uint   `json:"sender_id"`
	Title            string `json:"title"`
	Message          string `json:"message"`
	NotificationType string `json:"notification_type"`
	Status           string `json:"status"`
}

type NotificationDetails struct {
	NotificationsCountTotal int            `json:"notifications_count_total"`
	NotificationsCountNew   int            `json:"notifications_count_new"`
	NotificationMessage     string         `json:"notification_message"`
	NotificationMessageNew  string         `json:"notification_message_new"`
	Notification            []Notification `json:"notifications"`
}

type GetNotification struct {
	NotificationID uint `json:"notification_id"`
}
