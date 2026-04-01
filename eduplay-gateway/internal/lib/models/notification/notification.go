package model

import (
	dto "eduplay-gateway/internal/generated/clients/notification"
)

type NotificationFilter struct {
	Page      int64  `json:"page"`
	MaxOnPage int64  `json:"maxOnPage"`
	UserId    string `json:"userId"`
}

type Notifications struct {
	Notifications []NotificationInfo `json:"notifications"`
}

type NotificationInfo struct {
	NotificationId string             `json:"id"`
	Type           string             `json:"type"`
	Date           string             `json:"date"`
	FavEventStart  FavoriteEventStart `json:"favoriteEventStartExtra"`
	EvEnd          EventEnd           `json:"eventEndExtra"`
}

type FavoriteEventStart struct {
	EventId   string `json:"id"`
	EventName string `json:"eventName"`
}

type EventEnd struct {
	EventId            string `json:"id"`
	EventName          string `json:"eventName"`
	TimeLeft           string `json:"timeLeft"`
	NotStartedFavorite bool   `json:"notStartedFavorite"`
}

func NotificationsFromDto(notifications *dto.NotificationInfos) *Notifications {
	notifs := make([]NotificationInfo, 0)

	for _, notif := range notifications.Notifications {
		notifs = append(notifs, *NotificationFromDto(notif))
	}

	return &Notifications{Notifications: notifs}
}

func NotificationFromDto(notification *dto.NotificationInfo) *NotificationInfo {
	switch notification.Type {
	case "favoriteEventStart":
		return &NotificationInfo{
			NotificationId: notification.NotificationId,
			Type:           notification.Type,
			Date:           notification.Date.String(),
			FavEventStart: FavoriteEventStart{
				EventId:   notification.EventId,
				EventName: notification.EventName,
			},
		}
	case "eventEnd":
		return &NotificationInfo{
			NotificationId: notification.NotificationId,
			Type:           notification.Type,
			Date:           notification.Date.String(),
			EvEnd: EventEnd{
				EventId:            notification.EventId,
				EventName:          notification.EventName,
				TimeLeft:           notification.TimeLeft,
				NotStartedFavorite: notification.NotStartedFavorite,
			},
		}
	default:
		return nil
	}
}
